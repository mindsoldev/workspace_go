package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const startDateMarker string = "<startDate>"
const endDateMarker string = "<endDate>"

type forecastApi struct {
	ForecastURL           *string
	LastTenDaysURL        *string
	HistoricalURLTemplate *string
	latitude              *string
	longitude             *string
	refDate               *time.Time
}

/*
	type dataEntry struct {
		dataLabel string
		value     float32
	}

	type dayEntry struct {
		dayLabel    string
		dataEntries []dataEntry
	}

	type forecast struct {
		dayEntrys []dayEntry
	}
*/
func main() {
	forecastApi := new(forecastApi)
	emptyString := ""
	forecastApi.ForecastURL = &emptyString
	forecastApi.LastTenDaysURL = &emptyString
	forecastApi.HistoricalURLTemplate = &emptyString
	forecastApi.HistoricalURLTemplate = &emptyString
	lattitude := "47.497913"
	forecastApi.latitude = &lattitude
	longitude := "19.040236"
	forecastApi.longitude = &longitude
	now := time.Now()
	forecastApi.refDate = &now

	createForecastURL(*forecastApi)
	printResult(getCurrentForecast(*forecastApi))
	createLastTenDaysURL(*forecastApi)
	printResult(getLasttTenDaysForecast(*forecastApi))
	createHistoricalURLTemplate(*forecastApi)
	printResult(getHistoricalForecast(*forecastApi, "2023-01-10", "2023-01-12"))
}

func createForecastURL(forecastApi forecastApi) {
	forecastURL := "https://api.open-meteo.com/v1/forecast"
	forecastURL += "?" + "latitude=" + *forecastApi.latitude + "&longitude=" + *forecastApi.longitude
	forecastURL += "&current_weather=true"
	*forecastApi.ForecastURL = forecastURL
}

func createLastTenDaysURL(forecastApi forecastApi) {
	lastTenDaysURL := "https://api.open-meteo.com/v1/forecast"
	lastTenDaysURL += "?" + "latitude=" + *forecastApi.latitude + "&longitude=" + *forecastApi.longitude
	lastTenDaysURL += "&timezone=Europe/Budapest"
	lastTenDaysURL += "&past_days=10"
	lastTenDaysURL += "&daily=temperature_2m_min,temperature_2m_max"
	*forecastApi.LastTenDaysURL = lastTenDaysURL
}

func createHistoricalURLTemplate(forecastApi forecastApi) {
	historicalURLTemplate := "https://archive-api.open-meteo.com/v1/era5"
	historicalURLTemplate += "?" + "latitude=" + *forecastApi.latitude + "&longitude=" + *forecastApi.longitude
	historicalURLTemplate += "&timezone=Europe/Budapest"
	historicalURLTemplate += "&start_date=" + startDateMarker
	historicalURLTemplate += "&end_date=" + endDateMarker
	historicalURLTemplate += "&daily=temperature_2m_min,temperature_2m_max"
	*forecastApi.HistoricalURLTemplate = historicalURLTemplate
}

func getCurrentForecast(forecastApi forecastApi) (bool, string, []byte) {
	requestURL := *forecastApi.ForecastURL
	log.Println("Current forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)

	hasError := hasErrorInResponseBody(formattedBodyBytes)

	return hasError, "Current forecast", formattedBodyBytes
}

func getLasttTenDaysForecast(forecastApi forecastApi) (bool, string, []byte) {
	requestURL := *forecastApi.LastTenDaysURL
	log.Println("Last ten days forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)

	hasError := hasErrorInResponseBody(formattedBodyBytes)

	return hasError, "Last ten days forecast", formattedBodyBytes
}

func getHistoricalForecast(forecastApi forecastApi, startDate string, endDate string) (bool, string, []byte) {
	requestURL := *forecastApi.HistoricalURLTemplate
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

/*
func getMinusPlusOneDayForecast() {

}
*/

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
