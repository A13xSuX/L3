package service

import (
	"context"
	"l3/WarehouseControl/internal/models"
	"l3/WarehouseControl/internal/repository"
)

type AuditService struct {
	AuditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{
		AuditRepo: auditRepo,
	}
}

func (s *AuditService) GetByItemID(ctx context.Context, itemID int64) ([]models.HistoryItemResponse, error) {
	audit, err := s.AuditRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	auditResp := make([]models.HistoryItemResponse, len(audit))
	for i, row := range audit {
		rowResp := models.HistoryItemResponse{
			ID:                row.ID,
			ItemID:            row.ItemID,
			Action:            row.Action,
			OldData:           row.OldData,
			NewData:           row.NewData,
			ChangedByUsername: row.ChangedByUsername,
			ChangedByRole:     row.ChangedByRole,
			ChangedAt:         row.ChangedAt,
		}
		auditResp[i] = rowResp
	}
	return auditResp, nil
}
