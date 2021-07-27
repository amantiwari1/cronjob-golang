package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Type

type Weatherapi struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Current struct {
	LastUpdatedEpoch int64              `json:"last_updated_epoch"`
	LastUpdated      string             `json:"last_updated"`
	TempC            float64            `json:"temp_c"`
	TempF            float64            `json:"temp_f"`
	IsDay            int64              `json:"is_day"`
	Condition        Condition          `json:"condition"`
	WindMph          float64            `json:"wind_mph"`
	WindKph          float64            `json:"wind_kph"`
	WindDegree       int64              `json:"wind_degree"`
	WindDir          string             `json:"wind_dir"`
	PressureMB       float64            `json:"pressure_mb"`
	PressureIn       float64            `json:"pressure_in"`
	PrecipMm         float64            `json:"precip_mm"`
	PrecipIn         float64            `json:"precip_in"`
	Humidity         int64              `json:"humidity"`
	Cloud            int64              `json:"cloud"`
	FeelslikeC       float64            `json:"feelslike_c"`
	FeelslikeF       float64            `json:"feelslike_f"`
	VisKM            float64            `json:"vis_km"`
	VisMiles         float64            `json:"vis_miles"`
	Uv               float64            `json:"uv"`
	GustMph          float64            `json:"gust_mph"`
	GustKph          float64            `json:"gust_kph"`
	AirQuality       map[string]float64 `json:"air_quality"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int64  `json:"code"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Weather struct {
	gorm.Model
	Name        string  `json:"name"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	Temp_c      float64 `json:"temp_c"`
	Humidity    int64   `json:"humidity"`
	Cloud       int64   `json:"cloud"`
	Feelslike_c float64 `json:"feelslike_c"`
}

var (
	db *gorm.DB
)

// Connect to databasee
func Connect() {
	d, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = d
}

// get the database
func GetDB() *gorm.DB {
	return db
}

// init for connect, get db to auto migrate
func init() {
	Connect()
	db = GetDB()
	db.AutoMigrate(&Weather{})
}

// parse json data from get request data
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// fetch and save database function
func job() {
	weatherapi := &Weatherapi{}

	getJson("https://api.weatherapi.com/v1/current.json?key="+(os.Getenv("WEATHER_API_KEY"))+"&q=India&aqi=yes", weatherapi)

	weather := Weather{Name: weatherapi.Location.Name, Region: weatherapi.Location.Region, Country: weatherapi.Location.Country, Temp_c: weatherapi.Current.TempC, Humidity: weatherapi.Current.Humidity, Cloud: weatherapi.Current.Cloud, Feelslike_c: weatherapi.Current.FeelslikeC}

	db.Create(&weather)

	log.Println("Fetched and Saved Database Successfull")

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Running Cronjob")
	gocron.Every(1).Day().At("01:02").Do(job)
	gocron.Every(1).Day().At("01:03").Do(job)
	gocron.Every(1).Day().At("01:04").Do(job)

	<-gocron.Start()

}
