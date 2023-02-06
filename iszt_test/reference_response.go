package main

import (
	"io"
	"os"
	"time"
)

func saveApprovedResponse() {
	forecastApi := initForecastApi()
	refDate := time.Date(2023, 01, 28, 12, 00, 00, 000, time.Local)
	forecastApi.refDate = &refDate
	createMinusPlusOneDayURLTemplate(*forecastApi)
	requestUrl := replaceParametersInUrlTemplate(*forecastApi)
	response := executeQuery(requestUrl)
	defer response.Body.Close()
	out, err := os.Create("approved_response.txt")
	if err != nil {
		panic("File creation error")
	}
	defer out.Close()
	io.Copy(out, response.Body)
}

func saveApprovedResult() {
	forecastApi := initForecastApi()
	refDate := time.Date(2023, 01, 28, 12, 00, 00, 000, time.Local)
	forecastApi.refDate = &refDate
	createMinusPlusOneDayURLTemplate(*forecastApi)
	_, _, formattedResulBytearray := getMinusPlusOneDayForecast(*forecastApi)
	err := os.WriteFile("approved_result.txt", formattedResulBytearray, 0644)
	if err != nil {
		panic("File creation error")
	}
}
