package usecase

import (
	"encoding/json"
	"fmt"
	"go-knowledge/libs/golang/services/shared/go-events/events"
	"go-knowledge/services/eda/model-order-listener/internal/entity"
	repository "go-knowledge/services/eda/model-order-listener/internal/infra/database"
	inputDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/input"
	outputDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/output"
	"log"
)

type CreateModelOrderUseCase struct {
	modelOrderRepository repository.ModelOrderRepositoryInterface
	ErrorCreated         events.EventInterface
    FileOrderCreated     events.EventInterface
	EventDispatcher      events.EventDispatcherInterface
}

func NewCreateModelOrderUseCase(
	modelOrderRepository repository.ModelOrderRepositoryInterface,
	errorCreated events.EventInterface,
    fileOrderCreated     events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateModelOrderUseCase {
	return &CreateModelOrderUseCase{
		modelOrderRepository: modelOrderRepository,
		ErrorCreated:         errorCreated,
        FileOrderCreated:     fileOrderCreated,
		EventDispatcher:      eventDispatcher,
	}
}

func (u *CreateModelOrderUseCase) ProcessMessageChannel(msgCh <-chan []byte, listenerTag string) {
	u.ErrorCreated.SetTag(listenerTag)
	for msg := range msgCh {
		var msgDTO inputDTO.ModelOrderDTO
		err := json.Unmarshal(msg, &msgDTO)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			errMsg := outputDTO.ErrMsgDTO{
				Err:         err,
				Msg:         msg,
				ListenerTag: listenerTag,
			}
			u.ErrorCreated.AddPayload(errMsg)
		}
		log.Printf("Message received: %v", msgDTO)
		err = u.execute(msgDTO)

		if err != nil {
			log.Printf("Error processing message: %v", err)
			errMsg := outputDTO.ErrMsgDTO{
				Err:         err,
				Msg:         msg,
				ListenerTag: listenerTag,
			}
			u.ErrorCreated.AddPayload(errMsg)
		}
	}
	u.EventDispatcher.Dispatch(u.ErrorCreated)
}

func (u *CreateModelOrderUseCase) execute(msg inputDTO.ModelOrderDTO) error {
    files := GetFilesNameFromDTO(msg.Files)
	modelOrder, err := entity.NewModelOrder(
		msg.Costumer,
		msg.Context,
		ConvertSubcontextsDTOToEntity(msg.Subcontexts),
		msg.BucketName,
		files,
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
    for _, fileReference := range modelOrder.FilesReferences {
        fileOrder := outputDTO.FileOrderDTO{
            OrderID:     string(modelOrder.ID),
            ModelID:     string(modelOrder.ModelID),
            FileID:      string(fileReference.FileID),
            Costumer:    modelOrder.Costumer,
            Context:     modelOrder.Context,
            Subcontexts: ConvertSubcontextsEntityToDTO(modelOrder.Subcontexts),
            BucketName:  modelOrder.BucketName,
            Partition:   modelOrder.Partition,
        }
        u.FileOrderCreated.AddPayload(
            fileOrder,
        )
    }
    u.FileOrderCreated.SetTag(fmt.Sprintf("file-order-created:%s:%s", modelOrder.Costumer, modelOrder.ModelID))
    u.EventDispatcher.Dispatch(u.FileOrderCreated)
	return nil
}
