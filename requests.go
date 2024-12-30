package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Data []ResponseData `json:"data"`
	NextPage
}

type ResponseData struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type NextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	Uri    string `json:"uri"`
}

func GetProjects() error {
	url := fmt.Sprintf("https://app.asana.com/api/1.0/workspaces/%v/projects", os.Getenv("WORKSPACE_GID"))
	ticker := time.NewTicker(time.Minute * 5)

	for range ticker.C {
		projects, res, err := GetRequest(url)

		if res != nil && res.StatusCode == http.StatusTooManyRequests {
			retryAfter := res.Header.Get("Retry-After")

			dur, err := time.ParseDuration(retryAfter + "s")
			if err != nil {
				fmt.Println("Error parsing duration")
			}
			time.Sleep(dur)
		} else if err != nil {
			fmt.Println("Error getting projects: ", projects)
			return err
		} else {
			SaveDataToFile(projects, "projects")
		}
	}

	return nil
}

func GetUsers() error {
	url := "https://app.asana.com/api/1.0/users"
	ticker := time.NewTicker(time.Minute * 5)

	for range ticker.C {
		users, res, err := GetRequest(url)

		if res != nil && res.StatusCode == http.StatusTooManyRequests {
			retryAfter := res.Header.Get("Retry-After")

			dur, err := time.ParseDuration(retryAfter + "s")
			if err != nil {
				fmt.Println("Error parsing duration")
			}
			time.Sleep(dur)
		} else if err != nil {
			fmt.Println("Error getting users: ", users)
			return err
		} else {
			SaveDataToFile(users, "users")
		}
	}

	return nil
}

func GetRequest(url string) (Response, *http.Response, error) {
	var data Response
	var res *http.Response

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating http request: ", err)
		return data, res, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error getting data: ", err)
		return data, res, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding request body: ", err)
		return data, res, err
	}

	fmt.Printf("Data: %+v \n", data)

	return data, res, nil
}

func SaveDataToFile(data Response, entity string) {
	var fileName string
	if entity == "users" {
		fileName = "users.txt"
	} else if entity == "projects" {
		fileName = "projects.txt"
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not open %v", fileName)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)

	_, err = file.WriteString(string(jsonData))
	if err != nil {
		fmt.Printf("Could not write text %v \n", fileName)
	} else {
		fmt.Printf("Operation successful! Text has been appended %v \n", fileName)
	}
}
