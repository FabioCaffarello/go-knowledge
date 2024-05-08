package server



// type Server struct {
// 	quitCh         chan struct{}
// 	msgCh          chan amqp.Delivery
// 	tempCh         chan string
// 	rabbitMQ       *queue.RabbitMQ
// 	eventsNotifier *queue.RabbitMQNotifier
// 	eventsConsumer *queue.RabbitMQConsumer
// }

// func NewServer(rmq *queue.RabbitMQ, notifier *queue.RabbitMQNotifier, consumer *queue.RabbitMQConsumer) *Server {
//     return &Server{
//         rabbitMQ:       rmq,
//         eventsNotifier: notifier,
//         eventsConsumer: consumer,
//         quitCh:         make(chan struct{}),
//         msgCh:          make(chan amqp.Delivery, 128),
//         tempCh:         make(chan string, 128),
//     }
// }

// func (s *Server) Start() {
// 	fmt.Println("Server starting...")
// 	go s.consumeRabbitMQ()
// 	s.loop()
// }

// func (s *Server) consumeRabbitMQ() {
// 	go s.eventsConsumer.Consume(s.msgCh)
// 	for {
// 		select {
// 		case msg, ok := <-s.msgCh:
// 			if !ok {
// 				log.Println("Message channel closed, stopping consumption")
// 				return
// 			}
// 			msg.Ack(false)
// 			s.tempCh <- string(msg.Body)
// 		default:
// 		}
// 	}
// }

// func (s *Server) cleanupRabbitMQ() {
// 	// Perform cleanup, e.g., ack messages, close channels, etc.
// 	if s.rabbitMQ != nil {
// 		s.rabbitMQ.Close()
// 	}
// }

// func (s *Server) loop() {
// mainloop:
// 	for {
// 		select {
// 		case <-s.quitCh:
// 			fmt.Println("Quitting server...")
// 			s.cleanupRabbitMQ()
// 			break mainloop
// 		case msg := <-s.tempCh:
// 			s.handleMessage(msg)
// 		}
// 	}
// 	defer s.rabbitMQ.Close()
// 	fmt.Println("Server shutting down gracefully")
// }

// func (s *Server) handleMessage(msg string) {
// 	fmt.Printf("Handling message: %v\n", msg)
// }

// func (s *Server) quit() {
// 	s.quitCh <- struct{}{}
// }