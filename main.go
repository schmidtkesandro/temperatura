package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"unicode"

	"github.com/go-chi/chi/v5"
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
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}
type Response struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	r := chi.NewRouter()
	fmt.Println("Starting")
	r.Get("/cep/{cep}", GetTemperatureByCep)

	http.ListenAndServe(":8080", r)
}

func GetTemperatureByCep(w http.ResponseWriter, r *http.Request) {
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
	city := address.Localidade
	// Fetch weather information
	temperature, err := getTemperature(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to fetch temperature "+err.Error())
		return
	}

	response := Response{
		TempC: temperature.Current.TempC,
		TempF: celsiusToFahrenheit(temperature.Current.TempC),
		TempK: celsiusToKelvin(temperature.Current.TempC),
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAddress(cep string) (*Address, error) {
	var address Address

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return &address, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&address)

	if err != nil || address.Localidade == "" {
		return &address, fmt.Errorf("CEP not found")
	}

	return &address, nil
}

func getTemperature(city string) (*WeatherResponse, error) {
	apiKey := "aa5a63b17c16446cbdc24451240205"

	// Escape a cidade para lidar com caracteres especiais
	escapedCity := url.QueryEscape(city)
	link := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, escapedCity)

	resp, err := http.Get(link)

	if resp.StatusCode != 200 {
		link = fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, removeAcentos(city))
		resp, err = http.Get(link)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

type Config struct {
	WeatherAPIKey string `json:"weather_api_key"`
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
func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
