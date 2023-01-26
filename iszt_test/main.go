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
}

type dataEntry struct {
	DataLabel string
	Value     float64
}

type dayEntry struct {
	DayLabel    string
	DataEntries []dataEntry
}

type forecastEntry struct {
	DayEntries []dayEntry
}

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
	//dailyResultParts := []string{"temperature_2m_min", "temperature_2m_max"}
	dailyResultParts := []string{"temperature_2m_max"}
	forecastApi.dailyResultParts = &dailyResultParts
	forecastApi.timeNames = &[]string{"past", "now", "future"}

	createMinusPlusOneDayURLTemplate(*forecastApi)
	printResult(getMinusPlusOneDayForecast(*forecastApi))
	os.Exit(1)
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

	var result map[string]interface{}
	json.Unmarshal([]byte(formattedBodyBytes), &result)

	daily := result["daily"]

	var dataItems []dataEntry
	for _, partTitle := range *forecastApi.dailyResultParts {
		dataItems = append(dataItems, createDataItem(daily, partTitle, forecastApi))
	}

	var dayItems []dayEntry
	for i := 0; i < 3; i++ {
		dayItem := dayEntry{
			DayLabel:    (*forecastApi.timeNames)[i],
			DataEntries: dataItems,
		}
		dayItems = append(dayItems, dayItem)
	}

	forecast := forecastEntry{
		DayEntries: dayItems,
	}

	resulBytearray := convertStructToJson(forecast)
	formattedResulBytearray := formatJson(resulBytearray)

	return hasError, "Minus-plus one day forecast:", formattedResulBytearray
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

func createDataItem(daily interface{}, partName string, forecastApi forecastApi) dataEntry {

	temperatures := daily.(map[string]interface{})[partName]
	temperature0 := temperatures.([]interface{})[0]

	dataItem := dataEntry{
		DataLabel: (*forecastApi.dailyResultParts)[0],
		Value:     temperature0.(float64),
	}

	return dataItem
}

func convertStructToJson(structure forecastEntry) []byte {
	resulBytearray, err := json.Marshal(structure)
	if err != nil {
		log.Fatal(err)
	}

	return resulBytearray
}
