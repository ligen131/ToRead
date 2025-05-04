package controllers

import (
	"fmt"
	"sort"
	"strings"
	"to-read/controllers/auth"
	"to-read/model"
	"to-read/shared/llmprocessor"
	"to-read/utils/logs"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CollectionListRequest 收藏列表请求
type CollectionListRequest struct {
	Search string   `query:"search"`
	Tags   []string `query:"tags"`
}

// CollectionListItem 收藏列表项
type CollectionListItem struct {
	CollectionID uint32   `json:"collection_id"`
	Url          string   `json:"url"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	CreatedAt    int64    `json:"created_at"`
}

// CollectionListResponse 收藏列表响应
type CollectionListResponse struct {
	Collections []CollectionListItem `json:"collections"`
}

// CollectionListGET 获取收藏列表
func CollectionListGET(c echo.Context) error {
	logs.Debug("GET /collection/list")

	// 解析请求参数
	req := new(CollectionListRequest)
	
	req.Search = c.QueryParam("search")
	tagsParam := c.QueryParam("tags")
	
	if tagsParam != "" {
			req.Tags = strings.Split(tagsParam, ",")
	}

	// 获取用户ID
	claims, err := auth.GetClaimsFromHeader(c)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}
	userID := claims.UserID

	// 搜索收藏
	collections, err := model.SearchCollection(userID, req.Search, req.Tags)
	if err != nil {
		logs.Error("Failed to search collections",
			zap.Error(err),
			zap.Uint32("user_id", userID),
			zap.String("search", req.Search),
			zap.Any("tags", req.Tags))
		return ResponseInternalServerError(c, "Failed to search collections", err)
	}

	// 构建响应
	response := CollectionListResponse{
		Collections: make([]CollectionListItem, 0, len(collections)),
	}

	for _, collection := range collections {
		// 获取每个收藏关联的标签
		tags, err := model.GetCollectionTags(collection.ID)
		if err != nil {
			logs.Warn("Failed to get tags for collection",
				zap.Error(err),
				zap.Uint32("collection_id", collection.ID))
			// 继续处理其他收藏，不因为一个收藏的标签获取失败而中断整个请求
			continue
		}

		item := CollectionListItem{
			CollectionID: collection.ID,
			Url:          collection.Url,
			Type:         collection.Type,
			Title:        collection.Title,
			Description:  collection.Description,
			Tags:         tags,
			CreatedAt:    collection.CreatedAt.Unix(),
		}
		response.Collections = append(response.Collections, item)
	}

	return ResponseOK(c, response)
}

type CollectionAddRequest struct {
	Url string `json:"url"`
}

type CollectionAddResponse struct {
	CollectionID uint32   `json:"collection_id"`
	Url          string   `json:"url"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	CreateAt     int64    `json:"create_at"`
}

func CollectionAddPOST(c echo.Context) error {
	logs.Debug("POST /collection/add")

	collectionRequest := CollectionAddRequest{}
	_ok, err := Bind(c, &collectionRequest)
	if !_ok {
		return err
	}
	url := collectionRequest.Url
	if url == "" {
		return ResponseBadRequest(c, "Url is required", nil)
	}

	claims, err := auth.GetClaimsFromHeader(c)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}
	userID := claims.UserID

	// 检查URL是否已经被该用户收藏
	_, err = model.FindCollectionByUrl(userID, url)
	if err == nil {
		// URL已被收藏
		return ResponseBadRequest(c, "URL already collected", nil)
	} else if err != gorm.ErrRecordNotFound {
		// 发生其他错误
		logs.Error("Failed to check if URL is already collected", zap.Error(err), zap.String("url", url))
		return ResponseInternalServerError(c, "Failed to check if URL is already collected", err)
	}

	processor := llmprocessor.NewLLMProcessor()
	summary, err := processor.ProcessURLAuto(url)
	if err != nil {
		logs.Warn("Failed to process URL with LLM", zap.Error(err), zap.String("url", url))
		return ResponseInternalServerError(c, "Failed to process URL with LLM", err)
	}

	// 使用处理结果创建收藏
	collection, cwt, err := model.AddCollection(
		userID,
		url,
		summary.Type,
		summary.Title,
		summary.Description,
		summary.Tags,
	)

	if err != nil {
		logs.Warn("Failed to add collection", zap.Error(err), zap.String("url", url), zap.Any("collection", collection), zap.Any("cwt", cwt))
		return ResponseInternalServerError(c, "Add collection failed", err)
	}

	resp := CollectionAddResponse{
		CollectionID: collection.ID,
		Url:          collection.Url,
		Type:         collection.Type,
		Title:        collection.Title,
		Description:  collection.Description,
		Tags:         summary.Tags,
		CreateAt:     collection.CreatedAt.Unix(),
	}

	return ResponseOK(c, resp)
}

// CollectionSummaryRequest 收藏摘要请求
type CollectionSummaryRequest struct {
	Search string   `query:"search"`
	Tags   []string `query:"tags"`
}

// CollectionSummaryResponse 收藏摘要响应
type CollectionSummaryResponse struct {
	Summary string `json:"summary"`
}

// CollectionSummaryGET 获取收藏摘要
func CollectionSummaryGET(c echo.Context) error {
	logs.Debug("GET /collection/summary")

	// 解析请求参数
	req := new(CollectionListRequest)
	
	req.Search = c.QueryParam("search")
	tagsParam := c.QueryParam("tags")
	
	if tagsParam != "" {
			req.Tags = strings.Split(tagsParam, ",")
	}

	// 获取用户ID
	claims, err := auth.GetClaimsFromHeader(c)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}
	userID := claims.UserID

	// 搜索收藏
	collections, err := model.SearchCollection(userID, req.Search, req.Tags)
	if err != nil {
		logs.Error("Failed to search collections for summary",
			zap.Error(err),
			zap.Uint32("user_id", userID),
			zap.String("search", req.Search),
			zap.Any("tags", req.Tags))
		return ResponseInternalServerError(c, "Failed to search collections", err)
	}

	if len(collections) == 0 {
		return ResponseOK(c, CollectionSummaryResponse{
			Summary: "没有找到符合条件的收藏内容。",
		})
	}

	// 构建摘要文本
	summary, err := generateSummaryFromCollections(collections)
	if err != nil {
		logs.Error("Failed to generate summary from collections",
			zap.Error(err),
			zap.Int("collection_count", len(collections)))
		return ResponseInternalServerError(c, "Failed to generate summary", err)
	}

	return ResponseOK(c, CollectionSummaryResponse{
		Summary: summary,
	})
}

// 从收藏列表生成摘要
func generateSummaryFromCollections(collections []model.Collection) (string, error) {
	if len(collections) == 0 {
		return "没有找到符合条件的收藏内容。", nil
	}

	if len(collections) == 1 {
		return collections[0].Description, nil
	}

	// 创建LLM处理器
	processor := llmprocessor.NewLLMProcessor()

	// 构建输入文本
	var inputText strings.Builder

	// 收集标签信息以提供更多上下文
	allTags := make(map[string]int)
	for _, collection := range collections {
		tags, err := model.GetCollectionTags(collection.ID)
		if err == nil {
			for _, tag := range tags {
				allTags[tag]++
			}
		}
	}

	// 找出最常见的标签
	type tagFreq struct {
		tag   string
		count int
	}
	tagFreqs := make([]tagFreq, 0, len(allTags))
	for tag, count := range allTags {
		tagFreqs = append(tagFreqs, tagFreq{tag, count})
	}
	sort.Slice(tagFreqs, func(i, j int) bool {
		return tagFreqs[i].count > tagFreqs[j].count
	})

	// 取前5个最常见的标签
	topTags := []string{}
	for i, tf := range tagFreqs {
		if i >= 5 {
			break
		}
		topTags = append(topTags, tf.tag)
	}

	// 构建详细的输入文本
	inputText.WriteString("## 收藏内容分析任务\n\n")
	inputText.WriteString(fmt.Sprintf("请分析以下%d篇收藏内容，生成一份综合报告。请直接给出报告内容，无需说出你的思路等等。\n\n", len(collections)))

	if len(topTags) > 0 {
		inputText.WriteString("主要标签: " + strings.Join(topTags, ", ") + "\n\n")
	}

	inputText.WriteString("## 收藏内容列表\n\n")

	for i, collection := range collections {
		inputText.WriteString(fmt.Sprintf("### 文章 %d\n", i+1))
		inputText.WriteString(fmt.Sprintf("标题: %s\n", collection.Title))
		inputText.WriteString(fmt.Sprintf("类型: %s\n", collection.Type))

		// 获取标签
		tags, _ := model.GetCollectionTags(collection.ID)
		if len(tags) > 0 {
			inputText.WriteString(fmt.Sprintf("标签: %s\n", strings.Join(tags, ", ")))
		}

		inputText.WriteString(fmt.Sprintf("摘要: %s\n\n", collection.Description))

		// 限制输入长度，避免超过LLM的上下文长度限制
		if i >= 9 && len(collections) > 10 {
			inputText.WriteString(fmt.Sprintf("... 以及其他 %d 篇收藏内容\n", len(collections)-10))
			break
		}
	}

	// 构造一个特殊的虚拟URL，表示这是收藏摘要请求
	virtualURL := fmt.Sprintf("collection-summary://%d-articles", len(collections))

	// 使用专门的摘要配置
	summaryConfig := llmprocessor.ProcessorConfig{
		Enabled:     true,
		APIEndpoint: processor.GetConfig().TextProcessor.APIEndpoint,
		APIKey:      processor.GetConfig().TextProcessor.APIKey,
		Model:       processor.GetConfig().TextProcessor.Model,
		// 使用专门的摘要提示词
		Prompt:      getSummaryPrompt(),
		MaxTokens:   1500,
		Temperature: 1.0,
	}

	// 调用LLM处理器获取摘要
	summary, err := processor.GetContentSummaryFromText(inputText.String(), virtualURL, summaryConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get summary from LLM: %w", err)
	}

	return summary.Description, nil
}

// getSummaryPrompt 获取专门的摘要生成提示词
func getSummaryPrompt() string {
	return `你是一位专业的内容分析师和知识管理专家。你的任务是分析用户收藏的一组文章，并生成一份全面而有见解的综合报告。

请遵循以下指南:
1. 分析所有文章的主题和内容，找出共同的主题和关联
2. 识别关键的见解、观点和有价值的信息
3. 组织信息为一个连贯的摘要，而不是简单地列出每篇文章的内容
4. 如果文章之间有关联或互补的信息，请指出这些关联
5. 如果文章包含相互矛盾的观点，请客观地指出这些差异
6. 总结可能的行动要点或进一步探索的方向
7. 使用简体中文撰写摘要，语言应当清晰、专业且易于理解
8. 摘要应当是完整的段落，而不是要点列表

你的分析应该帮助用户理解这些收藏内容的整体价值和意义，而不仅仅是各个部分的简单总结。`
}

// CollectionTagResponse 收藏标签响应
type CollectionTagResponse struct {
	Tags []string `json:"tags"`
}

// CollectionTagGET 获取收藏标签
func CollectionTagGET(c echo.Context) error {
	logs.Debug("GET /collection/tag")

	// 获取用户ID
	claims, err := auth.GetClaimsFromHeader(c)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}
	userID := claims.UserID

	// 获取用户的所有标签
	tags, err := model.GetTagList(userID)
	if err != nil {
		logs.Error("Failed to get tag list", zap.Error(err), zap.Uint32("user_id", userID))
		return ResponseInternalServerError(c, "Failed to get tag list", err)
	}

	// 构建响应
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return ResponseOK(c, CollectionTagResponse{
		Tags: tagNames,
	})
}
