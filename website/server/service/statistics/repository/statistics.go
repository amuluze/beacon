// Package repository
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"server/pkg/database"
	"server/service/model"
	"server/service/schema"

	"github.com/google/wire"
)

var StatisticsRepositorySet = wire.NewSet(NewStatisticsRepository, wire.Bind(new(IStatisticsRepository), new(*StatisticsRepository)))

var _ IStatisticsRepository = (*StatisticsRepository)(nil)

type IStatisticsRepository interface {
	StatisticsQuery(context.Context) (model.Statistics, error)
	StatisticsUpdate(context.Context, schema.StatisticsUpdateArgs) error
	InstallationReport(context.Context, schema.InstallationReportArgs) error
}

type StatisticsRepository struct {
	DB *database.DB
}

func NewStatisticsRepository(db *database.DB) *StatisticsRepository {
	return &StatisticsRepository{DB: db}
}

func (s *StatisticsRepository) StatisticsQuery(ctx context.Context) (model.Statistics, error) {
	var statistics model.Statistics
	if err := s.DB.Model(&model.Statistics{}).First(&statistics).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statistics = model.Statistics{
				Times: 0,
			}
			s.DB.Model(&model.Statistics{}).Create(&statistics)
		}
	}
	return statistics, nil
}

func (s *StatisticsRepository) StatisticsUpdate(ctx context.Context, args schema.StatisticsUpdateArgs) error {
	// 更新 times 字段 + 1
	if err := s.DB.Model(&model.Statistics{}).Where("id = ?", args.ID).UpdateColumn("times", gorm.Expr("times + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func (s *StatisticsRepository) InstallationReport(ctx context.Context, args schema.InstallationReportArgs) error {
	report := model.InstallationReport{
		InstallID:     args.InstallID,
		Image:         args.Image,
		Version:       args.Version,
		PublicBaseURL: args.PublicBaseURL,
		InstallDir:    args.InstallDir,
		HTTPPort:      args.HTTPPort,
		ControlPort:   args.ControlPort,
		ContainerName: args.ContainerName,
		Hostname:      args.Hostname,
		ClientIP:      args.ClientIP,
		UserAgent:     args.UserAgent,
	}
	return s.DB.WithContext(ctx).Create(&report).Error
}
