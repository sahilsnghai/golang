package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/sahilsnghai/golang/Project11/internal/types"
	"github.com/valyala/fasthttp"
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

func WriteJson(ctx *fasthttp.RequestCtx, statusCode int, _error bool, v interface{}) error {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")

	jsonData, err := json.Marshal(map[string]interface{}{
		"version": map[string]string{"name": "askme-queue", "version": "v0.1"},
		"status":  map[string]string{"code": fmt.Sprint(statusCode), "value": fmt.Sprint(statusCode)},
		"data":    v,
		"error":   _error,
	})
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	_, err = ctx.Write(jsonData)
	return err
}

func MiddleWare(f types.APIFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		log.Println("Authorization checking")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			WriteJson(ctx, fasthttp.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": "Authorization header missing or Invalid Authorization header format",
			})
			return
		}

		var tokenString = authHeader[7:]
		parser := &jwt.Parser{}
		claims := jwt.MapClaims{}
		_, _, err := parser.ParseUnverified(tokenString, claims)
		if err != nil {
			WriteJson(ctx, fasthttp.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": fmt.Sprintf("Failed to decode token: %v", err),
			})
			return
		}

		log.Println("Token decoded")

		userVo, ok := claims["userVo"]
		if !ok {
			WriteJson(ctx, fasthttp.StatusUnauthorized, true, map[string]interface{}{
				"error":   true,
				"message": "userVo not found in token",
			})
			return
		}
		log.Printf("userVo: %+v", userVo)

		ctx.SetUserValue("userVo", userVo)

		f(ctx)
	}
}
