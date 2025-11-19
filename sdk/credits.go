package zaguansdk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// CreditsBalance represents the current credit balance and tier information.
type CreditsBalance struct {
	// CreditsRemaining is the number of credits remaining.
	CreditsRemaining int `json:"credits_remaining"`

	// CreditsTotal is the total credits allocated for the current period.
	CreditsTotal int `json:"credits_total,omitempty"`

	// CreditsUsed is the number of credits used in the current period.
	CreditsUsed int `json:"credits_used,omitempty"`

	// CreditsPercent is the percentage of credits remaining.
	CreditsPercent float64 `json:"credits_percent,omitempty"`

	// Tier is the user's subscription tier.
	// Examples: "free", "pro", "enterprise"
	Tier string `json:"tier"`

	// Bands are the model bands the user has access to.
	// Examples: ["A", "B", "C"]
	Bands []string `json:"bands"`

	// ResetDate is when the credits will reset (ISO 8601 format).
	ResetDate string `json:"reset_date,omitempty"`

	// StripePriceID is the Stripe price ID for the subscription.
	StripePriceID string `json:"stripe_price_id,omitempty"`

	// Warning is an optional warning message (e.g., low credits).
	Warning string `json:"warning,omitempty"`
}

// GetCreditsBalance retrieves the current credit balance and tier information.
//
// This endpoint requires authentication and returns the user's credit balance,
// tier, accessible bands, and reset date.
//
// Example:
//
//	balance, err := client.GetCreditsBalance(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Credits: %d/%d (%s tier)\n",
//		balance.CreditsRemaining, balance.CreditsTotal, balance.Tier)
func (c *Client) GetCreditsBalance(ctx context.Context, opts *RequestOptions) (*CreditsBalance, error) {
	c.log(ctx, LogLevelDebug, "getting credits balance")

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "GET",
		Path:   "/v1/credits/balance",
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var balance CreditsBalance
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &balance); err != nil {
		c.log(ctx, LogLevelError, "get credits balance request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "get credits balance request succeeded",
		"remaining", balance.CreditsRemaining,
		"tier", balance.Tier)

	return &balance, nil
}

// CreditsHistoryEntry represents a single credit usage entry.
type CreditsHistoryEntry struct {
	// ID is the unique identifier for this entry.
	ID string `json:"id"`

	// Timestamp is when the request was made (ISO 8601 format).
	Timestamp string `json:"timestamp"`

	// RequestID is the unique request identifier.
	RequestID string `json:"request_id"`

	// Model is the model used.
	Model string `json:"model"`

	// Provider is the provider name.
	Provider string `json:"provider"`

	// Band is the model band.
	Band string `json:"band"`

	// PromptTokens is the number of prompt tokens.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of completion tokens.
	CompletionTokens int `json:"completion_tokens"`

	// ReasoningTokens is the number of reasoning tokens.
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`

	// TotalTokens is the total number of tokens.
	TotalTokens int `json:"total_tokens"`

	// CreditsDebited is the number of credits charged.
	CreditsDebited int `json:"credits_debited"`

	// Cost is the cost in USD.
	Cost float64 `json:"cost,omitempty"`

	// LatencyMs is the request latency in milliseconds.
	LatencyMs int `json:"latency_ms,omitempty"`

	// Status is the request status.
	// Values: "success", "error", "rate_limited"
	Status string `json:"status"`

	// ErrorMessage is the error message (if status is "error").
	ErrorMessage string `json:"error_message,omitempty"`
}

// CreditsHistoryResponse represents the response from GET /v1/credits/history.
type CreditsHistoryResponse struct {
	// Entries is the list of credit usage entries.
	Entries []CreditsHistoryEntry `json:"entries"`

	// Total is the total number of entries available.
	Total int `json:"total,omitempty"`

	// HasMore indicates if there are more entries available.
	HasMore bool `json:"has_more,omitempty"`

	// NextCursor is the cursor for pagination.
	NextCursor string `json:"next_cursor,omitempty"`
}

// CreditsHistoryOptions contains options for querying credit history.
type CreditsHistoryOptions struct {
	// Limit is the maximum number of entries to return (default: 100).
	Limit int

	// Cursor is the pagination cursor.
	Cursor string

	// StartDate filters entries after this date (ISO 8601 format).
	StartDate string

	// EndDate filters entries before this date (ISO 8601 format).
	EndDate string

	// Model filters by model ID.
	Model string

	// Provider filters by provider name.
	Provider string

	// Band filters by band.
	Band string

	// Status filters by status.
	Status string
}

// GetCreditsHistory retrieves credit usage history with optional filtering.
//
// This endpoint supports pagination and filtering by date, model, provider, and band.
//
// Example:
//
//	history, err := client.GetCreditsHistory(ctx, &zaguansdk.CreditsHistoryOptions{
//		Limit: 50,
//		Model: "openai/gpt-4o",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range history.Entries {
//		fmt.Printf("%s: %d credits (%d tokens)\n",
//			entry.Timestamp, entry.CreditsDebited, entry.TotalTokens)
//	}
func (c *Client) GetCreditsHistory(ctx context.Context, historyOpts *CreditsHistoryOptions, opts *RequestOptions) (*CreditsHistoryResponse, error) {
	c.log(ctx, LogLevelDebug, "getting credits history")

	// Build request config
	reqCfg := internal.RequestConfig{
		Method:      "GET",
		Path:        "/v1/credits/history",
		QueryParams: make(map[string]string),
	}

	// Add query parameters from history options
	if historyOpts != nil {
		if historyOpts.Limit > 0 {
			reqCfg.QueryParams["limit"] = fmt.Sprintf("%d", historyOpts.Limit)
		}
		if historyOpts.Cursor != "" {
			reqCfg.QueryParams["cursor"] = historyOpts.Cursor
		}
		if historyOpts.StartDate != "" {
			reqCfg.QueryParams["start_date"] = historyOpts.StartDate
		}
		if historyOpts.EndDate != "" {
			reqCfg.QueryParams["end_date"] = historyOpts.EndDate
		}
		if historyOpts.Model != "" {
			reqCfg.QueryParams["model"] = historyOpts.Model
		}
		if historyOpts.Provider != "" {
			reqCfg.QueryParams["provider"] = historyOpts.Provider
		}
		if historyOpts.Band != "" {
			reqCfg.QueryParams["band"] = historyOpts.Band
		}
		if historyOpts.Status != "" {
			reqCfg.QueryParams["status"] = historyOpts.Status
		}
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var history CreditsHistoryResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &history); err != nil {
		c.log(ctx, LogLevelError, "get credits history request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "get credits history request succeeded",
		"count", len(history.Entries),
		"total", history.Total)

	return &history, nil
}

// CreditsStats represents aggregated credit statistics.
type CreditsStats struct {
	// Period is the time period for these stats.
	// Examples: "day", "week", "month", "all_time"
	Period string `json:"period"`

	// StartDate is the start of the period (ISO 8601 format).
	StartDate string `json:"start_date,omitempty"`

	// EndDate is the end of the period (ISO 8601 format).
	EndDate string `json:"end_date,omitempty"`

	// TotalCreditsUsed is the total credits used in this period.
	TotalCreditsUsed int `json:"total_credits_used"`

	// TotalRequests is the total number of requests.
	TotalRequests int `json:"total_requests"`

	// TotalTokens is the total number of tokens processed.
	TotalTokens int `json:"total_tokens"`

	// TotalCost is the total cost in USD.
	TotalCost float64 `json:"total_cost,omitempty"`

	// AverageLatencyMs is the average latency in milliseconds.
	AverageLatencyMs float64 `json:"average_latency_ms,omitempty"`

	// ByProvider contains stats broken down by provider.
	ByProvider map[string]ProviderStats `json:"by_provider,omitempty"`

	// ByModel contains stats broken down by model.
	ByModel map[string]ModelStats `json:"by_model,omitempty"`

	// ByBand contains stats broken down by band.
	ByBand map[string]BandStats `json:"by_band,omitempty"`

	// ByDay contains daily stats (if period is "month" or longer).
	ByDay []DailyStats `json:"by_day,omitempty"`
}

// ProviderStats represents statistics for a specific provider.
type ProviderStats struct {
	// Provider is the provider name.
	Provider string `json:"provider"`

	// CreditsUsed is the credits used for this provider.
	CreditsUsed int `json:"credits_used"`

	// Requests is the number of requests to this provider.
	Requests int `json:"requests"`

	// Tokens is the total tokens processed by this provider.
	Tokens int `json:"tokens"`

	// Cost is the total cost for this provider.
	Cost float64 `json:"cost,omitempty"`
}

// ModelStats represents statistics for a specific model.
type ModelStats struct {
	// Model is the model ID.
	Model string `json:"model"`

	// CreditsUsed is the credits used for this model.
	CreditsUsed int `json:"credits_used"`

	// Requests is the number of requests to this model.
	Requests int `json:"requests"`

	// Tokens is the total tokens processed by this model.
	Tokens int `json:"tokens"`

	// Cost is the total cost for this model.
	Cost float64 `json:"cost,omitempty"`
}

// BandStats represents statistics for a specific band.
type BandStats struct {
	// Band is the band identifier.
	Band string `json:"band"`

	// CreditsUsed is the credits used for this band.
	CreditsUsed int `json:"credits_used"`

	// Requests is the number of requests in this band.
	Requests int `json:"requests"`

	// Tokens is the total tokens processed in this band.
	Tokens int `json:"tokens"`
}

// DailyStats represents statistics for a single day.
type DailyStats struct {
	// Date is the date (ISO 8601 format, date only).
	Date string `json:"date"`

	// CreditsUsed is the credits used on this day.
	CreditsUsed int `json:"credits_used"`

	// Requests is the number of requests on this day.
	Requests int `json:"requests"`

	// Tokens is the total tokens processed on this day.
	Tokens int `json:"tokens"`
}

// CreditsStatsOptions contains options for querying credit statistics.
type CreditsStatsOptions struct {
	// Period is the time period to aggregate.
	// Values: "day", "week", "month", "all_time"
	Period string

	// StartDate filters stats after this date (ISO 8601 format).
	StartDate string

	// EndDate filters stats before this date (ISO 8601 format).
	EndDate string

	// GroupBy specifies how to group the stats.
	// Values: "provider", "model", "band", "day"
	GroupBy []string
}

// GetCreditsStats retrieves aggregated credit statistics.
//
// This endpoint provides detailed analytics about credit usage, including
// breakdowns by provider, model, band, and time period.
//
// Example:
//
//	stats, err := client.GetCreditsStats(ctx, &zaguansdk.CreditsStatsOptions{
//		Period: "month",
//		GroupBy: []string{"provider", "model"},
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Total credits used: %d\n", stats.TotalCreditsUsed)
//	for provider, pstats := range stats.ByProvider {
//		fmt.Printf("  %s: %d credits\n", provider, pstats.CreditsUsed)
//	}
func (c *Client) GetCreditsStats(ctx context.Context, statsOpts *CreditsStatsOptions, opts *RequestOptions) (*CreditsStats, error) {
	c.log(ctx, LogLevelDebug, "getting credits stats")

	// Build request config
	reqCfg := internal.RequestConfig{
		Method:      "GET",
		Path:        "/v1/credits/stats",
		QueryParams: make(map[string]string),
	}

	// Add query parameters from stats options
	if statsOpts != nil {
		if statsOpts.Period != "" {
			reqCfg.QueryParams["period"] = statsOpts.Period
		}
		if statsOpts.StartDate != "" {
			reqCfg.QueryParams["start_date"] = statsOpts.StartDate
		}
		if statsOpts.EndDate != "" {
			reqCfg.QueryParams["end_date"] = statsOpts.EndDate
		}
		if len(statsOpts.GroupBy) > 0 {
			reqCfg.QueryParams["group_by"] = strings.Join(statsOpts.GroupBy, ",")
		}
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var stats CreditsStats
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &stats); err != nil {
		c.log(ctx, LogLevelError, "get credits stats request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "get credits stats request succeeded",
		"total_credits", stats.TotalCreditsUsed,
		"total_requests", stats.TotalRequests)

	return &stats, nil
}

// ParseResetDate parses the reset date string into a time.Time.
func (b *CreditsBalance) ParseResetDate() (time.Time, error) {
	if b.ResetDate == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, b.ResetDate)
}

// DaysUntilReset calculates the number of days until credits reset.
func (b *CreditsBalance) DaysUntilReset() (int, error) {
	resetTime, err := b.ParseResetDate()
	if err != nil {
		return 0, err
	}
	if resetTime.IsZero() {
		return 0, nil
	}
	duration := time.Until(resetTime)
	return int(duration.Hours() / 24), nil
}

// IsLowCredits returns true if credits are running low (< 10%).
func (b *CreditsBalance) IsLowCredits() bool {
	return b.CreditsPercent < 10
}
