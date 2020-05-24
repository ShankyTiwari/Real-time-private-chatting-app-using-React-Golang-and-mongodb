package utils

import (
	"github.com/rs/cors"
)

// GetCorsConfig will return the CORS values
func GetCorsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                       // All origins
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}, // Allowing only get, just an example
	})
}
