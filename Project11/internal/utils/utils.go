package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/sahilsnghai/golang/Project11/internal/types"
)

func GetConfig() (*types.Config, error) {
	// basePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	fmt.Println("Error getting base path:", err)
	// 	return nil, err
	// }

	// filename := filepath.Join(basePath, "configuration/configuration.json")
	log.Println("reading configuration folder")
	filename, _ := filepath.Abs("configuration/configuration.json")
	//  := "configuration/configuration.json"

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config types.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &config, nil
}

func WriteJson(w http.ResponseWriter, status_code int, _error bool, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status_code)

	// resultsJSON, err := json.Marshal(v)
	// if err != nil {
	//     return nil, fmt.Errorf("error marshalling results to JSON: %v", err)
	// }

	// if jsonData, ok := v.([]byte); ok {
	// 	_, err := w.Write(jsonData)
	// 	return err
	// }

	// // If not a byte slice, encode it as JSON
	// return json.NewEncoder(w).Encode(map[string]any{"data": v})
	jsonData, err := json.Marshal(map[string]any{
		"version": map[string]string{"name": "askme-queue", "versionCode": "0.1"},
		"status":  map[string]string{"code": fmt.Sprint(status_code), "value": fmt.Sprint(status_code)},
		"data":    v,
		"error":   _error,
	})
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	// Write the JSON to the response
	_, err = w.Write(jsonData)
	return err
}

func MiddleWare(f types.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		log.Println("Authorization checking")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			WriteJson(w, http.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": fmt.Sprint("Authorization header missing or Invalid Authorization header format"),
			})
			return
		}

		var tokenString = authHeader[7:]

		parser := &jwt.Parser{}
		claims := jwt.MapClaims{}
		_, _, err := parser.ParseUnverified(tokenString, claims)
		if err != nil {
			WriteJson(w, http.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": fmt.Sprintf("Failed to decode token: %v", err),
			})
			return
		}

		log.Println("token decode")

		userVo, ok := claims["userVo"]
		if !ok {
			http.Error(w, "userVo not found in token", http.StatusUnauthorized)
			WriteJson(w, http.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": fmt.Sprint("userVo not found in token"),
			})
			return
		}
		log.Printf("userVo: %+v", userVo)

		ctx := context.WithValue(r.Context(), "userVo", userVo)
		if err := f(w, r.WithContext(ctx)); err != nil {
			WriteJson(w, http.StatusBadRequest, true, types.APIError{Error: err.Error()})
		}
	}

}
