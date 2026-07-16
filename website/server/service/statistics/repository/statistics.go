// Package repository
// Date:   2024/10/14 16:08
// Author: Amu
// Description:
package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// StatisticsQuery 读取唯一的统计记录；不存在时幂等播种为 Times=0。
// 所有底层错误（连接、SQL）必须向上冒泡，禁止吞错。
func (s *StatisticsRepository) StatisticsQuery(ctx context.Context) (model.Statistics, error) {
	var statistics model.Statistics
	err := s.DB.WithContext(ctx).
		Model(&model.Statistics{}).
		FirstOrCreate(&statistics, model.Statistics{Times: 0}).Error
	if err != nil {
		return model.Statistics{}, err
	}
	return statistics, nil
}

// StatisticsUpdate 将指定记录的 times 自增 1；
// 当目标记录不存在（RowsAffected=0）时返回 ErrRecordNotFound，避免静默失败。
func (s *StatisticsRepository) StatisticsUpdate(ctx context.Context, args schema.StatisticsUpdateArgs) error {
	res := s.DB.WithContext(ctx).
		Model(&model.Statistics{}).
		Where("id = ?", args.ID).
		UpdateColumn("times", gorm.Expr("times + ?", 1))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("statistics record not found")
	}
	return nil
}

// InstallationReport 落库一次安装上报；InstallID 的唯一约束 + OnConflict 幂等，
// 保证安装脚本重试/网络抖动不会产生重复记录。
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
	return s.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&report).Error
}
