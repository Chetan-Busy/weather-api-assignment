package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
)

type WeatherData struct {
	Data []struct {
		Name    string   `json:"name"`
		Status  []string `json:"status"`
		Weather string   `json:"weather"`
	} `json:"data"`
}

type Result struct {
	TotalPages uint `json:"total_pages"`
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`\d+`)
	num := re.FindString(s)
	return num
}

var resultMutex sync.Mutex
var result [][]string

func ExtractData(name string, pageNumber int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/weather/search?name=%s&page=%v", name, pageNumber)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making a get request", err.Error())
		return
	}
	var data map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	formattedJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("Error formatting JSON:", err)
	}
	var weatherData WeatherData

	jsonErr := json.Unmarshal(formattedJSON, &weatherData)
	if jsonErr != nil {
		fmt.Println("Error while unmarshal operation", jsonErr.Error())
		return
	}

	for _, item := range weatherData.Data {
		weather := extractNumber(item.Weather)
		wind := extractNumber(item.Status[0])
		humidity := extractNumber(item.Status[1])
		row := []string{item.Name, weather, wind, humidity}

		resultMutex.Lock()
		result = append(result, row)
		resultMutex.Unlock()
	}
}

func main() {
	fmt.Println("Weather Api Assignment")

	var name string
	fmt.Scanln(&name)

	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/weather/search?name=%s", name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Error ", readErr.Error())
		return
	}

	r := Result{}
	jsonErr := json.Unmarshal(body, &r)
	if jsonErr != nil {
		fmt.Println("Error", jsonErr.Error())
		return
	}

	fmt.Println(r.TotalPages)

	var wg sync.WaitGroup

	for i := 1; i <= int(r.TotalPages); i++ {
		wg.Add(1)
		go ExtractData(name, i, &wg)
	}

	wg.Wait()

	fmt.Println(result)
}
