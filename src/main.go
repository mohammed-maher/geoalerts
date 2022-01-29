package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type Location struct {
	Name      string
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

var availableLocations = []Location{
	{
		Name:      "Baxtiyari",
		Longitude: 36.21180032040089,
		Latitude:  43.99020072567531,
	},
	{
		Name:      "ANKAWA",
		Longitude: 36.218681,
		Latitude:  43.993753,
	},
	{
		Name:      "Tayrawa",
		Longitude: 36.19827565893692,
		Latitude:  44.011546155683796,
	},
}

const MINIMUM_DISTANCE = 0.5 //km

func main() {
	http.HandleFunc("/api", EligibleLocations)
	http.ListenAndServe(":8080", nil)
}

func EligibleLocations(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	ba, err := io.ReadAll(req.Body)
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	var loc Location
	err = json.Unmarshal(ba, &loc)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	nearbyLocations := []Location{}
	for _, av := range availableLocations {
		dst := av.Distance(loc)
		if dst < MINIMUM_DISTANCE {
			nearbyLocations = append(nearbyLocations, av)
		}
	}
	rsp, err := json.Marshal(nearbyLocations)
	if err != nil {
		fmt.Fprint(res, err.Error())
		return
	}
	fmt.Fprint(res, string(rsp))
}

func (origin Location) Distance(destination Location) float64 {
	degreesLat := degrees2radians(destination.Latitude - origin.Latitude)
	degreesLong := degrees2radians(destination.Longitude - origin.Longitude)
	a := (math.Sin(degreesLat/2)*math.Sin(degreesLat/2) +
		math.Cos(degrees2radians(origin.Latitude))*
			math.Cos(degrees2radians(destination.Latitude))*math.Sin(degreesLong/2)*
			math.Sin(degreesLong/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := radius * c
	return d
}

const radius = 6371 // Earth's mean radius in kilometers

func degrees2radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
