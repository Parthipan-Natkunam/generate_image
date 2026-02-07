package nanobanana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Parthipan-Natkunam/generate_image/pkg/generator"
)

const (
	defaultEndpoint = "https://api.nanobanana.im/v1/generate" // Change this to actual endpoint from the env
	providerName    = "nano-banana-pro"
)

type Provider struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

func New(apiKey string, opts ...ProviderOption) *Provider {
	p := &Provider{
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 60 * time.Second},
		endpoint: defaultEndpoint,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type ProviderOption func(*Provider)

func WithEndpoint(url string) ProviderOption {
	return func(p *Provider) {
		p.endpoint = url
	}
}

func WithClient(client *http.Client) ProviderOption {
	return func(p *Provider) {
		p.client = client
	}
}

func (p *Provider) Name() string {
	return providerName
}

// Generate sends a request to the Nano Banana API.
func (p *Provider) Generate(ctx context.Context, prompt string, opts ...generator.Option) ([]byte, string, error) {
	genOpts := &generator.GenerateOptions{
		Width:  1024,
		Height: 1024,
		Model:  "nano-banana-pro-v1", // Change this to an actual model
	}
	for _, opt := range opts {
		opt(genOpts)
	}

	reqBody := map[string]interface{}{
		"prompt":          prompt,
		"negative_prompt": genOpts.NegativePrompt,
		"width":           genOpts.Width,
		"height":          genOpts.Height,
		"model":           genOpts.Model,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "img-gen-cli/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("api returned error %d: %s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, contentType, nil
}
