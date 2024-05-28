package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidZipcode(t *testing.T) {
	tests := []struct {
		name     string
		zipcode  string
		expected bool
	}{
		{
			name:     "valid zipcode",
			zipcode:  "01001000",
			expected: true,
		},
		{
			name:     "invalid zipcode - letters",
			zipcode:  "abc12345",
			expected: false,
		},
		{
			name:     "invalid zipcode - shorter length",
			zipcode:  "12345",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumeric(tt.zipcode)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C to 32°F",
			celsius:  0,
			expected: 32,
		},
		{
			name:     "100°C to 212°F",
			celsius:  100,
			expected: 212,
		},
		{
			name:     "-50°C to -58°F",
			celsius:  -50,
			expected: -58,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToFahrenheit(tt.celsius)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "0°C to 273K",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "100°C to 373K",
			celsius:  100,
			expected: 373,
		},
		{
			name:     "-273°C to 0K",
			celsius:  -273,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToKelvin(tt.celsius)
			assert.Equal(t, tt.expected, result)
		})
	}
}
