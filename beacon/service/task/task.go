// Package task
package task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"beacon/pkg/secret"
	"beacon/pkg/utils"
	dingtalkclient "beacon/service/dingtalk/client"
	"beacon/service/model"
	"common/database"

	"github.com/patrickmn/go-cache"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var _ ITask = (*Task)(nil)

const (
	alarmDeduplicationWindow = 10 * time.Minute
	alarmSampleFreshness     = 90 * time.Second
	alarmWindowGrace         = 45 * time.Second
	alarmMaximumSampleGap    = 90 * time.Second
	containerStateFreshness  = 2 * time.Minute
)

type alarmWindowState uint8

const (
	alarmWindowUnknown alarmWindowState = iota
	alarmWindowNormal
	alarmWindowFiring
)

type metricSample struct {
	observedAt time.Time
	value      float64
}

type ITask interface {
	CPUAlarmTask(context.Context) error
	MemoryAlarmTask(context.Context) error
	DiskAlarmTask(context.Context) error
	ServiceTask(context.Context) error
}

// MailSender abstracts the email delivery so unit tests can stub it out
// without dialing a real SMTP server.
type MailSender interface {
	Send(title string, msg string) error
}

type Task struct {
	db             *database.DB
	cache          *cache.Cache
	mailSender     MailSender
	dingTalkSender dingtalkclient.Sender
	cipher         *secret.Cipher
	alarmMu        sync.Mutex
	now            func() time.Time
}

func NewTask(db *database.DB, ciphers ...*secret.Cipher) *Task {
	task := &Task{
		db:             db,
		cache:          cache.New(5*time.Minute, 10*time.Minute),
		dingTalkSender: dingtalkclient.NewSender(),
		now:            time.Now,
	}
	if len(ciphers) > 0 {
		task.cipher = ciphers[0]
	}
	return task
}

// SetMailSender overrides the default SMTP-based mail sender. Useful for
// tests; not expected to be called in production.
func (t *Task) SetMailSender(s MailSender) {
	t.mailSender = s
}

// SetDingTalkSender overrides the HTTP sender for focused notification tests.
func (t *Task) SetDingTalkSender(sender dingtalkclient.Sender) {
	t.dingTalkSender = sender
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

func alarmQueryStart(now time.Time, duration int) time.Time {
	if duration <= 0 {
		return now.Add(-alarmSampleFreshness)
	}
	return now.Add(-time.Duration(duration) * time.Minute)
}

// evaluateAlarmWindow distinguishes a valid normal window from missing or
// stale data. For a positive duration, samples must cover the whole window,
// arrive without a long gap, and all remain strictly above the threshold.
func evaluateAlarmWindow(samples []metricSample, now time.Time, duration int, threshold int) alarmWindowState {
	if len(samples) == 0 {
		return alarmWindowUnknown
	}

	if duration <= 0 {
		latest := samples[len(samples)-1]
		if latest.observedAt.Before(now.Add(-alarmSampleFreshness)) || latest.observedAt.After(now.Add(alarmWindowGrace)) {
			return alarmWindowUnknown
		}
		if utils.Decimal(latest.value) > float64(threshold) {
			return alarmWindowFiring
		}
		return alarmWindowNormal
	}

	windowDuration := time.Duration(duration) * time.Minute
	windowStart := now.Add(-windowDuration)
	grace := alarmWindowGrace
	if quarterWindow := windowDuration / 4; quarterWindow < grace {
		grace = quarterWindow
	}
	first := samples[0]
	latest := samples[len(samples)-1]
	if first.observedAt.After(windowStart.Add(grace)) || latest.observedAt.Before(now.Add(-grace)) || latest.observedAt.After(now.Add(alarmWindowGrace)) {
		return alarmWindowUnknown
	}

	previous := first.observedAt
	for index, sample := range samples {
		if index > 0 && sample.observedAt.Sub(previous) > alarmMaximumSampleGap {
			return alarmWindowUnknown
		}
		previous = sample.observedAt
		if utils.Decimal(sample.value) <= float64(threshold) {
			return alarmWindowNormal
		}
	}
	return alarmWindowFiring
}

func cpuSamples(rows []model.MonitorCPU) []metricSample {
	samples := make([]metricSample, 0, len(rows))
	for _, row := range rows {
		samples = append(samples, metricSample{observedAt: row.CreatedAt, value: row.CPUPercent})
	}
	return samples
}

func memorySamples(rows []model.MonitorMemory) []metricSample {
	samples := make([]metricSample, 0, len(rows))
	for _, row := range rows {
		samples = append(samples, metricSample{observedAt: row.CreatedAt, value: row.MemPercent})
	}
	return samples
}

func exceededThresholdMessage(resource string, duration int, threshold int) string {
	if duration <= 0 {
		return fmt.Sprintf("%s 使用率超过 %d%%", resource, threshold)
	}
	return fmt.Sprintf("%s 使用率连续 %d 分钟超过 %d%%", resource, duration, threshold)
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

	now := t.now()
	startTime := alarmQueryStart(now, threshold.Duration)
	for _, agentID := range agentIDs {
		var cpuData []model.MonitorCPU
		if err := t.db.WithContext(ctx).Model(&model.MonitorCPU{}).
			Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agentID, startTime, now.Add(alarmWindowGrace)).
			Order("created_at asc").Find(&cpuData).Error; err != nil {
			return err
		}

		key := "cpu:" + agentID
		switch evaluateAlarmWindow(cpuSamples(cpuData), now, threshold.Duration, threshold.Threshold) {
		case alarmWindowFiring:
			msg := exceededThresholdMessage(fmt.Sprintf("[%s] %s CPU", agentID, t.hostname(ctx, agentID)), threshold.Duration, threshold.Threshold)
			if err := t.triggerAlarm(ctx, key, msg, agentID); err != nil {
				return err
			}
		case alarmWindowNormal:
			msg := fmt.Sprintf("[%s] %s CPU 使用率已恢复至 %d%% 以下", agentID, t.hostname(ctx, agentID), threshold.Threshold)
			if err := t.resolveAlarm(ctx, key, msg, agentID); err != nil {
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

	now := t.now()
	startTime := alarmQueryStart(now, threshold.Duration)
	for _, agentID := range agentIDs {
		var memData []model.MonitorMemory
		if err := t.db.WithContext(ctx).Model(&model.MonitorMemory{}).
			Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agentID, startTime, now.Add(alarmWindowGrace)).
			Order("created_at asc").Find(&memData).Error; err != nil {
			return err
		}

		key := "memory:" + agentID
		switch evaluateAlarmWindow(memorySamples(memData), now, threshold.Duration, threshold.Threshold) {
		case alarmWindowFiring:
			msg := exceededThresholdMessage(fmt.Sprintf("[%s] %s 内存", agentID, t.hostname(ctx, agentID)), threshold.Duration, threshold.Threshold)
			if err := t.triggerAlarm(ctx, key, msg, agentID); err != nil {
				return err
			}
		case alarmWindowNormal:
			msg := fmt.Sprintf("[%s] %s 内存使用率已恢复至 %d%% 以下", agentID, t.hostname(ctx, agentID), threshold.Threshold)
			if err := t.resolveAlarm(ctx, key, msg, agentID); err != nil {
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

	now := t.now()
	startTime := alarmQueryStart(now, threshold.Duration)
	for _, agentID := range agentIDs {
		var diskData []model.MonitorDisk
		if err := t.db.WithContext(ctx).Model(&model.MonitorDisk{}).
			Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agentID, startTime, now.Add(alarmWindowGrace)).
			Order("device asc, created_at asc").
			Find(&diskData).Error; err != nil {
			return err
		}

		deviceOrder := make([]string, 0)
		deviceSamples := make(map[string][]metricSample)
		for _, item := range diskData {
			if _, exists := deviceSamples[item.Device]; !exists {
				deviceOrder = append(deviceOrder, item.Device)
			}
			deviceSamples[item.Device] = append(deviceSamples[item.Device], metricSample{observedAt: item.CreatedAt, value: item.DiskPercent})
		}
		for _, device := range deviceOrder {
			key := "disk:" + agentID + ":" + device
			switch evaluateAlarmWindow(deviceSamples[device], now, threshold.Duration, threshold.Threshold) {
			case alarmWindowFiring:
				msg := exceededThresholdMessage(fmt.Sprintf("[%s] %s 磁盘 %s", agentID, t.hostname(ctx, agentID), device), threshold.Duration, threshold.Threshold)
				if err := t.triggerAlarm(ctx, key, msg, agentID); err != nil {
					return err
				}
			case alarmWindowNormal:
				msg := fmt.Sprintf("[%s] %s 磁盘 %s 使用率已恢复至 %d%% 以下", agentID, t.hostname(ctx, agentID), device, threshold.Threshold)
				if err := t.resolveAlarm(ctx, key, msg, agentID); err != nil {
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

	now := t.now()
	freshAfter := now.Add(-containerStateFreshness)
	for _, agentID := range agentIDs {
		var containers []model.MonitorContainer
		latest := t.db.WithContext(ctx).Model(&model.MonitorContainer{}).
			Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agentID, freshAfter, now.Add(alarmWindowGrace)).
			Select("agent_id, name, MAX(created_at) AS created_at").
			Group("agent_id, name")
		if err := t.db.WithContext(ctx).Model(&model.MonitorContainer{}).
			Where("m_container.agent_id = ? AND m_container.created_at >= ? AND m_container.created_at <= ?", agentID, freshAfter, now.Add(alarmWindowGrace)).
			Joins("JOIN (?) latest ON latest.agent_id = m_container.agent_id AND latest.name = m_container.name AND latest.created_at = m_container.created_at", latest).
			Find(&containers).Error; err != nil {
			return err
		}

		for _, item := range containers {
			cacheKey := "container:" + agentID + ":" + item.Name
			if containerStateBytes, ok := t.cache.Get(cacheKey); ok {
				if containerStateBytes.(string) != item.State {
					msg := fmt.Sprintf("[%s] 容器 %s 的状态由 %s 变为 %s", agentID, item.Name, containerStateBytes.(string), item.State)
					if err := t.sendAlarmAudit(ctx, msg, agentID); err != nil {
						return err
					}
				}
			}
			t.cache.Set(cacheKey, item.State, 0)
		}
	}
	return nil
}

// triggerAlarm persists the alarm before attempting each notification channel.
// Delivery failures are isolated from each other and never roll back the audit.
func (t *Task) triggerAlarm(ctx context.Context, key string, msg string, agentID string) error {
	t.alarmMu.Lock()
	if _, ok := t.cache.Get(key); ok {
		t.cache.Set("alarm-state:"+key, "firing", cache.NoExpiration)
		t.alarmMu.Unlock()
		return nil
	}
	operateLog := model.Audit{
		Username: "system",
		AgentID:  agentID,
		Operate:  msg,
	}
	if err := t.db.WithContext(ctx).Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
		t.alarmMu.Unlock()
		return err
	}
	t.cache.Set(key, "true", alarmDeduplicationWindow)
	t.cache.Set("alarm-state:"+key, "firing", cache.NoExpiration)
	t.alarmMu.Unlock()
	t.notify(ctx, msg)
	return nil
}

func (t *Task) resolveAlarm(ctx context.Context, key string, msg string, agentID string) error {
	stateKey := "alarm-state:" + key
	t.alarmMu.Lock()
	if _, ok := t.cache.Get(stateKey); !ok {
		t.alarmMu.Unlock()
		return nil
	}
	operateLog := model.Audit{
		Username: "system",
		AgentID:  agentID,
		Operate:  msg,
	}
	if err := t.db.WithContext(ctx).Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
		t.alarmMu.Unlock()
		return err
	}
	t.cache.Delete(stateKey)
	t.cache.Delete(key)
	t.alarmMu.Unlock()
	t.notify(ctx, msg)
	return nil
}

func (t *Task) sendAlarmAudit(ctx context.Context, msg string, agentID string) error {
	operateLog := model.Audit{
		Username: "system",
		AgentID:  agentID,
		Operate:  msg,
	}
	if err := t.db.WithContext(ctx).Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
		return err
	}
	t.notify(ctx, msg)
	return nil
}

func (t *Task) notify(ctx context.Context, msg string) {
	title := "服务器告警"
	if err := t.sendMail(ctx, title, msg); err != nil {
		slog.Warn("send alarm notification failed", "channel", "mail", "err", err)
	}
	if err := t.sendDingTalk(ctx, title, msg); err != nil {
		slog.Warn("send alarm notification failed", "channel", "dingtalk", "err", err)
	}
}

func (t *Task) sendMail(ctx context.Context, title string, msg string) error {
	// If a test injected a MailSender, prefer it (avoids dialing a real SMTP
	// server during unit tests).
	if t.mailSender != nil {
		return t.mailSender.Send(title, msg)
	}
	var mail model.Mail
	if err := t.db.WithContext(ctx).Model(&model.Mail{}).First(&mail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	password, err := t.decryptAndMigrate(ctx, &mail, "password", mail.Password)
	if err != nil {
		return err
	}
	mail.Password = password
	dialer := gomail.NewDialer(mail.Server, mail.Port, mail.Sender, mail.Password)
	var sendErrors []error
	for _, recv := range strings.Split(mail.Receiver, ",") {
		recv = strings.TrimSpace(recv)
		if recv == "" {
			continue
		}
		mailMessage := gomail.NewMessage()
		mailMessage.SetHeader("From", mail.Sender)
		mailMessage.SetHeader("To", recv)
		mailMessage.SetHeader("Subject", title)
		mailMessage.SetBody("text/plain", msg)

		if err := dialer.DialAndSend(mailMessage); err != nil {
			sendErrors = append(sendErrors, err)
		}
	}
	return errors.Join(sendErrors...)
}

func (t *Task) sendDingTalk(ctx context.Context, title string, msg string) error {
	var setting model.DingTalk
	err := t.db.WithContext(ctx).
		Where("key = ?", model.DefaultDingTalkConfigKey).
		First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	webhook, err := t.decryptAndMigrate(ctx, &setting, "webhook", setting.Webhook)
	if err != nil {
		return err
	}
	secretValue, err := t.decryptAndMigrate(ctx, &setting, "secret", setting.Secret)
	if err != nil {
		return err
	}
	setting.Webhook = webhook
	setting.Secret = secretValue
	if !setting.Enabled || setting.Webhook == "" {
		return nil
	}
	return t.dingTalkSender.Send(ctx, dingtalkclient.Config{
		Webhook:      setting.Webhook,
		Secret:       setting.Secret,
		AtAll:        setting.AtAll,
	}, msg)
}

func (t *Task) decryptAndMigrate(ctx context.Context, target any, column, value string) (string, error) {
	if t.cipher == nil {
		return value, nil
	}
	plaintext, legacy, err := t.cipher.Decrypt(value)
	if err != nil {
		return "", err
	}
	if !legacy {
		return plaintext, nil
	}
	encrypted, err := t.cipher.Encrypt(plaintext)
	if err != nil {
		return "", err
	}
	if err := t.db.WithContext(ctx).Model(target).Update(column, encrypted).Error; err != nil {
		return "", err
	}
	return plaintext, nil
}
