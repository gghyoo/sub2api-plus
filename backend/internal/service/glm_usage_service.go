package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
)

// GLMModelSummary represents a model's total usage from the monitoring API.
type GLMModelSummary struct {
	ModelName   string `json:"modelName"`
	TotalTokens int64  `json:"totalTokens"`
	SortOrder   int    `json:"sortOrder"`
}

// GLMToolSummary represents a tool's total usage from the monitoring API.
type GLMToolSummary struct {
	ToolCode        string `json:"toolCode"`
	ToolName        string `json:"toolName"`
	TotalUsageCount int64  `json:"totalUsageCount"`
	SortOrder       int    `json:"sortOrder"`
}

// GLMQuotaLimit represents a quota limit entry from the monitoring API.
type GLMQuotaLimit struct {
	Type         string  `json:"type"`
	Unit         int     `json:"unit"`
	Number       int     `json:"number"`
	Percentage   float64 `json:"percentage"`
	Usage        int64   `json:"usage,omitempty"`
	CurrentValue int64   `json:"currentValue,omitempty"`
	Remaining    int64   `json:"remaining,omitempty"`
	NextResetMs  int64   `json:"nextResetTime,omitempty"`
	UsageDetails []struct {
		ModelCode string `json:"modelCode"`
		Usage     int64  `json:"usage"`
	} `json:"usageDetails,omitempty"`
}

// GLMUsageResponse combines all monitoring API responses.
type GLMUsageResponse struct {
	TotalModelCalls int64             `json:"total_model_calls"`
	TotalTokens     int64             `json:"total_tokens"`
	Models          []GLMModelSummary `json:"models"`
	Tools           []GLMToolSummary  `json:"tools"`
	QuotaLevel      string            `json:"quota_level"`
	QuotaLimits     []GLMQuotaLimit   `json:"quota_limits"`
}

// FetchGLMUsage fetches GLM Coding Plan usage from ZhiPu monitoring API.
// baseDomain should be scheme + host only (e.g. "https://open.bigmodel.cn").
func FetchGLMUsage(ctx context.Context, baseDomain, apiKey string) (*GLMUsageResponse, error) {
	client, err := httpclient.GetClient(httpclient.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}

	now := time.Now()
	startTime := formatGLMTime(time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(), 0, 0, 0, now.Location()))
	endTime := formatGLMTime(time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 59, 0, now.Location()))
	queryParams := "?startTime=" + url.QueryEscape(startTime) + "&endTime=" + url.QueryEscape(endTime)

	result := &GLMUsageResponse{}

	// Model usage
	modelData, err := glmGet(client, ctx, baseDomain+"/api/monitor/usage/model-usage"+queryParams, apiKey)
	if err != nil {
		return nil, fmt.Errorf("fetch model usage: %w", err)
	}
	if modelData != nil {
		parseModelUsage(modelData, result)
	}

	// Tool usage
	toolData, err := glmGet(client, ctx, baseDomain+"/api/monitor/usage/tool-usage"+queryParams, apiKey)
	if err != nil {
		return nil, fmt.Errorf("fetch tool usage: %w", err)
	}
	if toolData != nil {
		parseToolUsage(toolData, result)
	}

	// Quota limit (no time params)
	quotaData, err := glmGet(client, ctx, baseDomain+"/api/monitor/usage/quota/limit", apiKey)
	if err != nil {
		return nil, fmt.Errorf("fetch quota limit: %w", err)
	}
	if quotaData != nil {
		parseQuotaLimit(quotaData, result)
	}

	return result, nil
}

// parseModelUsage extracts model summaries from the model-usage API response.
// Response format: {"x_time":[...], "modelCallCount":[...], "tokensUsage":[...],
//
//	"totalUsage":{"totalModelCallCount":N, "totalTokensUsage":N, "modelSummaryList":[...]},
//	"modelSummaryList":[...]}
func parseModelUsage(data json.RawMessage, result *GLMUsageResponse) {
	var resp struct {
		TotalUsage struct {
			TotalModelCalls int64             `json:"totalModelCallCount"`
			TotalTokens     int64             `json:"totalTokensUsage"`
			Models          []GLMModelSummary `json:"modelSummaryList"`
		} `json:"totalUsage"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}
	result.TotalModelCalls = resp.TotalUsage.TotalModelCalls
	result.TotalTokens = resp.TotalUsage.TotalTokens
	result.Models = resp.TotalUsage.Models
}

// parseToolUsage extracts tool summaries from the tool-usage API response.
// Response format: {"x_time":[...], "totalUsage":{"toolSummaryList":[...]}, ...}
func parseToolUsage(data json.RawMessage, result *GLMUsageResponse) {
	var resp struct {
		TotalUsage struct {
			Tools []GLMToolSummary `json:"toolSummaryList"`
		} `json:"totalUsage"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}
	result.Tools = resp.TotalUsage.Tools
}

// parseQuotaLimit extracts quota limits from the quota/limit API response.
// Also merges MCP usageDetails into the tools list if not already present.
func parseQuotaLimit(data json.RawMessage, result *GLMUsageResponse) {
	var resp struct {
		Limits []GLMQuotaLimit `json:"limits"`
		Level  string          `json:"level"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}
	result.QuotaLimits = resp.Limits
	result.QuotaLevel = resp.Level

	// Merge MCP usageDetails into tools (e.g. zread may only appear here)
	existing := make(map[string]bool)
	for _, t := range result.Tools {
		existing[t.ToolCode] = true
	}
	for _, limit := range resp.Limits {
		if limit.Type == "TIME_LIMIT" {
			for idx, detail := range limit.UsageDetails {
				if !existing[detail.ModelCode] {
					result.Tools = append(result.Tools, GLMToolSummary{
						ToolCode:        detail.ModelCode,
						ToolName:        detail.ModelCode,
						TotalUsageCount: detail.Usage,
						SortOrder:       len(result.Tools) + idx + 1,
					})
					existing[detail.ModelCode] = true
				}
			}
		}
	}
}

// glmGet makes a GET request to the monitoring API and returns the "data" field.
func glmGet(client *http.Client, ctx context.Context, apiURL, apiKey string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Accept-Language", "en-US,en")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return envelope.Data, nil
}

// formatGLMTime formats a time as "2006-01-02 15:04:05".
func formatGLMTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// IsBigModelBaseURL checks if a base URL points to a ZhiPu/bigmodel endpoint.
func IsBigModelBaseURL(baseURL string) bool {
	return strings.Contains(strings.ToLower(baseURL), "bigmodel")
}

// IsMiniMaxBaseURL checks if a base URL points to a MiniMax endpoint.
func IsMiniMaxBaseURL(baseURL string) bool {
	return strings.Contains(strings.ToLower(baseURL), "minimax")
}

// IsKimiBaseURL checks if a base URL points to a Kimi endpoint.
func IsKimiBaseURL(baseURL string) bool {
	lower := strings.ToLower(baseURL)
	return strings.Contains(lower, "kimi.com") || strings.Contains(lower, "moonshot")
}

// MiniMaxModelRemain represents a model's remaining quota from the MiniMax API.
// NOTE: current_interval_usage_count and current_weekly_usage_count are REMAINING counts, not used.
type MiniMaxModelRemain struct {
	ModelName                 string `json:"model_name"`
	CurrentIntervalTotalCount int64  `json:"current_interval_total_count"` // 5h window total quota
	CurrentIntervalUsageCount int64  `json:"current_interval_usage_count"` // 5h window REMAINING
	StartTime                 int64  `json:"start_time"`
	EndTime                   int64  `json:"end_time"`                   // 5h window reset time (ms epoch)
	RemainsTime               int64  `json:"remains_time"`               // ms remaining until 5h reset
	CurrentWeeklyTotalCount   int64  `json:"current_weekly_total_count"` // weekly total quota
	CurrentWeeklyUsageCount   int64  `json:"current_weekly_usage_count"` // weekly REMAINING
	WeeklyStartTime           int64  `json:"weekly_start_time"`
	WeeklyEndTime             int64  `json:"weekly_end_time"`     // weekly reset time (ms epoch)
	WeeklyRemainsTime         int64  `json:"weekly_remains_time"` // ms remaining until weekly reset
}

// MiniMaxUsageResponse combines all MiniMax Coding Plan usage data.
type MiniMaxUsageResponse struct {
	Models []MiniMaxModelRemain `json:"models"`
}

// KimiWindow represents a rate-limit window definition.
type KimiWindow struct {
	Duration int    `json:"duration"`
	TimeUnit string `json:"timeUnit"`
}

// KimiDetail represents the usage detail for a window.
type KimiDetail struct {
	Limit     string `json:"limit"`
	Remaining string `json:"remaining"`
	ResetTime string `json:"resetTime"`
}

// KimiLimit represents a single limit entry from the Kimi usage API.
type KimiLimit struct {
	Window KimiWindow `json:"window"`
	Detail KimiDetail `json:"detail"`
}

// KimiUsageResponse represents the response from Kimi Coding Plan usage API.
type KimiUsageResponse struct {
	Usage struct {
		Limit     string `json:"limit"`
		Used      string `json:"used"`
		Remaining string `json:"remaining"`
		ResetTime string `json:"resetTime"`
	} `json:"usage"`
	Limits []KimiLimit `json:"limits"`
}

// CodingPlanUsageResponse is a unified response for GLM, MiniMax, and Kimi coding plan usage.
type CodingPlanUsageResponse struct {
	Platform string                `json:"platform"` // "glm", "minimax", or "kimi"
	GLM      *GLMUsageResponse     `json:"glm,omitempty"`
	MiniMax  *MiniMaxUsageResponse `json:"minimax,omitempty"`
	Kimi     *KimiUsageResponse    `json:"kimi,omitempty"`
}

// FetchMiniMaxUsage fetches MiniMax Coding Plan usage from the MiniMax API.
func FetchMiniMaxUsage(ctx context.Context, apiKey string) (*MiniMaxUsageResponse, error) {
	client, err := httpclient.GetClient(httpclient.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}

	apiURL := "https://www.minimaxi.com/v1/api/openplatform/coding_plan/remains"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	// API returns {"model_remains": [...]} directly (no nested data field)
	var apiResp struct {
		ModelRemains []MiniMaxModelRemain `json:"model_remains"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &MiniMaxUsageResponse{
		Models: apiResp.ModelRemains,
	}, nil
}

// FetchKimiUsage fetches Kimi Coding Plan usage from the Kimi API.
func FetchKimiUsage(ctx context.Context, apiKey string) (*KimiUsageResponse, error) {
	client, err := httpclient.GetClient(httpclient.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}

	apiURL := "https://api.kimi.com/coding/v1/usages"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	var result KimiUsageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}

// ExtractBaseDomain extracts scheme://host from a URL string.
func ExtractBaseDomain(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}
	return parsed.Scheme + "://" + parsed.Host, nil
}
