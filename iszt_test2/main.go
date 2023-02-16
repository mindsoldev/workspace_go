package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/mohae/struct2csv"
	"github.com/tidwall/gjson"
)

type hourlyItem struct {
	Time             time.Time
	WindDirection10  int
	WindDirection80  int
	WindDirection120 int
	WindDirection180 int
	WindSpeed10      float64
	WindSpeed80      float64
	WindSpeed120     float64
	WindSpeed180     float64
}

type dailyItem struct {
	Day         string
	HourlyItems []hourlyItem
}

type dailyItemAveraged struct {
	Day        string
	HourlyItem hourlyItem
}

const resultTimePattern string = "2006-01-02T15:04"
const datePattern string = "2006-01-02"

type parameters struct {
	url         string
	csvFileName string
}

var params parameters

func main() {

	params = parameters{
		url:         "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&past_days=10&hourly=windspeed_10m,windspeed_80m,windspeed_120m,windspeed_180m,winddirection_10m,winddirection_80m,winddirection_120m,winddirection_180m",
		csvFileName: "winds.csv",
	}

	resp, err := http.Get(params.url)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytesArray, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	days := collectDatas(bodyBytesArray)

	daysAveraged := dataAggregation(days)

	printTable(daysAveraged)

	buff := createCsv(daysAveraged)

	fmt.Println(buff.String())
	saveCsvToFile(buff)
}

func collectDatas(bodyBytesArray []byte) []dailyItem {
	hourly := gjson.Get(string(bodyBytesArray), "hourly")
	times := hourly.Get("time").Array()
	windSpeed10 := hourly.Get("windspeed_10m").Array()
	windSpeed80 := hourly.Get("windspeed_80m").Array()
	windSpeed120 := hourly.Get("windspeed_120m").Array()
	windSpeed180 := hourly.Get("windspeed_180m").Array()
	windDirection10 := hourly.Get("winddirection_10m").Array()
	windDirection80 := hourly.Get("winddirection_80m").Array()
	windDirection120 := hourly.Get("winddirection_120m").Array()
	windDirection180 := hourly.Get("winddirection_180m").Array()

	days := []dailyItem{}
	hourlyItems := []hourlyItem{}
	beforeTime := parseTimeFromJsonResult(times[0])
	for i, timeResult := range times {
		time := parseTimeFromJsonResult(timeResult)
		if !isSomeDay(beforeTime, time) {
			daylyI := dailyItem{
				Day:         beforeTime.Format(datePattern),
				HourlyItems: hourlyItems,
			}
			days = append(days, daylyI)
			hourlyItems = []hourlyItem{}
		}

		hourlyI := hourlyItem{
			Time:             time,
			WindDirection10:  int(windDirection10[i].Int()),
			WindDirection80:  int(windDirection80[i].Int()),
			WindDirection120: int(windDirection120[i].Int()),
			WindDirection180: int(windDirection180[i].Int()),
			WindSpeed10:      windSpeed10[i].Float(),
			WindSpeed80:      windSpeed80[i].Float(),
			WindSpeed120:     windSpeed120[i].Float(),
			WindSpeed180:     windSpeed180[i].Float(),
		}

		hourlyItems = append(hourlyItems, hourlyI)
		beforeTime = time
	}

	return days
}

func parseTimeFromJsonResult(timeResult gjson.Result) time.Time {
	time, err := time.Parse(resultTimePattern, timeResult.Str)
	if err != nil {
		log.Fatal(err)
	}

	return time
}

func isSomeDay(t1, t2 time.Time) bool {
	var result bool = false
	year1, month1, day1 := t1.Date()
	year2, month2, day2 := t2.Date()
	if year1 == year2 && month1 == month2 && day1 == day2 {
		result = true
	}

	return result
}

func dataAggregation(days []dailyItem) []dailyItemAveraged {
	daysAveraged := []dailyItemAveraged{}
	for _, day := range days {
		hourlyItems := day.HourlyItems
		var windSpeed10Sum float64 = 0.0
		var windSpeed80Sum float64 = 0.0
		var windSpeed120Sum float64 = 0.0
		var windSpeed180Sum float64 = 0.0

		windDirection10Sum := 0
		windDirection80Sum := 0
		windDirection120Sum := 0
		windDirection180Sum := 0
		for _, hItem := range hourlyItems {
			windSpeed10Sum += hItem.WindSpeed10
			windSpeed80Sum += hItem.WindSpeed80
			windSpeed120Sum += hItem.WindSpeed120
			windSpeed180Sum += hItem.WindSpeed180
			windDirection10Sum += hItem.WindDirection10
			windDirection80Sum += hItem.WindDirection80
			windDirection120Sum += hItem.WindDirection120
			windDirection180Sum += hItem.WindDirection180
		}

		lenght := len(hourlyItems)
		hourlyItemAveraged := hourlyItem{
			Time:             hourlyItems[0].Time,
			WindDirection10:  windDirection10Sum / lenght,
			WindDirection80:  windDirection80Sum / lenght,
			WindDirection120: windDirection120Sum / lenght,
			WindDirection180: windDirection180Sum / lenght,
			WindSpeed10:      roundTo2decimal(windSpeed10Sum / float64(lenght)),
			WindSpeed80:      roundTo2decimal(windSpeed80Sum / float64(lenght)),
			WindSpeed120:     roundTo2decimal(windSpeed120Sum / float64(lenght)),
			WindSpeed180:     roundTo2decimal(windSpeed180Sum / float64(lenght)),
		}

		daylyItemAveraged := dailyItemAveraged{
			Day:        hourlyItemAveraged.Time.Format(datePattern),
			HourlyItem: hourlyItemAveraged,
		}

		daysAveraged = append(daysAveraged, daylyItemAveraged)
	}

	return daysAveraged
}

func roundTo2decimal(value float64) float64 {
	return math.Round(value*100) / 100
}

func printTable(daysAveraged []dailyItemAveraged) {
	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("Day              Direction        |         Speed             |\n")
	fmt.Printf("          |  10m|  80m| 120m| 180m|   10m|   80m|  120m|  180m|\n")

	for _, dayAveraged := range daysAveraged {
		fmt.Printf("%s|  %d|  %d|  %d|  %d| %5.2f| %5.2f| %5.2f| %5.2f|\n",
			dayAveraged.Day,
			dayAveraged.HourlyItem.WindDirection10,
			dayAveraged.HourlyItem.WindDirection80,
			dayAveraged.HourlyItem.WindDirection120,
			dayAveraged.HourlyItem.WindDirection180,
			dayAveraged.HourlyItem.WindSpeed10,
			dayAveraged.HourlyItem.WindSpeed80,
			dayAveraged.HourlyItem.WindSpeed120,
			dayAveraged.HourlyItem.WindSpeed180)
	}

	fmt.Println("---------------------------------------------------------------")
}

func createCsv(daysAveraged []dailyItemAveraged) *bytes.Buffer {
	buff := &bytes.Buffer{}

	w := struct2csv.NewWriter(buff)

	err := w.Write([]string{"Day", "Dir_10m", "Dir_80m", "Dir_120m", "Dir_180m", "Speed_10m", "Speed_80m", "Speed_120m", "Speed_180m"})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range daysAveraged {
		err = w.WriteStruct(v)
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()

	return buff
}

func saveCsvToFile(buff *bytes.Buffer) {
	file, err := os.Create(params.csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(buff.Bytes())
}
