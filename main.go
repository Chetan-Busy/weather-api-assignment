package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type WeatherData struct {
	Data []struct {
		Name    string   `json:"name"`
		Status  []string `json:"status"`
		Weather string   `json:"weather"`
	} `json:"data"`
	// TotalPages int `json:"total_pages"`
}

type Result struct {
	TotalPages uint `json:"total_pages"`
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`\d+`)
	num := re.FindString(s)
	return num
}

var result [][]string

func ExtractData(name string, pageNumber int) {
	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/weather/search?name=%s&page=%v", name, int(pageNumber))

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

	json.Unmarshal([]byte(formattedJSON), &weatherData)

	for _, item := range weatherData.Data {
		weather := extractNumber(item.Weather)
		wind := extractNumber(item.Status[0])
		humidity := extractNumber(item.Status[1])
		row := []string{item.Name, weather, wind, humidity}
		result = append(result, row)
	}

}

func main() {
	fmt.Println("Weather Api Assignment")

	var name string
	fmt.Scanln(&name)

	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/weather/search?name=%s", name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	r := Result{}
	jsonErr := json.Unmarshal(body, &r)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	fmt.Println(r.TotalPages)

	for i := 1; i <= int(r.TotalPages); i++ {
		ExtractData(name, i)
	}

	fmt.Println(result)

}
