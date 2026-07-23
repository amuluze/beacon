package repository

import (
	"context"

	"beacon/service/model"
	"common/database"

	"github.com/google/wire"
)

var DingTalkRepositorySet = wire.NewSet(NewDingTalkRepository, wire.Bind(new(IDingTalkRepository), new(*DingTalkRepository)))

type IDingTalkRepository interface {
	Query(context.Context) (model.DingTalk, error)
	Save(context.Context, *model.DingTalk) error
}

type DingTalkRepository struct {
	DB *database.DB
}

func NewDingTalkRepository(db *database.DB) *DingTalkRepository {
	return &DingTalkRepository{DB: db}
}

func (r *DingTalkRepository) Query(ctx context.Context) (model.DingTalk, error) {
	var setting model.DingTalk
	err := r.DB.WithContext(ctx).
		Where("key = ?", model.DefaultDingTalkConfigKey).
		First(&setting).Error
	return setting, err
}

func (r *DingTalkRepository) Save(ctx context.Context, setting *model.DingTalk) error {
	return r.DB.WithContext(ctx).Save(setting).Error
}
