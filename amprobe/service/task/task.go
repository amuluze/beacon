// Package task
package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"amprobe/pkg/utils"
	"amprobe/service/model"
	"common/database"

	"github.com/patrickmn/go-cache"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var _ ITask = (*Task)(nil)

type ITask interface {
	CPUAlarmTask(context.Context) error
	MemoryAlarmTask(context.Context) error
	DiskAlarmTask(context.Context) error
	ServiceTask(context.Context) error
}

type Task struct {
	db    *database.DB
	cache *cache.Cache
}

func NewTask(db *database.DB) *Task {
	return &Task{
		db:    db,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (t *Task) CPUAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}

	// Read host info from local DB
	var hostInfo model.MonitorHost
	var hostname string
	if err := t.db.Model(&model.MonitorHost{}).Order("timestamp desc").First(&hostInfo).Error; err == nil {
		hostname = hostInfo.Hostname
	}

	// Read CPU data from local DB
	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute).Unix()
	var cpuData []model.MonitorCPU
	if err := t.db.Model(&model.MonitorCPU{}).Where("timestamp > ?", time.Unix(startTime, 0)).Order("timestamp asc").Find(&cpuData).Error; err != nil {
		return err
	}

	total := 0.0
	for _, item := range cpuData {
		total += item.CPUPercent
	}
	if len(cpuData) > 0 && int(utils.Decimal(total/float64(len(cpuData)))*100) > threshold.Threshold {
		return t.triggerAlarm("cpu", fmt.Sprintf("%s CPU 使用率连续 %d 分钟超过 %d%%", hostname, threshold.Duration, threshold.Threshold))
	}
	return nil
}

func (t *Task) MemoryAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}

	var hostInfo model.MonitorHost
	var hostname string
	if err := t.db.Model(&model.MonitorHost{}).Order("timestamp desc").First(&hostInfo).Error; err == nil {
		hostname = hostInfo.Hostname
	}

	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute).Unix()
	var memData []model.MonitorMemory
	if err := t.db.Model(&model.MonitorMemory{}).Where("timestamp > ?", time.Unix(startTime, 0)).Order("timestamp asc").Find(&memData).Error; err != nil {
		return err
	}

	total := 0.0
	for _, item := range memData {
		total += item.MemPercent
	}
	if len(memData) > 0 && int(utils.Decimal(total/float64(len(memData)))*100) > threshold.Threshold {
		return t.triggerAlarm("memory", fmt.Sprintf("%s 内存使用率连续 %d 分钟超过 %d%%", hostname, threshold.Duration, threshold.Threshold))
	}
	return nil
}

func (t *Task) DiskAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}

	var hostInfo model.MonitorHost
	var hostname string
	if err := t.db.Model(&model.MonitorHost{}).Order("timestamp desc").First(&hostInfo).Error; err == nil {
		hostname = hostInfo.Hostname
	}

	// Get latest disk info by device
	var diskData []model.MonitorDisk
	if err := t.db.Model(&model.MonitorDisk{}).Group("device").Order("timestamp desc").Find(&diskData).Error; err != nil {
		return err
	}

	diskMap := make(map[string]struct{})
	for _, item := range diskData {
		if _, ok := diskMap[item.Device]; ok {
			continue
		}
		diskMap[item.Device] = struct{}{}
		if int(utils.Decimal(item.DiskPercent)*100) > threshold.Threshold {
			return t.triggerAlarm("disk", fmt.Sprintf("%s 磁盘 %s 使用率超过 %d%%", hostname, item.Device, threshold.Threshold))
		}
	}
	return nil
}

func (t *Task) ServiceTask(ctx context.Context) error {
	// Read latest container states from local DB
	var containers []model.MonitorContainer
	if err := t.db.Model(&model.MonitorContainer{}).Order("created_at desc").Find(&containers).Error; err != nil {
		return err
	}

	for _, item := range containers {
		if containerStateBytes, ok := t.cache.Get(item.Name); ok {
			if containerStateBytes.(string) != item.State {
				msg := fmt.Sprintf("容器 %s 的状态由 %s 变为 %s", item.Name, containerStateBytes.(string), item.State)
				if err := t.sendAlarmAudit(msg); err != nil {
					return err
				}
			}
		}
		t.cache.Set(item.Name, item.State, 0)
	}
	return nil
}

func (t *Task) triggerAlarm(key string, msg string) error {
	return t.db.RunInTransaction(func(tx *gorm.DB) error {
		operateLog := model.Audit{
			Username: "system",
			Operate:  msg,
		}
		if err := tx.Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
			return err
		}
		if _, ok := t.cache.Get(key); ok {
			return nil
		}
		if err := t.sendMail(msg); err != nil {
			return err
		}
		t.cache.Set(key, "true", 10*time.Minute)
		return nil
	})
}

func (t *Task) sendAlarmAudit(msg string) error {
	return t.db.RunInTransaction(func(tx *gorm.DB) error {
		operateLog := model.Audit{
			Username: "system",
			Operate:  msg,
		}
		if err := tx.Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
			return err
		}
		if err := t.sendMail(msg); err != nil {
			return err
		}
		return nil
	})
}

func (t *Task) sendMail(msg string) error {
	var mail model.Mail
	if err := t.db.Model(&model.Mail{}).First(&mail).Error; err != nil {
		return err
	}
	dialer := gomail.NewDialer(mail.Server, mail.Port, mail.Sender, mail.Password)
	for _, recv := range strings.Split(mail.Receiver, ",") {
		mailMessage := gomail.NewMessage()
		mailMessage.SetHeader("From", mail.Sender)
		mailMessage.SetHeader("To", recv)
		mailMessage.SetHeader("Subject", "服务器告警")
		mailMessage.SetBody("text/plain", msg)

		if err := dialer.DialAndSend(mailMessage); err != nil {
			return err
		}
	}
	return nil
}
