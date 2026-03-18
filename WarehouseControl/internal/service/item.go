package service

import (
	"context"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"
	"l3/WarehouseControl/internal/repository"
)

type ItemService struct {
	ItemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{
		ItemRepo: itemRepo,
	}
}

func (s *ItemService) Create(ctx context.Context, itemReq *models.CreateItemRequest) (*models.ItemResponse, error) {
	if itemReq.Title == "" {
		return nil, customErrs.ErrTitleEmpty
	}
	if itemReq.Sku == "" {
		return nil, customErrs.ErrSkuEmpty
	}
	if itemReq.Quantity < 0 {
		return nil, customErrs.ErrQuantityNotPositive
	}

	item := models.Item{
		Title:    itemReq.Title,
		Sku:      itemReq.Sku,
		Quantity: itemReq.Quantity,
	}

	err := s.ItemRepo.Create(ctx, &item)
	if err != nil {
		return nil, err
	}

	resp := models.ItemResponse{
		ID:        item.ID,
		Title:     item.Title,
		Sku:       item.Sku,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return &resp, err
}

func (s *ItemService) GetAll(ctx context.Context) ([]models.ItemResponse, error) {
	items, err := s.ItemRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	itemsResp := make([]models.ItemResponse, len(items))
	for i, item := range items {
		itemResp := models.ItemResponse{
			ID:        item.ID,
			Title:     item.Title,
			Sku:       item.Sku,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		itemsResp[i] = itemResp
	}
	return itemsResp, nil
}

func (s *ItemService) Update(ctx context.Context, id int64, itemReq *models.UpdateItemRequest) error {
	if itemReq.Title == "" {
		return customErrs.ErrTitleEmpty
	}
	if itemReq.Sku == "" {
		return customErrs.ErrSkuEmpty
	}
	if itemReq.Quantity < 0 {
		return customErrs.ErrQuantityNotPositive
	}
	item := &models.Item{
		Title:    itemReq.Title,
		Sku:      itemReq.Sku,
		Quantity: itemReq.Quantity,
	}
	err := s.ItemRepo.Update(ctx, id, item)
	return err
}

func (s *ItemService) Delete(ctx context.Context, id int64) error {
	return s.ItemRepo.Delete(ctx, id)
}
