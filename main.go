package main

import (
	"net/http"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		panic("OPENAI_API_KEY not set in .env file")
	}

	router := gin.Default()

	// 信頼するプロキシを設定
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/api/openai", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+openaiApiKey).
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			Post("https://api.openai.com/v1/chat/completions")

		if err != nil {
			log.Println("Request error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Request body:", requestBody)
		log.Println("Response status:", resp.Status())
		log.Println("Response body:", resp.String())

		if resp.StatusCode() == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Data(resp.StatusCode(), "application/json", resp.Body())
	})

	router.Run(":3000")
}
