// Package service
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package service

import (
	"amprobe/service/mail/repository"
	"amprobe/service/schema"
	"context"
	"github.com/google/wire"
)

var MailServiceSet = wire.NewSet(NewMailService, wire.Bind(new(IMailService), new(*MailService)))

var _ IMailService = (*MailService)(nil)

type IMailService interface {
	MailQuery(context.Context) (schema.Mail, error)
	MailCreate(context.Context, schema.MailCreateArgs) error
	MailUpdate(context.Context, schema.MailUpdateArgs) error
	MailDelete(context.Context, schema.MailDeleteArgs) error
	MailTest(context.Context, schema.MailTestArgs) error
}

type MailService struct {
	MailRepository repository.IMailRepository
}

func NewMailService(mailRepository repository.IMailRepository) *MailService {
	return &MailService{MailRepository: mailRepository}
}

func (m *MailService) MailQuery(ctx context.Context) (schema.Mail, error) {
	result := schema.Mail{}
	reply, err := m.MailRepository.MailQuery(ctx)
	if err != nil {
		return result, err
	}
	result.ID = reply.ID
	result.Server = reply.Server
	result.Port = reply.Port
	result.Sender = reply.Sender
	result.Receiver = reply.Receiver
	return result, nil
}

func (m *MailService) MailCreate(ctx context.Context, args schema.MailCreateArgs) error {
	return m.MailRepository.MailCreate(ctx, args)
}

func (m *MailService) MailUpdate(ctx context.Context, args schema.MailUpdateArgs) error {
	return m.MailRepository.MailUpdate(ctx, args)
}

func (m *MailService) MailDelete(ctx context.Context, args schema.MailDeleteArgs) error {
	return m.MailRepository.MailDelete(ctx, args)
}

func (m *MailService) MailTest(ctx context.Context, args schema.MailTestArgs) error {
	return m.MailRepository.MailTest(ctx, args)
}
