package main

import (
	"log"
	"todolist/internal/server"
)

func main() {
	log.Println("Starting TodoList server...")
	
	s := server.NewServer()
	if err := s.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}