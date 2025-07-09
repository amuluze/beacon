// Package repository
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package repository

import (
	"context"
	"crypto/tls"

	"amprobe/pkg/database"
	"amprobe/service/model"
	"amprobe/service/schema"

	"github.com/google/wire"
	"gopkg.in/gomail.v2"
)

var MailRepositorySet = wire.NewSet(NewMailRepository, wire.Bind(new(IMailRepository), new(*MailRepository)))

var _ IMailRepository = (*MailRepository)(nil)

type IMailRepository interface {
	MailQuery(context.Context) (model.Mail, error)
	MailCreate(context.Context, schema.MailCreateArgs) error
	MailUpdate(context.Context, schema.MailUpdateArgs) error
	MailDelete(context.Context, schema.MailDeleteArgs) error
	MailTest(context.Context, schema.MailTestArgs) error
}

type MailRepository struct {
	DB *database.DB
}

func NewMailRepository(db *database.DB) *MailRepository {
	return &MailRepository{DB: db}
}

func (m *MailRepository) MailQuery(ctx context.Context) (model.Mail, error) {
	var mail model.Mail
	err := m.DB.Model(&model.Mail{}).First(&mail).Error
	return mail, err
}

func (m *MailRepository) MailCreate(ctx context.Context, args schema.MailCreateArgs) error {
	mail := model.Mail{
		Server:   args.Server,
		Port:     args.Port,
		Sender:   args.Sender,
		Password: args.Password,
		Receiver: args.Receiver,
	}
	return m.DB.Model(&model.Mail{}).Create(&mail).Error
}

func (m *MailRepository) MailUpdate(ctx context.Context, args schema.MailUpdateArgs) error {
	params := make(map[string]interface{})
	if args.Server != "" {
		params["server"] = args.Server
	}
	if args.Port != 0 {
		params["port"] = args.Port
	}
	if args.Sender != "" {
		params["sender"] = args.Sender
	}
	if args.Password != "" {
		params["password"] = args.Password
	}
	if args.Receiver != "" {
		params["receiver"] = args.Receiver
	}
	return m.DB.Model(&model.Mail{}).Where("id = ?", args.ID).Updates(params).Error
}

func (m *MailRepository) MailDelete(ctx context.Context, args schema.MailDeleteArgs) error {
	return m.DB.Model(&model.Mail{}).Where("id = ?", args.ID).Delete(&model.Mail{}).Error
}

func (m *MailRepository) MailTest(ctx context.Context, args schema.MailTestArgs) error {
	var mail model.Mail
	err := m.DB.Model(&model.Mail{}).First(&mail).Error
	if err != nil {
		return err
	}
	content := `邮件告警测试消息`
	message := gomail.NewMessage()
	// 增加发件人别名
	message.SetHeader("From", mail.Sender) // 发件人
	message.SetHeader("To", args.Receiver) // 收件人
	message.SetHeader("Subject", "服务器告警")  //定义的主题
	message.SetBody("text/plain", content) //发送的文本消息

	dialer := gomail.NewDialer(mail.Server, mail.Port, mail.Sender, mail.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(message)
}
