package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	createForecastURL()
	printResult(getCurrentForecast(ForecastURL))
	createLastTenDaysURL()
	printResult(getLasttTenDaysForecast(LastTenDaysURL))
	createHistoricalURLTemplate()
	printResult(getHistoricalForecast(HistoricalURLTemplate, "2023-01-10", "2023-01-12"))
}

var ForecastURL string
var LastTenDaysURL string
var HistoricalURLTemplate string

var latitude string = "47.497913"
var longitude string = "19.040236"

const startDateMarker string = "<startDate>"
const endDateMarker string = "<endDate>"

func createForecastURL() {
	ForecastURL += "https://api.open-meteo.com/v1/forecast"
	ForecastURL += "?" + "latitude=" + latitude + "&longitude=" + longitude
	ForecastURL += "&current_weather=true"
}

func createLastTenDaysURL() {
	LastTenDaysURL += "https://api.open-meteo.com/v1/forecast"
	LastTenDaysURL += "?" + "latitude=" + latitude + "&longitude=" + longitude
	LastTenDaysURL += "&timezone=Europe/Budapest"
	LastTenDaysURL += "&past_days=10"
	LastTenDaysURL += "&daily=temperature_2m_min,temperature_2m_max"
}

func createHistoricalURLTemplate() {
	HistoricalURLTemplate += "https://archive-api.open-meteo.com/v1/era5"
	HistoricalURLTemplate += "?" + "latitude=" + latitude + "&longitude=" + longitude
	HistoricalURLTemplate += "&timezone=Europe/Budapest"
	HistoricalURLTemplate += "&start_date=" + startDateMarker
	HistoricalURLTemplate += "&end_date=" + endDateMarker
	HistoricalURLTemplate += "&daily=temperature_2m_min,temperature_2m_max"
}

func getCurrentForecast(requestURL string) (bool, string, []byte) {
	log.Println("Current forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)

	hasError := hasErrorInResponseBody(formattedBodyBytes)

	return hasError, "Current forecast", formattedBodyBytes
}

func getLasttTenDaysForecast(requestURL string) (bool, string, []byte) {
	log.Println("Last ten days forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)

	hasError := hasErrorInResponseBody(formattedBodyBytes)

	return hasError, "Last ten days forecast", formattedBodyBytes
}

func getHistoricalForecast(requestURL string, startDate string, endDate string) (bool, string, []byte) {
	requestURL = strings.Replace(requestURL, startDateMarker, startDate, 1)
	requestURL = strings.Replace(requestURL, endDateMarker, endDate, 1)
	log.Println("Forecast history URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)

	hasError := hasErrorInResponseBody(formattedBodyBytes)

	return hasError, "Historical forecast", formattedBodyBytes
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

func hasErrorInResponseBody(formattedBodyBytes []byte) bool {
	var result map[string]any
	json.Unmarshal([]byte(formattedBodyBytes), &result)
	_, hasError := result["error"]
	return hasError
}

func printResult(hasError bool, message string, result []byte) {
	errorMessage := ""
	if hasError {
		errorMessage = " error"
	}
	fmt.Printf("\n%s%s:\n%s\n\n", message, errorMessage, result)
}
