package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Plot IDs
var plotIDs = []string{
	"84c7857a-93f6-4896-ab97-29fcd875157e",
	"bd7196c6-5944-4245-a0cd-8dbb9f46191e",
	"ab8f6eb6-79a4-4f9d-8be5-ef4129d5fc2c",
	"de6e212f-7886-45f2-bbc4-cba7c44674ce",
	"06821ed7-53d3-4419-b4a9-f01243f681ac",
	"6eed81a1-bdea-4161-9233-2d9cd728965d",
	"76b31c5b-19bb-4995-833b-35b6fdf87a69",
	"334acf13-c2bd-4752-b930-a1332d45e55c",
	"36a1a12f-e600-47e3-b36b-4c718bb707c9",
}

// Animal IDs
var animalIDs = []string{
	"4dc30f2a-8ea4-4f84-b4ea-58bdc8d5cb41",
	"0a03bff4-f9a8-4f8d-8297-ca92cc740f8d",
	"11f40230-c605-4eb0-ba28-9928f4676ff3",
	"65c52a49-453f-4b7c-92c5-6cbc6ddcbdfb",
	"6a540ff9-3b78-4915-bbec-4487380bb006",
	"c2c3d4aa-ef64-45bf-b0c1-7259cd826b10",
}

// Headers
var headers = map[string]string{
	"accept":             "application/json, text/plain, */*",
	"accept-language":    "vi,vi-VN;q=0.9,fr-FR;q=0.8,fr;q=0.7,en-US;q=0.6,en;q=0.5",
	"authorization":      "tma query_id=AAHxAvVhAgAAAPEC9WH-NwCF&user=%7B%22id%22%3A5938414321%2C%22first_name%22%3A%22Viet%22%2C%22last_name%22%3A%22Hoang%22%2C%22username%22%3A%22anhviet1312%22%2C%22language_code%22%3A%22vi%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F-2xihp_MFnFSH6TKZdZ3Nj82iZYhRRvGjU81vnswwa1UZrJlKJmuKbO-nhDjad3E.svg%22%7D&auth_date=1757477764&signature=e40IXhJc6CMv_8usvVVGZd9gorMlpVZcQJV5GTZnMIBBCOH3G81KPtw_5GeJDh-S2vmGHlAE55X1v440x3q2BA&hash=5bd0b723ec80032e79d6bd3b4b030ca6c2379d4a490a1f8ec3f5768e37d0209c",
	"content-type":       "application/json",
	"origin":             "https://clf.greenlandgame.xyz",
	"priority":           "u=1, i",
	"referer":            "https://clf.greenlandgame.xyz/",
	"sec-ch-ua":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": `"Linux"`,
	"sec-fetch-dest":     "empty",
	"sec-fetch-mode":     "cors",
	"sec-fetch-site":     "same-site",
	"user-agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

// HTTP client
var client = &http.Client{
	Timeout: 30 * time.Second,
}

func makeRequest(method, url string, data []byte) *http.Response {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating request for %s %s: %v", method, url, err)
		return nil
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error in %s %s: %v", method, url, err)
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("%s %s: Status %d, Response: %s", method, url, resp.StatusCode, string(body))
	return resp
}

func sowSeeds() {
	log.Println("Starting sowing process...")
	for _, plotID := range plotIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/farming/plot/%s/sowing", plotID)
		data, _ := json.Marshal(map[string]int{"seed_id": 4})
		makeRequest("Sowing", url, data)
	}
}

func claimRewards() {
	log.Println("Starting claim process for plots...")
	for _, plotID := range plotIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/farming/plot/%s/claim", plotID)
		makeRequest("Claiming plot", url, nil)
	}
}

func feedAnimals() {
	log.Println("Starting feeding process for animals...")
	for _, animalID := range animalIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/animal/feed/%s", animalID)
		makeRequest("Feeding animal", url, nil)
	}
}

func claimAnimals() {
	log.Println("Starting claiming process for animals...")
	for _, animalID := range animalIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/animal/claim/%s", animalID)
		makeRequest("Claiming animal", url, nil)
	}
}

func runSchedule() {
	ict := time.FixedZone("ICT", 7*3600)
	log.Printf("Starting new plot cycle at %s", time.Now().In(ict).Format("2006-01-02 15:04:05 -0700"))
	sowSeeds()
	log.Println("Waiting 10 minutes before claiming plots...")
	time.Sleep(10 * time.Minute)
	claimRewards()
	log.Println("Plot process completed.")
}

func runAnimalSchedule() {
	for {
		feedAnimals()
		log.Println("Waiting 30 minutes before claiming animals...")
		time.Sleep(30 * time.Minute)
		claimAnimals()
		log.Println("Animal process completed, waiting 1 second before next cycle...")
		time.Sleep(1 * time.Second)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/kaihealthcheck" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "Ok")
}

func init() {
	log.SetFlags(0)
}

func main() {
	go func() {
		http.HandleFunc("/kaihealthcheck", healthCheckHandler)
		log.Println("Starting health check server on :8082...")
		if err := http.ListenAndServe(":8082", nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// go runAnimalSchedule()
	for {
		runSchedule()
		log.Println("Waiting 1 second before next plot cycle...")
		time.Sleep(1 * time.Second)
	}
}
