// Package task
// Date:   2024/10/15 11:19
// Author: Amu
// Description:
package task

import (
	"amprobe/pkg/rpc"
	"amprobe/pkg/utils"
	"amprobe/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
	"context"
	"fmt"
	"strings"
	"time"

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
	rpcClient *rpc.Client
	db        *database.DB
	cache     *cache.Cache
}

func NewTask(cli *rpc.Client, db *database.DB) *Task {
	return &Task{
		rpcClient: cli,
		db:        db,
		cache:     cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (t *Task) CPUAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}
	hostArgs := rpcSchema.HostInfoArgs{}
	var hostReply rpcSchema.HostInfoReply
	err := t.rpcClient.Call(ctx, "HostInfo", hostArgs, &hostReply)
	if err != nil {
		return err
	}
	args := rpcSchema.CPUAlarmQueryArgs{
		StartTime: time.Now().Add(-(time.Duration(threshold.Duration)) * time.Minute).UnixMicro(),
		EndTime:   time.Now().UnixMicro(),
	}
	var reply rpcSchema.CPUAlarmQueryReply
	err = t.rpcClient.Call(ctx, "CPUAlarmQuery", args, &reply)
	if err != nil {
		return err
	}
	total := 0.0
	for _, item := range reply.Data {
		total += item.Value
	}
	// Duration 时间段内，如果 CPU 使用率的的平均值大于告警阈值，则告警
	if int(utils.Decimal(total)*100) > threshold.Threshold {

		err := t.db.RunInTransaction(func(tx *gorm.DB) error {
			msg := fmt.Sprintf("%s CPU 使用率连续 %d 分钟超过 %d", hostReply.Hostname, threshold.Duration, threshold.Threshold)
			// 存储系统日志
			operateLog := model.Audit{
				Username: "system",
				Operate:  msg,
			}
			if err := tx.Model(&model.Audit{}).Create(&operateLog).Error; err != nil {
				return err
			}
			// 发送告警邮件
			if err := t.sendMail(msg); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) MemoryAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}
	hostArgs := rpcSchema.HostInfoArgs{}
	var hostReply rpcSchema.HostInfoReply
	err := t.rpcClient.Call(ctx, "HostInfo", hostArgs, &hostReply)
	if err != nil {
		return err
	}
	args := rpcSchema.MemoryAlarmQueryArgs{
		StartTime: time.Now().Add(-(time.Duration(threshold.Duration)) * time.Minute).UnixMicro(),
		EndTime:   time.Now().UnixMicro(),
	}
	var reply rpcSchema.MemoryAlarmQueryReply
	err = t.rpcClient.Call(ctx, "MemoryAlarmQuery", args, &reply)
	if err != nil {
		return err
	}
	total := 0.0
	for _, item := range reply.Data {
		total += item.Value
	}
	if int(utils.Decimal(total)*100) > threshold.Threshold {
		err := t.db.RunInTransaction(func(tx *gorm.DB) error {
			msg := fmt.Sprintf("%s 内存使用率连续 %d 分钟 %d", hostReply.Hostname, threshold.Duration, threshold.Threshold)
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
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) DiskAlarmTask(ctx context.Context) error {
	var threshold model.AlarmThreshold
	if err := t.db.Model(&model.AlarmThreshold{}).First(&threshold).Error; err != nil {
		return err
	}
	hostArgs := rpcSchema.HostInfoArgs{}
	var hostReply rpcSchema.HostInfoReply
	err := t.rpcClient.Call(ctx, "HostInfo", hostArgs, &hostReply)
	if err != nil {
		return err
	}
	args := rpcSchema.DiskAlarmQueryArgs{
		StartTime: time.Now().Add(-(time.Duration(threshold.Duration)) * time.Minute).UnixMicro(),
		EndTime:   time.Now().UnixMicro(),
	}
	var reply rpcSchema.DiskAlarmQueryReply
	err = t.rpcClient.Call(ctx, "DiskAlarmQuery", args, &reply)
	if err != nil {
		return err
	}
	var diskMap = make(map[string]struct{})
	for _, item := range reply.Data {
		if _, ok := diskMap[item.Device]; ok {
			continue
		}
		diskMap[item.Device] = struct{}{}
		if int(utils.Decimal(item.DiskPercent)*100) > threshold.Threshold {
			err := t.db.RunInTransaction(func(tx *gorm.DB) error {
				msg := fmt.Sprintf("%s 磁盘 %s 使用率超过 %d", hostReply.Hostname, item.Device, threshold.Threshold)
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
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Task) ServiceTask(ctx context.Context) error {
	args := rpcSchema.ServiceAlarmQueryArgs{}
	var reply rpcSchema.ServiceAlarmQueryReply
	err := t.rpcClient.Call(ctx, "ServiceAlarmQuery", args, &reply)
	if err != nil {
		return err
	}
	for _, item := range reply.Data {
		if containerStateBytes, ok := t.cache.Get(item.Name); ok {
			if containerStateBytes.(string) != item.State {
				err := t.db.RunInTransaction(func(tx *gorm.DB) error {
					msg := fmt.Sprintf("容器 %s 的状态由 %s 变为 %s", item.Name, containerStateBytes.(string), item.State)
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
				if err != nil {
					return err
				}
			}
		}
		t.cache.Set(item.Name, item.State, 0)
	}
	return nil
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
