package database

import "go-knowledge/services/eda/model-order-listener/internal/entity"

type ModelOrderRepositoryInterface interface {
	Create(eventMessage *entity.ModelOrder) error
	FindAll() ([]*entity.ModelOrder, error)
	FindByID(id string) (*entity.ModelOrder, error)
    DeleteByID(id string) error
}
