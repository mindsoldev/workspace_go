package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	getCurrentForecast()
	getLasttTenDaysForecast()
	getHistoricalForecast("2023-01-10", "2023-01-12")
}

var forecastURL string = "https://api.open-meteo.com/v1/forecast"
var lastTenDaysURL string = "https://api.open-meteo.com/v1/forecast"
var historicalURL string = "https://archive-api.open-meteo.com/v1/era5"

var latitude string = "47.497913"
var longitude string = "19.040236"

func getCurrentForecast() {
	requestURL := forecastURL + "?" + "latitude=" + latitude + "&longitude=" + longitude
	requestURL += "&current_weather=true"
	log.Println("Current forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes := readBody(response)

		formattedBodyBytes := formatJson(bodyBytes)

		fmt.Printf("\nCurrent forecast:\n%s\n\n", formattedBodyBytes)
	}
}

func getLasttTenDaysForecast() {
	requestURL := lastTenDaysURL + "?" + "latitude=" + latitude + "&longitude=" + longitude
	requestURL += "&timezone=Europe/Budapest"
	requestURL += "&past_days=10"
	requestURL += "&daily=temperature_2m_min,temperature_2m_max"
	log.Println("Last ten days forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes := readBody(response)

		formattedBodyBytes := formatJson(bodyBytes)

		fmt.Printf("\nLast ten days forecast: \n%s\n\n", formattedBodyBytes)
	}
}

func getHistoricalForecast(startDate string, endDate string) {
	requestURL := historicalURL + "?" + "latitude=" + latitude + "&longitude=" + longitude
	requestURL += "&timezone=Europe/Budapest"
	requestURL += "&start_date=" + startDate
	requestURL += "&end_date=" + endDate
	requestURL += "&daily=temperature_2m_min,temperature_2m_max"
	log.Println("Forecast history URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes := readBody(response)

		formattedBodyBytes := formatJson(bodyBytes)

		fmt.Printf("\nHistorical forecast from: %s to %s:\n%s\n\n", startDate, endDate, formattedBodyBytes)
	}
}

func executeQuery(requestURL string) *http.Response {
	response, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func readBody(response *http.Response) []byte {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes
}

func formatJson(jsonBytes []byte) []byte {
	var res bytes.Buffer
	err := json.Indent(&res, jsonBytes, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return res.Bytes()
}
