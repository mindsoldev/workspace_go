package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const startDateMarker string = "<startDate>"
const endDateMarker string = "<endDate>"
const dailyresultMarker string = "<dailyresult>"

type forecastApi struct {
	ForecastURL                *string
	LastTenDaysURL             *string
	HistoricalURLTemplate      *string
	MinusPlusOneDayURLTemplate *string
	latitude                   *string
	longitude                  *string
	refDate                    *time.Time
	dailyResultParts           *[]string
	timeNames                  *[]string
	packedToArrayTitles        *[]string
}

type DataEntries map[string]float64
type DayEntries map[string]interface{}

func main() {
	forecastApi := new(forecastApi)
	emptyString := ""
	forecastApi.ForecastURL = &emptyString
	forecastApi.LastTenDaysURL = &emptyString
	forecastApi.HistoricalURLTemplate = &emptyString
	forecastApi.HistoricalURLTemplate = &emptyString
	forecastApi.MinusPlusOneDayURLTemplate = &emptyString
	lattitude := "47.497913"
	forecastApi.latitude = &lattitude
	longitude := "19.040236"
	forecastApi.longitude = &longitude
	now := time.Now()
	forecastApi.refDate = &now
	// dailyResultParts := []string{"temperature_2m_min", "temperature_2m_max"}
	dailyResultParts := []string{"temperature_2m_max"}
	forecastApi.dailyResultParts = &dailyResultParts
	forecastApi.timeNames = &[]string{"past", "now", "future"}
	forecastApi.packedToArrayTitles = &[]string{"past", "future"}

	createMinusPlusOneDayURLTemplate(*forecastApi)
	hasError := printResult(getMinusPlusOneDayForecast(*forecastApi))
	var exitCode int
	if hasError {
		exitCode = 1
	} else {
		exitCode = 0
	}
	os.Exit(exitCode)
}

func createMinusPlusOneDayURLTemplate(forecastApi forecastApi) {
	historicalURLTemplate := "https://api.open-meteo.com/v1/forecast"
	historicalURLTemplate += "?" + "latitude=" + *forecastApi.latitude + "&longitude=" + *forecastApi.longitude
	historicalURLTemplate += "&timezone=Europe/Budapest"
	historicalURLTemplate += "&start_date=" + startDateMarker
	historicalURLTemplate += "&end_date=" + endDateMarker
	historicalURLTemplate += "&daily=" + dailyresultMarker
	*forecastApi.HistoricalURLTemplate = historicalURLTemplate
}

func getMinusPlusOneDayForecast(forecastApi forecastApi) (bool, string, []byte) {
	requestURL := *forecastApi.MinusPlusOneDayURLTemplate
	startDate := forecastApi.refDate.AddDate(0, 0, -1)
	endDate := forecastApi.refDate.AddDate(0, 0, 1)
	requestURL = strings.Replace(requestURL, startDateMarker, startDate.Format("2006-01-02"), 1)
	requestURL = strings.Replace(requestURL, endDateMarker, endDate.Format("2006-01-02"), 1)
	dailyResultTemplate := createDailyResultTemplate(*forecastApi.dailyResultParts)
	requestURL = strings.Replace(requestURL, dailyresultMarker, dailyResultTemplate, 1)
	//log.Println("Minus-plus one day forecast URL: " + requestURL)

	response := executeQuery(requestURL)
	defer response.Body.Close()

	bodyBytes := readBody(response)

	formattedBodyBytes := formatJson(bodyBytes)
	hasError := hasErrorInResponseBody(formattedBodyBytes)

	var formattedResulBytearray []byte
	if !hasError {

		var result map[string]interface{}
		json.Unmarshal([]byte(formattedBodyBytes), &result)

		daily := result["daily"]
		days := make(DayEntries)

		for i := 0; i < 3; i++ {
			dataEntries := make(DataEntries)

			for _, partTitle := range *forecastApi.dailyResultParts {
				data := daily.(map[string]interface{})[partTitle].([]interface{})[i]
				dataEntries[partTitle] = data.(float64)
			}

			title := (*forecastApi.timeNames)[i]
			if Contains(*forecastApi.packedToArrayTitles, title) {
				days[title] = []interface{}{dataEntries}
			} else {
				days[title] = dataEntries
			}
		}

		resulBytearray := convertResultMapToJson(days)
		formattedResulBytearray = formatJson(resulBytearray)
	} else {
		formattedResulBytearray = formattedBodyBytes
	}

	return hasError, "Minus-plus one day forecast", formattedResulBytearray
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

func printResult(hasError bool, message string, result []byte) bool {
	errorMessage := ""
	if hasError {
		errorMessage = " error"
	}
	fmt.Printf("\n%s%s:\n%s\n\n", message, errorMessage, result)

	return hasError
}

func createDailyResultTemplate(dailyResultParts []string) string {
	var dailyResultTemplate string
	for _, template := range dailyResultParts {
		if len(dailyResultTemplate) > 0 {
			dailyResultTemplate += ","
		}
		dailyResultTemplate += template
	}
	return dailyResultTemplate
}

func convertResultMapToJson(structure DayEntries) []byte {
	resulBytearray, err := json.Marshal(structure)
	if err != nil {
		log.Fatal(err)
	}

	return resulBytearray
}

func Contains[T comparable](s []T, e T) bool {
	// source: https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
