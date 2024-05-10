package inmemorydatabase

import (
	"go-knowledge/libs/golang/resources/database/in-memory/go-doc-db-client/client"
	"go-knowledge/services/eda/model-order-listener/internal/entity"
	"log"
)

type ModelOrderRepository struct {
	databaseName   string
	client         *client.Client
	collectionName string
}

func NewModelOrderRepository(
	databaseName string,
	client *client.Client,
	collectionName string,
) *ModelOrderRepository {
	repository := &ModelOrderRepository{
		databaseName:   databaseName,
		client:         client,
		collectionName: collectionName,
	}
	repository.client.CreateCollection(collectionName)
	return repository
}

func (r *ModelOrderRepository) Create(modelOrder *entity.ModelOrder) error {
	log.Printf("Save model order in memory")
	modelOrderMap := map[string]interface{}{
		"_id":              modelOrder.ID,
		"model_id":         modelOrder.ModelID,
		"customer":         modelOrder.Costumer,
		"context":          modelOrder.Context,
		"subcontexts":      modelOrder.Subcontexts,
		"bucket_name":      modelOrder.BucketName,
		"files_references": modelOrder.FilesReferences,
		"partition":        modelOrder.Partition,
		"created_at":       modelOrder.CreatedAt,
	}

	err := r.client.InsertOne(r.collectionName, modelOrderMap)
	if err != nil {
		return err
	}
	return nil
}

func (r *ModelOrderRepository) FindAll() ([]*entity.ModelOrder, error) {
	log.Printf("Find all model orders in memory")
	documents, err := r.client.FindAll(r.collectionName)
	if err != nil {
		return nil, err
	}
	modelOrders := make([]*entity.ModelOrder, 0, len(documents))
	for _, document := range documents {
		var result entity.ModelOrder
		if err := mapToModelOrder(document, &result); err != nil {
			return nil, err
		}
		modelOrders = append(modelOrders, &result)
	}
	return modelOrders, nil
}

func (r *ModelOrderRepository) FindByID(id string) (*entity.ModelOrder, error) {
	log.Printf("Find model order by ID in memory")
	document, err := r.client.FindOne(r.collectionName, id)
	if err != nil {
		return nil, err
	}
	var result entity.ModelOrder
	if err := mapToModelOrder(document, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
