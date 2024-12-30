package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading env variables: ", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := GetProjects(); err != nil {
			fmt.Println("Error getting projects: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := GetUsers(); err != nil {
			fmt.Println("Error getting users: ", err)
		}
	}()

	wg.Wait()
}
