package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct{}

func NewWeather() *Weather {
	return &Weather{}
}

type Class struct {
	Name     string   `json:"name"`
	EnName   string   `json:"enName"`
	Parent   string   `json:"parent"`
	Children []string `json:"children"`
}

type Area struct {
	Centers map[string]struct {
		Name       string   `json:"name"`
		EnName     string   `json:"enName"`
		OfficeName string   `json:"officeName"`
		Children   []string `json:"children"`
	} `json:"centers"`
	Offices map[string]struct {
		Name       string   `json:"name"`
		EnName     string   `json:"enName"`
		OfficeName string   `json:"officeName"`
		Parent     string   `json:"parent"`
		Children   []string `json:"children"`
	} `json:"offices"`
	Class10s map[string]Class `json:"class10s"`
	Class15s map[string]Class `json:"class15s"`
	Class20s map[string]Class `json:"class20s"`
}

type AmedasArea struct {
	Type   string    `json:"type"`
	Elems  string    `json:"elems"`
	Lat    []float32 `json:"lat"`
	Lon    []float32 `json:"lon"`
	Alt    []int     `json:"alt"`
	KjName string    `json:"kjName"`
	KnName string    `json:"knName"`
	EnName string    `json:"enName"`
}

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type Amedas struct {
	PrefNumber        int       `json:"prefNumber"`
	ObservationNumber int       `json:"observationNumber"`
	Temp              []float32 `json:"temp"`
	Sum10m            []float32 `json:"sum10m"`
	Sum1h             []float32 `json:"sum1h"`
	Precipitation10m  []float32 `json:"precipitation10m"`
	Precipitation1h   []float32 `json:"precipitation1h"`
	Precipitation3h   []float32 `json:"precipitation3h"`
	Precipitation24h  []float32 `json:"precipitation24h"`
	WindDirection     []float32 `json:"windDirection"`
	Wind              []float32 `json:"wind"`
	MaxTempTime       Time      `json:"maxTempTime"`
	MaxTemp           []float32 `json:"maxTemp"`
	MinTempTime       Time      `json:"minTempTime"`
	MinTemp           []float32 `json:"minTemp"`
	GustTime          Time      `json:"gustTime"`
	GustDirection     []float32 `json:"gustDirection"`
	Gust              []float32 `json:"gust"`
}

func (Weather) GetWeather(location string) (Amedas, error) {
	res, _ := http.Get("https://www.jma.go.jp/bosai/amedas/const/amedastable.json")
	body, _ := io.ReadAll(res.Body)
	var amedasArea map[string]AmedasArea
	_ = json.Unmarshal(body, &amedasArea)
	for id, data := range amedasArea {
		if data.KjName == location {
			now := time.Now()
			res, _ = http.Get(fmt.Sprintf("https://www.jma.go.jp/bosai/amedas/data/point/%s/%s.json", id, func() string {
				if temp := now.Hour() % 3; temp != 0 {
					return now.Add(-time.Duration(temp) * time.Hour).Format("20060102_15")
				} else {
					return now.Format("20060102_15")
				}
			}()))
			body, _ = io.ReadAll(res.Body)
			var amedas map[string]Amedas
			_ = json.Unmarshal(body, &amedas)

			i := 1
			for _, data := range amedas {
				if i == len(amedas) {
					return data, nil
				} else {
					i++
				}
			}
		}
	}
	return Amedas{}, fmt.Errorf("Unknown error")
}
