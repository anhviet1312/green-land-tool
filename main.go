package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"authorization":      "tma query_id=AAHxAvVhAgAAAPEC9WFofORf&user=%7B%22id%22%3A5938414321%2C%22first_name%22%3A%22Viet%22%2C%22last_name%22%3A%22Hoang%22%2C%22username%22%3A%22anhviet1312%22%2C%22language_code%22%3A%22vi%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F-2xihp_MFnFSH6TKZdZ3Nj82iZYhRRvGjU81vnswwa1UZrJlKJmuKbO-nhDjad3E.svg%22%7D&auth_date=1757296594&signature=cMfNKNN_5KCmc_UZ_s9_vFoAj4Ku5aBag-G_ffF1lfqjuNc1g0qQ9shmr8_bCLpErfhPs0ACwi2yxr6B8V-mDw&hash=04ab0edbe682fbc2d47fc1264e9077613a785a902c4bb75b559c631a2e12b9f1",
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

// HTTP client with timeout
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// makeRequest handles HTTP PATCH requests with retries
func makeRequest(method, url string, data []byte) *http.Response {
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("Error creating request for %s %s: %v\n", method, url, err)
			time.Sleep(5 * time.Second)
			continue
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error in %s %s: %v\n", method, url, err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("%s %s: Status %d, Response: %s\n", method, url, resp.StatusCode, string(body))
		return resp
	}
	fmt.Printf("Failed %s %s after 3 attempts\n", method, url)
	return nil
}

// sowSeeds sends sowing requests for all plots
func sowSeeds() {
	fmt.Println("Starting sowing process...")
	for _, plotID := range plotIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/farming/plot/%s/sowing", plotID)
		data, _ := json.Marshal(map[string]int{"seed_id": 4})
		makeRequest("Sowing", url, data)
		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}

// claimRewards sends claim requests for all plots
func claimRewards() {
	fmt.Println("Starting claim process for plots...")
	for _, plotID := range plotIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/farming/plot/%s/claim", plotID)
		makeRequest("Claiming plot", url, nil)
		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}

// feedAnimals sends feed requests for all animals
func feedAnimals() {
	fmt.Println("Starting feeding process for animals...")
	for _, animalID := range animalIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/animal/feed/%s", animalID)
		makeRequest("Feeding animal", url, nil)
		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}

// claimAnimals sends claim requests for all animals
func claimAnimals() {
	fmt.Println("Starting claiming process for animals...")
	for _, animalID := range animalIDs {
		url := fmt.Sprintf("https://service.greenlandgame.xyz/api/v1/animal/claim/%s", animalID)
		makeRequest("Claiming animal", url, nil)
		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}

// runSchedule handles the plot sowing and claiming cycle
func runSchedule() {
	// Print current time in +07:00 timezone
	ict := time.FixedZone("ICT", 7*3600)
	currentTime := time.Now().In(ict).Format("2006-01-02 15:04:05 -0700")
	fmt.Printf("Starting new plot cycle at %s\n", currentTime)

	sowSeeds()
	fmt.Println("Waiting 10 minutes before claiming plots...")
	time.Sleep(10 * time.Minute)
	claimRewards()
	fmt.Println("Plot process completed.")
}

// runAnimalSchedule handles the animal feeding and claiming cycle
func runAnimalSchedule() {
	for {
		feedAnimals()
		fmt.Println("Waiting 30 minutes before claiming animals...")
		time.Sleep(30 * time.Minute)
		claimAnimals()
		fmt.Println("Animal process completed, waiting 1 second before next cycle...")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Start animal schedule in a separate goroutine
	// go runAnimalSchedule()

	// Run plot schedule in main goroutine
	for {
		runSchedule()
		fmt.Println("Waiting 1 second before next plot cycle...")
		time.Sleep(1 * time.Second)
	}
}
