package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks para AddressGetter e WeatherGetter
type MockAddressGetter struct {
	mock.Mock
}

func (m *MockAddressGetter) GetAddress(cep string) (Address, error) {
	args := m.Called(cep)
	return args.Get(0).(Address), args.Error(1)
}

type MockWeatherGetter struct {
	mock.Mock
}

func (m *MockWeatherGetter) GetTemperature(address Address) (float64, error) {
	args := m.Called(address)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetTemperature(t *testing.T) {
	router := chi.NewRouter()

	// Substituir com mocks
	mockAddressGetter := new(MockAddressGetter)
	mockWeatherGetter := new(MockWeatherGetter)

	router.Get("/cep/{cep}", GetTemperature)

	// Teste com CEP válido
	mockAddressGetter.On("GetAddress", "12345678").Return(Address{Localidade: "São Paulo", Uf: "SP"}, nil)
	mockWeatherGetter.On("GetTemperature", Address{Localidade: "São Paulo", Uf: "SP"}).Return(22.0, nil)

	request, _ := http.NewRequest("GET", "/cep/12345678", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	var weatherResp WeatherResponse
	json.NewDecoder(response.Body).Decode(&weatherResp)
	assert.Equal(t, 22.0, weatherResp.TempC)

	// Teste com CEP inválido
	request, _ = http.NewRequest("GET", "/cep/1234", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

	// Teste com erro na busca do endereço
	err := ""
	mockAddressGetter.On("GetAddress", "87654321").Return(Address{}, err)
	request, _ = http.NewRequest("GET", "/cep/87654321", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)

	// Teste com erro na busca da temperatura
	mockAddressGetter.On("GetAddress", "12345678").Return(Address{Localidade: "São Paulo", Uf: "SP"}, nil)
	mockWeatherGetter.On("GetTemperature", Address{Localidade: "São Paulo", Uf: "SP"}).Return(0.0, fmt.Errorf("API error"))
	request, _ = http.NewRequest("GET", "/cep/12345678", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
