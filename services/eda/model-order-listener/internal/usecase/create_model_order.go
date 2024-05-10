package usecase

import (
	"encoding/json"
	"go-knowledge/services/eda/model-order-listener/internal/entity"
	repository "go-knowledge/services/eda/model-order-listener/internal/infra/database"
	inputDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/input"
	"log"
)

type CreateModelOrderUseCase struct {
	modelOrderRepository repository.ModelOrderRepositoryInterface
}

func NewCreateModelOrderUseCase(
	modelOrderRepository repository.ModelOrderRepositoryInterface,
) *CreateModelOrderUseCase {
	return &CreateModelOrderUseCase{
		modelOrderRepository: modelOrderRepository,
	}
}

func (u *CreateModelOrderUseCase) ProcessMessageChannel(msgCh <-chan []byte, quitCh chan struct{}) {
	for msg := range msgCh {
		var msgDTO inputDTO.ModelOrderDTO
		err := json.Unmarshal(msg, &msgDTO)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
		log.Printf("Message received: %v", msgDTO)
		err = u.execute(msgDTO)
        
        if err != nil {
            log.Printf("Error processing message: %v", err)
        }
	}
}

func (u *CreateModelOrderUseCase) execute(msg inputDTO.ModelOrderDTO) error {
	modelOrder, err := entity.NewModelOrder(
        msg.Costumer,
        msg.Context,
        ConvertSubcontextsDTOToEntity(msg.Subcontexts),
        msg.BucketName,
        ConvertFilesReferencesDTOToEntity(msg.FilesReferences),	
        msg.Partition,
    )
    if err != nil {
        log.Printf("Error creating model order: %v", err)
        return err
    }
    err = u.modelOrderRepository.Create(modelOrder)
    if err != nil {
        log.Printf("Error saving model order: %v", err)
        return err
    }

    return nil   
}
