package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"unicode"

	"github.com/go-chi/chi"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type AddressGetter interface {
	GetAddress(cep string) (Address, error)
}

type WeatherGetter interface {
	GetTemperature(address Address) (float64, error)
}

func main() {
	r := chi.NewRouter()
	fmt.Println("Starting")
	r.Get("/cep/{cep}", GetTemperature)

	http.ListenAndServe(":8080", r)
}

func GetTemperature(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	// Check if the cep has 8 digits
	if len(cep) != 8 || !isNumeric(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "invalid zipcode")
		return
	}

	// Fetch address information
	address, err := getAddress(cep)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "can not find zipcode")
		return
	}

	// Fetch weather information
	temperature, err := getTemperature(address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to fetch temperature")
		return
	}

	response := WeatherResponse{
		TempC: temperature,
		TempF: temperature*1.8 + 32,
		TempK: temperature + 273.15,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAddress(cep string) (Address, error) {
	var address Address

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return address, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&address)

	if err != nil || address.Localidade == "" {
		return address, fmt.Errorf("CEP not found")
	}

	return address, nil
}

func getTemperature(address Address) (float64, error) {
	// You should replace 'YOUR_API_KEY' with your WeatherAPI key

	config, err := loadConfig("config.json")
	if err != nil {
		return 0, err
	}
	apiKey := config.WeatherAPIKey

	// Escape a cidade para lidar com caracteres especiais
	//escapedCity := url.QueryEscape(address.Localidade)
	escapedCity := "Brasília"
	link := "http://api.weatherapi.com/v1/current.json?q=" + escapedCity + "&key=" + apiKey

	resp, _ := http.Get(link)
	//fmt.Println(err)
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		link = "http://api.weatherapi.com/v1/current.json?q=" + removeAcentos(address.Localidade) + "&key=" + apiKey
	}
	fmt.Println(link)
	resp, err = http.Get(link)
	//fmt.Println(resp)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var weatherData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	//fmt.Println("err 1", err)
	//fmt.Println(resp.Body)

	if err != nil {
		return 0, err
	}
	currentData, ok := weatherData["current"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("current data not found in weather response")
	}
	fmt.Println(weatherData)
	tempC, ok := currentData["temp_c"].(float64)
	if !ok {
		return 0, err
	}

	return tempC, nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

type Config struct {
	WeatherAPIKey string `json:"weather_api_key"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
func removeAcentos(input string) string {
	// Cria um transformador para normalizar a string
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		// Retorna true para remover todos os caracteres que não são ASCII
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)

	// Aplica a transformação na string de entrada
	result, _, _ := transform.String(t, input)

	return result
}
