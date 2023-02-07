package main

import (
	"io"
	"os"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Forecast", func() {

	Describe("Run query", func() {
		Context("When the query done", func() {
			It("Should be equals with the approwed response", func() {
				approvedContent, err := os.ReadFile("approved_response.txt")
				Expect(err).ShouldNot(HaveOccurred())

				cleanedApprovedContent := removeGeneretionTime(string(approvedContent))

				forecastApi := initForecastApi()
				refDate := time.Date(2023, 01, 28, 12, 00, 00, 000, time.Local)
				forecastApi.refDate = &refDate
				createMinusPlusOneDayURLTemplate(*forecastApi)
				requestUrl := replaceParametersInUrlTemplate(*forecastApi)
				response := executeQuery(requestUrl)
				body, err := io.ReadAll(response.Body)
				defer response.Body.Close()

				cleanedBody := removeGeneretionTime(string(body))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(cleanedBody).To(Equal(cleanedApprovedContent))
			})
		})
	})

	Describe("Processing data", func() {
		Context("When the query processed", func() {
			It("Should be equals with the approwed result", func() {
				approvedContent, err := os.ReadFile("approved_result.txt")
				Expect(err).ShouldNot(HaveOccurred())

				forecastApi := initForecastApi()
				refDate := time.Date(2023, 01, 28, 12, 00, 00, 000, time.Local)
				forecastApi.refDate = &refDate
				createMinusPlusOneDayURLTemplate(*forecastApi)
				_, _, formattedResulBytearray := getMinusPlusOneDayForecast(*forecastApi)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(formattedResulBytearray)).To(Equal(string(approvedContent)))
			})
		})
	})
})

func removeGeneretionTime(text string) string {
	s := regexp.MustCompile(`"generationtime_ms\": ?[\d]{1}[.]{1}[\d]*,`).Split(string(text), 2)
	return s[0] + s[1]
}
