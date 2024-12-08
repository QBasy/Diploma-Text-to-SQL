package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func startService(serviceName, cmd string) error {
	fmt.Printf("Starting %s service...\n", serviceName)
	command := exec.Command(cmd)
	err := command.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %v", serviceName, err)
	}
	fmt.Printf("%s service started.\n", serviceName)
	return nil
}

func main() {

	if err := startService("Database Service", "go build ./database-service/database-service"); err != nil {
		log.Fatalf("Error starting Database Service: %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := startService("Text-to-SQL Service", "python3 ./text-to-sql-service/main.py"); err != nil {
		log.Fatalf("Error starting Text-to-SQL service: %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := startService("API", "go build ./API/API"); err != nil {
		log.Fatalf("Error starting API service: %v", err)
	}

	select {}
}
