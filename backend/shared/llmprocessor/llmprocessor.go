package llmprocessor

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"to-read/utils/logs"

	"go.uber.org/zap"
)

type LLMConfig struct {
	TextProcessor  ProcessorConfig `yaml:"textProcessor"`
	ImageProcessor ProcessorConfig `yaml:"imageProcessor"`
	VideoProcessor ProcessorConfig `yaml:"videoProcessor"`
}

type ProcessorConfig struct {
	Enabled     bool    `yaml:"enabled"`
	APIEndpoint string  `yaml:"apiEndpoint"`
	APIKey      string  `yaml:"apiKey"`
	Model       string  `yaml:"model"`
	Prompt      string  `yaml:"prompt"`
	MaxTokens   int     `yaml:"maxTokens"`
	Temperature float64 `yaml:"temperature"`
}

// ContentType 内容类型
type ContentType string

const (
	TextType  ContentType = "text"
	ImageType ContentType = "image"
	VideoType ContentType = "video"
)

// ContentSummary 内容摘要
type ContentSummary struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// Processor LLM处理器接口
type Processor interface {
	ProcessURL(url string) (ContentSummary, error)
}

// LLMProcessor 大语言模型处理器
type LLMProcessor struct {
	config LLMConfig
	client *http.Client
}

var config LLMConfig

func InitLLMProcessor(cfg LLMConfig) error {
	config = cfg

	return nil
}

// NewLLMProcessor 创建新的LLM处理器
func NewLLMProcessor() *LLMProcessor {
	return &LLMProcessor{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessURL 处理URL并返回内容摘要
func (p *LLMProcessor) ProcessURL(url string, contentType ContentType) (ContentSummary, error) {
	// 根据内容类型选择处理器配置
	var processorConfig ProcessorConfig
	switch contentType {
	case TextType:
		processorConfig = p.config.TextProcessor
	case ImageType:
		processorConfig = p.config.ImageProcessor
	case VideoType:
		processorConfig = p.config.VideoProcessor
	default:
		return ContentSummary{}, errors.New("unsupported content type")
	}

	// 检查处理器是否启用
	if !processorConfig.Enabled {
		return ContentSummary{}, fmt.Errorf("%s processor is not enabled", contentType)
	}

	// 获取网页内容
	content, err := p.fetchURLContent(url)
	if err != nil {
		return ContentSummary{}, fmt.Errorf("failed to fetch URL content: %w", err)
	}

	logs.Debug("Fetching URL content", zap.String("url", url), zap.String("content", content))

	// 调用LLM API处理内容
	summary, err := p.callLLMAPI(content, url, processorConfig)
	if err != nil {
		return ContentSummary{}, fmt.Errorf("failed to process content with LLM: %w", err)
	}

	// 设置内容类型
	summary.Type = string(contentType)

	logs.Debug("Processing URL content", zap.String("url", url), zap.Any("summary", summary))

	return summary, nil
}

// fetchURLContent 获取URL内容
func (p *LLMProcessor) fetchURLContent(url string) (string, error) {
	// 对于文本格式的网页，先通过 r.jina.ai 转换成 Markdown 格式
	if strings.HasPrefix(url, "http") && !p.isImageURL(url) {
		jinaURL := fmt.Sprintf("https://r.jina.ai/%s", url)
		resp, err := p.client.Get(jinaURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return string(body), nil
		}
		// 如果 jina 转换失败，回退到直接获取
		logs.Warn("Failed to convert via jina.ai, falling back to direct fetch", zap.String("url", url))
	}

	// 对于图片URL，转换成base64
	if p.isImageURL(url) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
		req.Header.Set("Sec-Fetch-Mode", "no-cors")
		req.Header.Set("Sec-Fetch-Dest", "image")
		// req.Header.Set("Referer", "https://www.google.com/")

		// 发送请求
		resp, err := p.client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
			logs.Info("Fetching image content failed", zap.String("url", url), zap.String("body", string(data)))
			return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// 读取图片内容
		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// 这里返回特殊格式，表示这是一个图片内容
		return fmt.Sprintf("__IMAGE_BASE64__:%s", p.encodeToBase64(imageData, resp.Header.Get("Content-Type"))), nil
	}

	// 默认处理方式
	resp, err := p.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// isImageURL 判断URL是否为图片
func (p *LLMProcessor) isImageURL(url string) bool {
	lowerURL := strings.ToLower(url)
	return strings.HasSuffix(lowerURL, ".jpg") ||
		strings.HasSuffix(lowerURL, ".jpeg") ||
		strings.HasSuffix(lowerURL, ".png") ||
		strings.HasSuffix(lowerURL, ".gif") ||
		strings.HasSuffix(lowerURL, ".webp")
}

// encodeToBase64 将图片数据编码为base64
func (p *LLMProcessor) encodeToBase64(data []byte, contentType string) string {
	if contentType == "" {
		contentType = "image/png" // 默认类型
	}
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data))
}

// OpenAIRequest OpenAI API请求
type OpenAIRequest struct {
	Model        string           `json:"model"`
	Messages     []OpenAIMessage  `json:"messages"`
	Functions    []OpenAIFunction `json:"functions,omitempty"`
	FunctionCall string           `json:"function_call,omitempty"`
	MaxTokens    int              `json:"max_tokens,omitempty"`
	Temperature  float64          `json:"temperature,omitempty"`
}

// OpenAIMessage OpenAI消息
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIFunction OpenAI函数
type OpenAIFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// OpenAIResponse OpenAI API响应
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role         string `json:"role"`
			Content      string `json:"content"`
			FunctionCall struct {
				Name      string `json:"name"`
				Arguments string `json:"arguments"`
			} `json:"function_call,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// callLLMAPI 调用LLM API
func (p *LLMProcessor) callLLMAPI(content, url string, config ProcessorConfig) (ContentSummary, error) {
	// 检查内容是否为图片
	if strings.HasPrefix(content, "__IMAGE_BASE64__:") {
		// 处理图片内容
		imageBase64 := strings.TrimPrefix(content, "__IMAGE_BASE64__:")

		// 第一步：让LLM详细解释图片内容
		imageDescription, err := p.getImageDescription(imageBase64, url, config)
		if err != nil {
			return ContentSummary{}, fmt.Errorf("failed to get image description: %w", err)
		}

		logs.Debug("Image description", zap.String("description", imageDescription))

		// 第二步：将图片描述作为文本再次请求LLM，获取结构化摘要
		return p.getContentSummaryFromText(imageDescription, url, config, false)
	} else {
		// 直接处理文本内容
		return p.getContentSummaryFromText(content, url, config, false)
	}
}

// getImageDescription 获取图片的详细描述
func (p *LLMProcessor) getImageDescription(imageBase64, url string, config ProcessorConfig) (string, error) {
	// 构建系统提示
	systemPrompt := "你是一个专业的图片分析专家。请详细描述这张图片的内容，包括主要对象、场景、文字、颜色、风格等方面。尽可能全面和具体地描述，不要遗漏任何可见的细节。"

	// 构建包含图片的请求
	imageMessage := map[string]interface{}{
		"role": "user",
		"content": []map[string]interface{}{
			{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": imageBase64,
				},
			},
			{
				"type": "text",
				"text": "请详细描述这张图片的内容。这张图片来自URL: " + url,
			},
		},
	}

	// 构建请求体
	requestMap := map[string]interface{}{
		"model": config.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "system",
				"content": systemPrompt,
			},
			imageMessage,
		},
		"max_tokens":  1000,
		"temperature": 0.5,
	}

	requestBody, err := json.Marshal(requestMap)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.APIEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	logs.Debug("Sending image description request to OpenAI API", zap.String("request", string(requestBody)))

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 记录响应内容到日志
	logs.Debug("API response body for image description", zap.String("body", string(bodyBytes)))

	// 解析响应
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", err
	}

	// 检查是否有响应
	if len(apiResponse.Choices) == 0 {
		return "", errors.New("no response from LLM API for image description")
	}

	return apiResponse.Choices[0].Message.Content, nil
}

// getContentSummaryFromText 从文本内容获取结构化摘要
func (p *LLMProcessor) getContentSummaryFromText(content, url string, config ProcessorConfig, isSummaryRequest bool) (ContentSummary, error) {
	var messages []OpenAIMessage

	// 系统提示
	messages = append(messages, OpenAIMessage{
		Role:    "system",
		Content: config.Prompt,
	})

	// 用户消息
	messages = append(messages, OpenAIMessage{
		Role:    "user",
		Content: fmt.Sprintf("URL: %s\n\nContent: %s", url, content),
	})

	enable_function_call := "auto"
	if isSummaryRequest {
		enable_function_call = "none"
	}

	// 构建OpenAI API请求
	request := OpenAIRequest{
		Model:        config.Model,
		Messages:     messages,
		Functions:    getFunctionDefinition(),
		FunctionCall: enable_function_call,
		MaxTokens:    config.MaxTokens,
		Temperature:  config.Temperature,
	}

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return ContentSummary{}, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.APIEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return ContentSummary{}, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	logs.Debug("Sending content summary request to OpenAI API", zap.String("request", string(requestBody)))

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return ContentSummary{}, err
	}
	defer resp.Body.Close()

	if isSummaryRequest {
		return processTextResponse(resp)
	}
	return processAPIResponse(resp)
}

// getFunctionDefinition 获取函数定义
func getFunctionDefinition() []OpenAIFunction {
	return []OpenAIFunction{
		{
			Name:        "extract_content_summary",
			Description: "提取网页内容的摘要信息，并以简体中文返回结果",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "内容的标题（使用简体中文）",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "内容的简洁摘要（使用简体中文）",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "与内容相关的标签（使用简体中文）",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"required": []string{"title", "description", "tags"},
			},
		},
	}
}

// processAPIResponse 处理API响应
func processAPIResponse(resp *http.Response) (ContentSummary, error) {
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return ContentSummary{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContentSummary{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// 记录响应内容到日志
	logs.Debug("API response body", zap.String("body", string(bodyBytes)))

	// 重新创建一个新的io.ReadCloser供后续代码使用
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 解析响应
	var apiResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return ContentSummary{}, err
	}

	// 检查是否有响应
	if len(apiResponse.Choices) == 0 {
		return ContentSummary{}, errors.New("no response from LLM API")
	}

	// 解析函数调用结果
	functionCall := apiResponse.Choices[0].Message.FunctionCall
	if functionCall.Name != "extract_content_summary" {
		return ContentSummary{}, fmt.Errorf("unexpected function call: %s", functionCall.Name)
	}

	// 解析函数参数
	var summary ContentSummary
	if err := json.Unmarshal([]byte(functionCall.Arguments), &summary); err != nil {
		return ContentSummary{}, err
	}

	return summary, nil
}

// DetectContentType 检测URL的内容类型
func (p *LLMProcessor) DetectContentType(url string) (ContentType, error) {
	// 简单的基于URL后缀的检测
	lowerURL := strings.ToLower(url)

	// 检测图片
	if strings.HasSuffix(lowerURL, ".jpg") ||
		strings.HasSuffix(lowerURL, ".jpeg") ||
		strings.HasSuffix(lowerURL, ".png") ||
		strings.HasSuffix(lowerURL, ".gif") ||
		strings.HasSuffix(lowerURL, ".webp") {
		return ImageType, nil
	}

	// 检测视频
	if strings.HasSuffix(lowerURL, ".mp4") ||
		strings.HasSuffix(lowerURL, ".avi") ||
		strings.HasSuffix(lowerURL, ".mov") ||
		strings.HasSuffix(lowerURL, ".webm") ||
		strings.Contains(lowerURL, "youtube.com/watch") ||
		strings.Contains(lowerURL, "youtu.be/") {
		return VideoType, nil
	}

	// 默认为文本
	return TextType, nil
}

// ProcessURLAuto 自动检测内容类型并处理URL
func (p *LLMProcessor) ProcessURLAuto(url string) (ContentSummary, error) {
	contentType, err := p.DetectContentType(url)
	if err != nil {
		return ContentSummary{}, err
	}

	logs.Info("Detected content type", zap.String("url", url), zap.String("type", string(contentType)))
	return p.ProcessURL(url, contentType)
}

// GetContentSummaryFromText 从文本内容获取结构化摘要（公开版本）
func (p *LLMProcessor) GetContentSummaryFromText(content, url string, config ProcessorConfig) (ContentSummary, error) {
	return p.getContentSummaryFromText(content, url, config, true)
}

// GetConfig 获取LLM处理器配置
func (p *LLMProcessor) GetConfig() LLMConfig {
	return p.config
}

// processTextResponse 处理文本响应（用于摘要请求）
func processTextResponse(resp *http.Response) (ContentSummary, error) {
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return ContentSummary{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContentSummary{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// 记录响应内容到日志
	logs.Debug("API response body for text response", zap.String("body", string(bodyBytes)))

	// 解析响应
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return ContentSummary{}, err
	}

	// 检查是否有响应
	if len(apiResponse.Choices) == 0 {
		return ContentSummary{}, errors.New("no response from LLM API")
	}

	// 创建摘要对象
	summary := ContentSummary{
		Type:        "summary",
		Title:       "收藏内容综合摘要",
		Description: apiResponse.Choices[0].Message.Content,
		Tags:        []string{"摘要"},
	}

	return summary, nil
}
