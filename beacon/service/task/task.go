// Package task
package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"beacon/pkg/utils"
	"beacon/service/model"
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

func (t *Task) threshold(alarmType string) (model.AlarmThreshold, error) {
	var threshold model.AlarmThreshold
	err := t.db.Model(&model.AlarmThreshold{}).Where("type = ?", alarmType).First(&threshold).Error
	return threshold, err
}

func (t *Task) agentIDs(ctx context.Context) ([]string, error) {
	var agents []model.Agent
	if err := t.db.WithContext(ctx).Model(&model.Agent{}).Order("agent_id asc").Find(&agents).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(agents))
	seen := make(map[string]struct{})
	for _, agent := range agents {
		if agent.AgentID == "" {
			continue
		}
		ids = append(ids, agent.AgentID)
		seen[agent.AgentID] = struct{}{}
	}
	if len(ids) > 0 {
		return ids, nil
	}
	var monitorIDs []string
	if err := t.db.WithContext(ctx).Model(&model.MonitorHost{}).
		Where("agent_id <> ?", "").
		Distinct("agent_id").
		Pluck("agent_id", &monitorIDs).Error; err != nil {
		return nil, err
	}
	for _, agentID := range monitorIDs {
		if _, ok := seen[agentID]; !ok {
			ids = append(ids, agentID)
		}
	}
	return ids, nil
}

func (t *Task) hostname(ctx context.Context, agentID string) string {
	var hostInfo model.MonitorHost
	if err := t.db.WithContext(ctx).Model(&model.MonitorHost{}).
		Where("agent_id = ?", agentID).
		Order("timestamp desc").
		First(&hostInfo).Error; err == nil {
		return hostInfo.Hostname
	}
	return agentID
}

func (t *Task) CPUAlarmTask(ctx context.Context) error {
	threshold, err := t.threshold("cpu")
	if err != nil {
		return err
	}
	agentIDs, err := t.agentIDs(ctx)
	if err != nil {
		return err
	}

	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute).Unix()
	for _, agentID := range agentIDs {
		var cpuData []model.MonitorCPU
		if err := t.db.WithContext(ctx).Model(&model.MonitorCPU{}).
			Where("agent_id = ? AND timestamp > ?", agentID, time.Unix(startTime, 0)).
			Order("timestamp asc").Find(&cpuData).Error; err != nil {
			return err
		}

		total := 0.0
		for _, item := range cpuData {
			total += item.CPUPercent
		}
		if len(cpuData) > 0 && int(utils.Decimal(total/float64(len(cpuData)))*100) > threshold.Threshold {
			msg := fmt.Sprintf("[%s] %s CPU 使用率连续 %d 分钟超过 %d%%", agentID, t.hostname(ctx, agentID), threshold.Duration, threshold.Threshold)
			if err := t.triggerAlarm("cpu:"+agentID, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Task) MemoryAlarmTask(ctx context.Context) error {
	threshold, err := t.threshold("memory")
	if err != nil {
		return err
	}
	agentIDs, err := t.agentIDs(ctx)
	if err != nil {
		return err
	}

	startTime := time.Now().Add(-time.Duration(threshold.Duration) * time.Minute).Unix()
	for _, agentID := range agentIDs {
		var memData []model.MonitorMemory
		if err := t.db.WithContext(ctx).Model(&model.MonitorMemory{}).
			Where("agent_id = ? AND timestamp > ?", agentID, time.Unix(startTime, 0)).
			Order("timestamp asc").Find(&memData).Error; err != nil {
			return err
		}

		total := 0.0
		for _, item := range memData {
			total += item.MemPercent
		}
		if len(memData) > 0 && int(utils.Decimal(total/float64(len(memData)))*100) > threshold.Threshold {
			msg := fmt.Sprintf("[%s] %s 内存使用率连续 %d 分钟超过 %d%%", agentID, t.hostname(ctx, agentID), threshold.Duration, threshold.Threshold)
			if err := t.triggerAlarm("memory:"+agentID, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Task) DiskAlarmTask(ctx context.Context) error {
	threshold, err := t.threshold("disk")
	if err != nil {
		return err
	}
	agentIDs, err := t.agentIDs(ctx)
	if err != nil {
		return err
	}

	for _, agentID := range agentIDs {
		var diskData []model.MonitorDisk
		latest := t.db.WithContext(ctx).Model(&model.MonitorDisk{}).
			Where("agent_id = ?", agentID).
			Select("agent_id, device, MAX(timestamp) AS timestamp").
			Group("agent_id, device")
		if err := t.db.WithContext(ctx).Model(&model.MonitorDisk{}).
			Where("m_disk.agent_id = ?", agentID).
			Joins("JOIN (?) latest ON latest.agent_id = m_disk.agent_id AND latest.device = m_disk.device AND latest.timestamp = m_disk.timestamp", latest).
			Find(&diskData).Error; err != nil {
			return err
		}

		for _, item := range diskData {
			if int(utils.Decimal(item.DiskPercent)*100) > threshold.Threshold {
				msg := fmt.Sprintf("[%s] %s 磁盘 %s 使用率超过 %d%%", agentID, t.hostname(ctx, agentID), item.Device, threshold.Threshold)
				if err := t.triggerAlarm("disk:"+agentID+":"+item.Device, msg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *Task) ServiceTask(ctx context.Context) error {
	agentIDs, err := t.agentIDs(ctx)
	if err != nil {
		return err
	}

	for _, agentID := range agentIDs {
		var containers []model.MonitorContainer
		latest := t.db.WithContext(ctx).Model(&model.MonitorContainer{}).
			Where("agent_id = ?", agentID).
			Select("agent_id, name, MAX(timestamp) AS timestamp").
			Group("agent_id, name")
		if err := t.db.WithContext(ctx).Model(&model.MonitorContainer{}).
			Where("m_container.agent_id = ?", agentID).
			Joins("JOIN (?) latest ON latest.agent_id = m_container.agent_id AND latest.name = m_container.name AND latest.timestamp = m_container.timestamp", latest).
			Find(&containers).Error; err != nil {
			return err
		}

		for _, item := range containers {
			cacheKey := "container:" + agentID + ":" + item.Name
			if containerStateBytes, ok := t.cache.Get(cacheKey); ok {
				if containerStateBytes.(string) != item.State {
					msg := fmt.Sprintf("[%s] 容器 %s 的状态由 %s 变为 %s", agentID, item.Name, containerStateBytes.(string), item.State)
					if err := t.sendAlarmAudit(msg); err != nil {
						return err
					}
				}
			}
			t.cache.Set(cacheKey, item.State, 0)
		}
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
