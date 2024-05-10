package usecase

import (
	"go-knowledge/libs/golang/resources/go-rabbitmq/queue"
	"go-knowledge/services/eda/event-listener/internal/application/dto/input"
	"go-knowledge/services/eda/event-listener/internal/entity"
	repository "go-knowledge/services/eda/event-listener/internal/infra/database"
	"log"
	"sync"
	"time"
)

type OutputsListenerUseCase struct {
	eventMessageRepository repository.EventMessageRepositoryInterface
	eventsNotifier         *queue.RabbitMQNotifier
}

func NewOutputsListenerUseCase(
	eventMessageRepository repository.EventMessageRepositoryInterface,
	eventsNotifier *queue.RabbitMQNotifier,
) *OutputsListenerUseCase {
	return &OutputsListenerUseCase{
		eventMessageRepository: eventMessageRepository,
		eventsNotifier:         eventsNotifier,
	}
}

func handleCreateOutput(data string, respCh chan Response, wg *sync.WaitGroup) {
	log.Printf("Creating output: %s", data)
	time.Sleep(200 * time.Millisecond)
	respCh <- Response{data: "output created", err: nil}
	wg.Done()
}

func handleUpdateProcessingWallet(data string, respCh chan Response, wg *sync.WaitGroup) {
	log.Printf("Updating processing wallet: %s", data)
	time.Sleep(100 * time.Millisecond)
	respCh <- Response{data: "processing wallet updated", err: nil}
	wg.Done()
}

type Response struct {
	data any
	err  error
}

func (u *OutputsListenerUseCase) Execute(msg input.MessageDTO) (string, error) {
	mockData := "mock data"
	msgEntity, err := entity.NewEventMessage(msg.Data)
	if err != nil {
		log.Printf("Error creating event message: %v", err)
		return "", err
	}
	// store in database in memory
	err = u.eventMessageRepository.CreateEventMessage(msgEntity)
	if err != nil {
		log.Printf("Error storing event message: %v", err)
		return "", err
	}

	var (
		respCh = make(chan Response, 2)
		wg     = &sync.WaitGroup{}
	)
	go handleCreateOutput(mockData, respCh, wg)
	go handleUpdateProcessingWallet(mockData, respCh, wg)
	wg.Add(2)
	wg.Wait()
	close(respCh)

	for resp := range respCh {
		if resp.err != nil {
			log.Printf("Error handling output: %v", resp.err)
			return "", resp.err
		}
		log.Printf("Response: %v", resp.data)
	}

	log.Printf("Received message: %+v", msgEntity)
	return "", nil
}
