package usecase

import (
	"encoding/json"
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
	EventDispatcher      events.EventDispatcherInterface
}

func NewCreateModelOrderUseCase(
	modelOrderRepository repository.ModelOrderRepositoryInterface,
	errorCreated events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateModelOrderUseCase {
	return &CreateModelOrderUseCase{
		modelOrderRepository: modelOrderRepository,
		ErrorCreated:         errorCreated,
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
