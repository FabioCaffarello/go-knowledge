package inmemorydatabase

import (
	"encoding/json"
	"go-knowledge/services/eda/model-order-listener/internal/entity"
)

func mapToModelOrder(document map[string]interface{}, documentEntity *entity.ModelOrder) error {
	documentBytes, err := json.Marshal(document)
    if err != nil {
        return err
    }
    err = json.Unmarshal(documentBytes, &documentEntity)
    if err != nil {
        return err
    }
    return nil
}
