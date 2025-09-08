package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatGPTRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type ChatGPTResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    openaiApiKey := os.Getenv("OPENAI_API_KEY")
    if openaiApiKey == "" {
        log.Fatal("OPENAI_API_KEY not set in .env file")
    }

    router := gin.Default()

    // CORSミドルウェアを追加
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    router.Use(cors.New(config))

    router.POST("/api/openai", func(c *gin.Context) {
        var requestBody ChatGPTRequest
        if err := c.ShouldBindJSON(&requestBody); err != nil {
            log.Println("Failed to bind JSON:", err)
            log.Println("Request body:", c.Request.Body)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        log.Printf("Received request: %+v\n", requestBody)

        if len(requestBody.Messages) == 0 {
            log.Println("No messages in request")
            c.JSON(http.StatusBadRequest, gin.H{"error": "No messages in request"})
            return
        }

        // モデルが設定されていない場合のデフォルト値
        if requestBody.Model == "" {
            requestBody.Model = "gpt-3.5-turbo"
        }

        // メッセージの内容を修正
        newMessage := Message{
            Role: "user",
            Content: requestBody.Messages[0].Content + "\n\nネットワークに関する、〇か×で答えられる二者択一形式の問題を1つだけ生成してください。複数の問題は絶対に生成しないでください。以下の形式で厳密に出力してください：\n\n問題: [ここに1つの問題文を入れてください]\n正解: [〇または×]\n解説: [ここに解説を入れてください]",
        }

        // 新しいリクエストボディを作成
        newRequestBody := ChatGPTRequest{
            Model:    requestBody.Model,
            Messages: []Message{newMessage},
        }

        log.Printf("Sending request to OpenAI: %+v\n", newRequestBody)

        client := resty.New()
        resp, err := client.R().
            SetHeader("Authorization", "Bearer "+openaiApiKey).
            SetHeader("Content-Type", "application/json").
            SetBody(newRequestBody).
            Post("https://api.openai.com/v1/chat/completions")

        if err != nil {
            log.Println("Request error:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        log.Println("Response status:", resp.Status())
        log.Println("Response body:", resp.String())

        if resp.StatusCode() != http.StatusOK {
            log.Printf("Unexpected status code: %d", resp.StatusCode())
            c.JSON(resp.StatusCode(), gin.H{"error": fmt.Sprintf("Unexpected status code: %d", resp.StatusCode())})
            return
        }

        var chatGPTResponse ChatGPTResponse
        err = json.Unmarshal(resp.Body(), &chatGPTResponse)
        if err != nil {
            log.Println("Failed to parse ChatGPT response:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ChatGPT response"})
            return
        }

        if len(chatGPTResponse.Choices) == 0 {
            log.Println("No choices in ChatGPT response")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "No choices in ChatGPT response"})
            return
        }

        content := chatGPTResponse.Choices[0].Message.Content
        log.Println("Original content:", content)

        questions := strings.Split(content, "\n\n")
        if len(questions) > 0 {
            content = questions[0]
        }

        log.Println("Processed content:", content)

        re := regexp.MustCompile(`問題:\s*(.+)\s*\n正解:\s*(.+)(?:\s*\n解説:\s*(.+))?`)
        matches := re.FindStringSubmatch(content)

        if len(matches) >= 3 {
            formattedQuestion := "問題: " + matches[1] + "\n正解: " + matches[2]
            if len(matches) > 3 && matches[3] != "" {
                formattedQuestion += "\n解説: " + matches[3]
            }
            chatGPTResponse.Choices[0].Message.Content = formattedQuestion
            log.Println("Formatted question:", formattedQuestion)
        } else {
            log.Println("Failed to extract question, answer, and explanation")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract question, answer, and explanation"})
            return
        }

        c.JSON(http.StatusOK, chatGPTResponse)
    })

    router.Run(":3000")
}