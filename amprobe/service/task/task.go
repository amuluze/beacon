// Package task
package task

import (
	"context"
	"fmt"
	"log/slog"
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
	threshold, err := t.alarmThreshold(ctx, "cpu")
	if err != nil {
		return err
	}

	agentIDs, err := t.agentIDs(ctx, &model.MonitorCPU{})
	if err != nil {
		return err
	}

	for _, agentID := range agentIDs {
		if err := t.cpuAlarmTaskForAgent(ctx, agentID, threshold); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) cpuAlarmTaskForAgent(ctx context.Context, agentID string, threshold model.AlarmThreshold) error {
	// Read CPU data from local DB
	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute)
	var cpuData []model.MonitorCPU
	if err := t.db.WithContext(ctx).Model(&model.MonitorCPU{}).
		Where("agent_id = ? AND timestamp > ?", agentID, startTime).
		Order("timestamp asc").
		Find(&cpuData).Error; err != nil {
		return err
	}

	total := 0.0
	for _, item := range cpuData {
		total += item.CPUPercent
	}
	if len(cpuData) > 0 && int(utils.Decimal(total/float64(len(cpuData)))*100) > threshold.Threshold {
		return t.triggerAlarm(
			alarmCacheKey(agentID, "cpu"),
			fmt.Sprintf("%s CPU 使用率连续 %d 分钟超过 %d%%", t.agentLabel(ctx, agentID), threshold.Duration, threshold.Threshold),
		)
	}
	return nil
}

func (t *Task) MemoryAlarmTask(ctx context.Context) error {
	threshold, err := t.alarmThreshold(ctx, "memory")
	if err != nil {
		return err
	}

	agentIDs, err := t.agentIDs(ctx, &model.MonitorMemory{})
	if err != nil {
		return err
	}

	for _, agentID := range agentIDs {
		if err := t.memoryAlarmTaskForAgent(ctx, agentID, threshold); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) memoryAlarmTaskForAgent(ctx context.Context, agentID string, threshold model.AlarmThreshold) error {
	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute)
	var memData []model.MonitorMemory
	if err := t.db.WithContext(ctx).Model(&model.MonitorMemory{}).
		Where("agent_id = ? AND timestamp > ?", agentID, startTime).
		Order("timestamp asc").
		Find(&memData).Error; err != nil {
		return err
	}

	total := 0.0
	for _, item := range memData {
		total += item.MemPercent
	}
	if len(memData) > 0 && int(utils.Decimal(total/float64(len(memData)))*100) > threshold.Threshold {
		return t.triggerAlarm(
			alarmCacheKey(agentID, "memory"),
			fmt.Sprintf("%s 内存使用率连续 %d 分钟超过 %d%%", t.agentLabel(ctx, agentID), threshold.Duration, threshold.Threshold),
		)
	}
	return nil
}

func (t *Task) DiskAlarmTask(ctx context.Context) error {
	threshold, err := t.alarmThreshold(ctx, "disk")
	if err != nil {
		return err
	}

	agentIDs, err := t.agentIDs(ctx, &model.MonitorDisk{})
	if err != nil {
		return err
	}

	for _, agentID := range agentIDs {
		if err := t.diskAlarmTaskForAgent(ctx, agentID, threshold); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) diskAlarmTaskForAgent(ctx context.Context, agentID string, threshold model.AlarmThreshold) error {
	// Get latest disk info by device
	var diskData []model.MonitorDisk
	if err := t.db.WithContext(ctx).Model(&model.MonitorDisk{}).
		Where("agent_id = ?", agentID).
		Order("timestamp desc").
		Find(&diskData).Error; err != nil {
		return err
	}

	diskMap := make(map[string]struct{})
	for _, item := range diskData {
		if _, ok := diskMap[item.Device]; ok {
			continue
		}
		diskMap[item.Device] = struct{}{}
		if int(utils.Decimal(item.DiskPercent)*100) > threshold.Threshold {
			return t.triggerAlarm(
				alarmCacheKey(agentID, "disk", item.Device),
				fmt.Sprintf("%s 磁盘 %s 使用率超过 %d%%", t.agentLabel(ctx, agentID), item.Device, threshold.Threshold),
			)
		}
	}
	return nil
}

func (t *Task) ServiceTask(ctx context.Context) error {
	// Read latest container states from local DB
	var containers []model.MonitorContainer
	if err := t.db.WithContext(ctx).Model(&model.MonitorContainer{}).Order("agent_id asc, name asc, timestamp desc").Find(&containers).Error; err != nil {
		return err
	}

	for _, item := range containers {
		if item.AgentID == "" {
			continue
		}
		stateKey := alarmCacheKey(item.AgentID, "service-state", item.Name)
		if containerStateBytes, ok := t.cache.Get(stateKey); ok {
			if containerStateBytes.(string) != item.State {
				msg := fmt.Sprintf("%s 容器 %s 的状态由 %s 变为 %s", t.agentLabel(ctx, item.AgentID), item.Name, containerStateBytes.(string), item.State)
				if err := t.triggerAlarm(alarmCacheKey(item.AgentID, "service", item.Name), msg); err != nil {
					return err
				}
			}
		}
		t.cache.Set(stateKey, item.State, 0)
	}
	return nil
}

func (t *Task) alarmThreshold(ctx context.Context, alarmType string) (model.AlarmThreshold, error) {
	var threshold model.AlarmThreshold
	err := t.db.WithContext(ctx).Model(&model.AlarmThreshold{}).Where("type = ?", alarmType).First(&threshold).Error
	return threshold, err
}

func (t *Task) agentIDs(ctx context.Context, model interface{}) ([]string, error) {
	var agentIDs []string
	if err := t.db.WithContext(ctx).Model(model).
		Distinct("agent_id").
		Where("agent_id <> ?", "").
		Order("agent_id asc").
		Pluck("agent_id", &agentIDs).Error; err != nil {
		return nil, err
	}
	return agentIDs, nil
}

func (t *Task) agentLabel(ctx context.Context, agentID string) string {
	var hostInfo model.MonitorHost
	if err := t.db.WithContext(ctx).Model(&model.MonitorHost{}).
		Where("agent_id = ?", agentID).
		Order("timestamp desc").
		First(&hostInfo).Error; err == nil && hostInfo.Hostname != "" {
		return fmt.Sprintf("Agent %s(%s)", agentID, hostInfo.Hostname)
	}
	return fmt.Sprintf("Agent %s", agentID)
}

func alarmCacheKey(parts ...string) string {
	return strings.Join(parts, ":")
}

func (t *Task) triggerAlarm(key string, msg string) error {
	if err := t.sendAlarmAudit(msg); err != nil {
		return err
	}
	if _, ok := t.cache.Get(key); ok {
		return nil
	}
	if err := t.sendMail(msg); err != nil {
		slog.Error("send alarm mail failed", "key", key, "error", err)
		return nil
	}
	t.cache.Set(key, "true", 10*time.Minute)
	return nil
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
