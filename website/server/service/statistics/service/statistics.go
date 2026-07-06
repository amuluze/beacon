// Package service
// Date: 2025/02/12 15:25:53
// Author: Amu
// Description:
package service

import (
	"context"
	"github.com/google/wire"
	"server/service/schema"
	"server/service/statistics/repository"
)

var StatisticsServiceSet = wire.NewSet(NewStatisticsService, wire.Bind(new(IStatisticsService), new(*StatisticsService)))

var _ IStatisticsService = (*StatisticsService)(nil)

type IStatisticsService interface {
	StatisticsQuery(context.Context) (schema.StatisticsQueryReply, error)
	StatisticsUpdate(context.Context, schema.StatisticsUpdateArgs) (schema.StatisticsUpdateReply, error)
	InstallationReport(context.Context, schema.InstallationReportArgs) (schema.InstallationReportReply, error)
}

type StatisticsService struct {
	StatisticsRepository repository.IStatisticsRepository
}

func NewStatisticsService(repo repository.IStatisticsRepository) *StatisticsService {
	return &StatisticsService{StatisticsRepository: repo}
}

func (s *StatisticsService) StatisticsQuery(ctx context.Context) (schema.StatisticsQueryReply, error) {
	result := schema.StatisticsQueryReply{}
	reply, err := s.StatisticsRepository.StatisticsQuery(ctx)
	if err != nil {
		return result, err
	}
	result.Data.ID = reply.ID
	result.Data.Times = reply.Times
	return result, nil
}

func (s *StatisticsService) StatisticsUpdate(ctx context.Context, args schema.StatisticsUpdateArgs) (schema.StatisticsUpdateReply, error) {
	result := schema.StatisticsUpdateReply{}
	return result, s.StatisticsRepository.StatisticsUpdate(ctx, args)
}

func (s *StatisticsService) InstallationReport(ctx context.Context, args schema.InstallationReportArgs) (schema.InstallationReportReply, error) {
	result := schema.InstallationReportReply{}
	return result, s.StatisticsRepository.InstallationReport(ctx, args)
}
