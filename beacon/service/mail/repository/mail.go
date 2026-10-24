// Package repository
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package repository

import (
	alarmtitle "beacon/service/alarm/title"
	"beacon/pkg/secret"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
	"context"

	"github.com/google/wire"
	"gopkg.in/gomail.v2"
)

var MailRepositorySet = wire.NewSet(NewMailRepository, wire.Bind(new(IMailRepository), new(*MailRepository)))

var _ IMailRepository = (*MailRepository)(nil)

type IMailRepository interface {
	MailQuery(context.Context) (model.Mail, error)
	MailCreate(context.Context, schema.MailCreateArgs) error
	MailUpdate(context.Context, schema.MailUpdateArgs) error
	MailTest(context.Context, schema.MailTestArgs, string) error
}

type MailRepository struct {
	DB     *database.DB
	Cipher *secret.Cipher
}

func NewMailRepository(db *database.DB, cipher *secret.Cipher) *MailRepository {
	return &MailRepository{DB: db, Cipher: cipher}
}

func (m *MailRepository) MailQuery(ctx context.Context) (model.Mail, error) {
	var mail model.Mail
	err := m.DB.WithContext(ctx).Model(&model.Mail{}).First(&mail).Error
	if err != nil {
		return mail, err
	}
	password, legacy, err := m.Cipher.Decrypt(mail.Password)
	if err != nil {
		return mail, err
	}
	if legacy {
		encrypted, err := m.Cipher.Encrypt(password)
		if err != nil {
			return mail, err
		}
		if err := m.DB.WithContext(ctx).Model(&mail).Update("password", encrypted).Error; err != nil {
			return mail, err
		}
	}
	mail.Password = password
	return mail, err
}

func (m *MailRepository) MailCreate(ctx context.Context, args schema.MailCreateArgs) error {
	password, err := m.Cipher.Encrypt(args.Password)
	if err != nil {
		return err
	}
	mail := model.Mail{
		Server:   args.Server,
		Port:     args.Port,
		Sender:   args.Sender,
		Password: password,
		Receiver: args.Receiver,
	}
	return m.DB.WithContext(ctx).Model(&model.Mail{}).Create(&mail).Error
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
		password, err := m.Cipher.Encrypt(args.Password)
		if err != nil {
			return err
		}
		params["password"] = password
	}
	if args.Receiver != "" {
		params["receiver"] = args.Receiver
	}
	return m.DB.WithContext(ctx).Model(&model.Mail{}).Where("id = ?", args.ID).Updates(params).Error
}

func (m *MailRepository) MailTest(ctx context.Context, args schema.MailTestArgs, title string) error {
	mail, err := m.MailQuery(ctx)
	if err != nil {
		return err
	}
	message := buildMailTestMessage(mail, args, title)
	dialer := gomail.NewDialer(mail.Server, mail.Port, mail.Sender, mail.Password)

	return dialer.DialAndSend(message)
}

func buildMailTestMessage(mail model.Mail, args schema.MailTestArgs, title string) *gomail.Message {
	content := `邮件告警测试消息`
	message := gomail.NewMessage()
	// 增加发件人别名
	message.SetHeader("From", mail.Sender)                    // 发件人
	message.SetHeader("To", args.Receiver)                    // 收件人
	message.SetHeader("Subject", alarmtitle.OrDefault(title)) // 定义的主题
	message.SetBody("text/plain", content)                    //发送的文本消息
	return message
}
