package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WeatherAPI struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	var args = os.Args
	var loc = args[1]

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=4e221ed9011447169fa162606242908&q=" + loc + "&days=7&aqi=yes&alerts=yes")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("WeatherAPI not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// print responsebody
	// fmt.Println(string(body))

	var weather WeatherAPI
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	txt := getSymbol(current.Condition.Text)

	fmt.Print("\n---------- Current Weather -----------\n\n")
	fmt.Printf("%s , %s : %.0f°C %s \n",
		location.Name,
		location.Country,
		current.TempC,
		txt,
	)

	fmt.Print("\n---------- Forecast for today -----------\n\n")
	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before((time.Now())) {
			continue
		}
		txt := getSymbol(hour.Condition.Text)

		fmt.Printf("%s - %.0f°C, %.0f%% %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			txt)
	}

}

func getSymbol(s string) string {

	switch {

	case strings.Contains(s, "rain"):
		s = "⛅"
	case strings.Contains(s, "sun"):
		s = "☀️"

	}
	return s

}
