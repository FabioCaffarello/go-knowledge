package main

import (
	"fmt"
	"time"
)

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"
// )

// func FetchUserID(ctx context.Context) (string, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
// 	defer cancel()

// 	val := ctx.Value("requestID").(string)
// 	fmt.Printf("Processing request %v\n", val)

// 	type result struct {
// 		userID string
// 		err    error
// 	}

// 	resultCh := make(chan result, 1)

// 	go func() {
// 		res, err := thirdPartyHTTPCall()
// 		resultCh <- result{res, err}
// 	}()

// 	select {
// 	case <-ctx.Done():
// 		return "", ctx.Err()
// 	case res := <-resultCh:
// 		if res.err != nil {
// 			return "", res.err
// 		}
// 		return res.userID, nil
// 	}
// }

// func thirdPartyHTTPCall() (string, error) {
// 	time.Sleep(50 * time.Millisecond)
// 	return "User ID 1", nil
// }

// func main() {
// 	start := time.Now()
// 	ctx := context.WithValue(context.Background(), "requestID", "user123")
// 	userID, err := FetchUserID(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to fetch user ID: %v", err)
// 	}
// 	fmt.Printf("The response took %v: %v", time.Since(start), userID)
// }

type Server struct {
	quitCh chan struct{}
	msgCh  chan string
}

func NewServer() *Server {
	return &Server{
		quitCh: make(chan struct{}),
		msgCh:  make(chan string, 128),
	}
}

func (s *Server) sendMessage(msg string) {
	s.msgCh <- msg
}

func (s *Server) Start() {
	fmt.Println("Server starting...")
	s.loop()
}

func (s *Server) loop() {
mainloop:
	for {
		select {
		case <-s.quitCh:
			fmt.Println("Quitting server...")
			break mainloop
		case msg := <-s.msgCh:
			s.handleMessage(msg)
		}
	}
	fmt.Println("Server shutting down gracefully")
}

func (s *Server) handleMessage(msg string) {
	fmt.Printf("Handling message: %v\n", msg)
}

func (s *Server) quit() {
	s.quitCh <- struct{}{}
}

func main() {
	server := NewServer()

	go func() {
		for i := 0; i < 10; i++ {
			server.sendMessage(fmt.Sprintf("Message %v", i))
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		server.quit()
	}()

	server.Start()
}
