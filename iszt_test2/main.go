package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type forecast struct {
	Latitude              string
	Longitude             string
	Generationtime_ms     float64
	Utc_offset_seconds    int
	Timezone              string
	Timezone_abbreviation string
	Elevation             int
	HorlyUnits            horlyUnits
	Times                 []timeItem
	WindSpeeds            []windspeed
	WindDirections        []winddirection
}

type horlyUnits struct {
	Time               string
	Windspeed_10m      string
	Windspeed_80m      string
	Windspeed_120m     string
	Windspeed_180m     string
	Winddirection_10m  string
	Winddirection_80m  string
	Winddirection_120m string
	Winddirection_180m string
}

type timeItem struct {
	timeValue string
}

type windspeed struct {
	speed float64
}

type winddirection struct {
	speed float64
}

func main() {

	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&past_days=10&hourly=windspeed_10m,windspeed_80m,windspeed_120m,windspeed_180m,winddirection_10m,winddirection_80m,winddirection_120m,winddirection_180m")
	if err != nil {
		log.Fatal(err)
	}

	bodyBytesArray, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var buf interface{}
	if err := json.Unmarshal(bodyBytesArray, &buf); err != nil {
		log.Fatal(err)
	}
	data := buf.(map[string]interface{})["hourly"]
	times := data.(map[string]interface{})["time"]
	windSpeed10 := data.(map[string]interface{})["windspeed_10m"]
	windSpeed80 := data.(map[string]interface{})["windspeed_80m"]
	windSpeed120 := data.(map[string]interface{})["windspeed_120m"]
	windSpeed180 := data.(map[string]interface{})["windspeed_180m"]
	windDirection10 := data.(map[string]interface{})["winddirection_10m"]
	windDirection80 := data.(map[string]interface{})["winddirection_80m"]
	windDirection120 := data.(map[string]interface{})["winddirection_120m"]
	windDirection180 := data.(map[string]interface{})["winddirection_180m"]

	fmt.Println(times)
	fmt.Println(windSpeed10)
	fmt.Println(windSpeed80)
	fmt.Println(windSpeed120)
	fmt.Println(windSpeed180)
	fmt.Println(windDirection10)
	fmt.Println(windDirection80)
	fmt.Println(windDirection120)
	fmt.Println(windDirection180)

	// var result bytes.Buffer
	// err = json.Indent(&result, bodyBytesArray, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)
}
