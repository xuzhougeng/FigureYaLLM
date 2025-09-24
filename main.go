package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Module represents a FigureYa module
type Module struct {
	Module      string `json:"module"`
	需求描述        string `json:"需求描述"`
	实用场景        string `json:"实用场景"`
	图片类型        string `json:"图片类型"`
	LLMStatus   string `json:"llm_status"`
}

// DocumentData represents the JSON structure
type DocumentData struct {
	GeneratedAt string   `json:"generated_at"`
	Modules     []Module `json:"modules"`
}

// RecommendationRequest represents the incoming request
type RecommendationRequest struct {
	Query string `json:"query" binding:"required"`
}

// RecommendationResponse represents the response
type RecommendationResponse struct {
	Query         string   `json:"query"`
	Recommendations []ModuleRecommendation `json:"recommendations"`
	Explanation   string   `json:"explanation"`
}

// ModuleRecommendation represents a recommended module
type ModuleRecommendation struct {
	Module      string  `json:"module"`
	Description string  `json:"description"`
	UseCase     string  `json:"useCase"`
	ChartType   string  `json:"chartType"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
}

// OpenAIRequest for LLM API calls
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Message for chat completion
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse for LLM API response
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

// RecommendationService handles the recommendation logic
type RecommendationService struct {
	modules    []Module
	openaiURL  string
	openaiKey  string
	model      string
}

func main() {
	// Load configuration
	openaiURL := getEnv("BASE_URL", getEnv("OPENAI_URL", "https://api.openai.com/v1/chat/completions"))
	if !strings.HasSuffix(openaiURL, "/chat/completions") && !strings.Contains(openaiURL, "/chat/completions") {
		if strings.HasSuffix(openaiURL, "/") {
			openaiURL = openaiURL + "v1/chat/completions"
		} else {
			openaiURL = openaiURL + "/v1/chat/completions"
		}
	}
	openaiKey := getEnv("OPENAI_API_KEY", "")
	model := getEnv("MODEL", "gpt-3.5-turbo")

	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY is required")
	}

	// Load modules data
	modules, err := loadModules("figureya_docs_llm.json")
	if err != nil {
		log.Fatal("Failed to load modules:", err)
	}

	// Create recommendation service
	service := &RecommendationService{
		modules:   modules,
		openaiURL: openaiURL,
		openaiKey: openaiKey,
		model:     model,
	}

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Routes
	r.POST("/recommend", service.handleRecommendation)
	r.GET("/modules", service.handleListModules)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "modules_loaded": len(service.modules)})
	})

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s with %d modules loaded", port, len(modules))
	log.Fatal(r.Run(":" + port))
}

func (s *RecommendationService) handleRecommendation(c *gin.Context) {
	var req RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	recommendations, explanation, err := s.getRecommendations(req.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := RecommendationResponse{
		Query:           req.Query,
		Recommendations: recommendations,
		Explanation:     explanation,
	}

	c.JSON(200, response)
}

func (s *RecommendationService) handleListModules(c *gin.Context) {
	c.JSON(200, gin.H{
		"total":   len(s.modules),
		"modules": s.modules,
	})
}

func (s *RecommendationService) getRecommendations(query string) ([]ModuleRecommendation, string, error) {
	// Build context for LLM
	context := s.buildModulesContext()

	prompt := fmt.Sprintf(`根据用户查询，从提供的FigureYa模块中推荐最相关的3-5个模块。

用户查询: %s

可用模块信息:
%s

请分析用户需求，推荐最相关的模块，并解释推荐理由。

请按以下JSON格式回复：
{
  "recommendations": [
    {
      "module": "模块名",
      "score": 0.95,
      "reason": "推荐理由"
    }
  ],
  "explanation": "整体推荐说明"
}`, query, context)

	// Call LLM API
	response, err := s.callLLM(prompt)
	if err != nil {
		return nil, "", fmt.Errorf("LLM API call failed: %v", err)
	}

	// Parse LLM response
	return s.parseLLMResponse(response)
}

func (s *RecommendationService) buildModulesContext() string {
	var builder strings.Builder
	for _, module := range s.modules {
		builder.WriteString(fmt.Sprintf("模块: %s\n", module.Module))
		builder.WriteString(fmt.Sprintf("需求描述: %s\n", module.需求描述))
		builder.WriteString(fmt.Sprintf("实用场景: %s\n", module.实用场景))
		builder.WriteString(fmt.Sprintf("图片类型: %s\n", module.图片类型))
		builder.WriteString("\n")
	}
	return builder.String()
}

func (s *RecommendationService) callLLM(prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model:       s.model,
		Temperature: 0.3,
		MaxTokens:   2000,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的数据可视化助手，能够根据用户需求推荐合适的FigureYa模块。请仔细分析用户查询并给出准确的推荐。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.openaiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openaiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", err
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func (s *RecommendationService) parseLLMResponse(response string) ([]ModuleRecommendation, string, error) {
	// Try to extract JSON from response
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 {
		return nil, "", fmt.Errorf("invalid JSON response from LLM")
	}

	jsonStr := response[start : end+1]

	var parsed struct {
		Recommendations []struct {
			Module string  `json:"module"`
			Score  float64 `json:"score"`
			Reason string  `json:"reason"`
		} `json:"recommendations"`
		Explanation string `json:"explanation"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, "", fmt.Errorf("failed to parse LLM response: %v", err)
	}

	// Build recommendations with full module info
	var recommendations []ModuleRecommendation
	moduleMap := make(map[string]Module)
	for _, m := range s.modules {
		moduleMap[m.Module] = m
	}

	for _, rec := range parsed.Recommendations {
		if module, exists := moduleMap[rec.Module]; exists {
			recommendations = append(recommendations, ModuleRecommendation{
				Module:      module.Module,
				Description: module.需求描述,
				UseCase:     module.实用场景,
				ChartType:   module.图片类型,
				Score:       rec.Score,
				Reason:      rec.Reason,
			})
		}
	}

	return recommendations, parsed.Explanation, nil
}

func loadModules(filename string) ([]Module, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	modules := data["modules"].([]interface{})
	var validModules []Module

	for _, moduleInterface := range modules {
		moduleMap := moduleInterface.(map[string]interface{})

		// Check if status is "ok"
		if status, exists := moduleMap["llm_status"]; !exists || status != "ok" {
			continue
		}

		// Extract fields with safe type assertions
		module := Module{
			Module:    getString(moduleMap, "module"),
			需求描述:      getString(moduleMap, "需求描述"),
			实用场景:      getString(moduleMap, "实用场景"),
			图片类型:      getString(moduleMap, "图片类型"),
			LLMStatus: getString(moduleMap, "llm_status"),
		}

		validModules = append(validModules, module)
	}

	return validModules, nil
}

// Helper function to safely extract string from map
func getString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}