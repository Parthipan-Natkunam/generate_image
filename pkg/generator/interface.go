package generator

import (
	"context"
)

type ImageGenerator interface {
	// Generate creates an image based on the prompt and options.
	// Returns the image data (bytes), the content type (e.g. "image/png"), and any error encountered.
	Generate(ctx context.Context, prompt string, opts ...Option) ([]byte, string, error)
	
	// Name returns the unique identifier for the provider.
	Name() string
}

type GenerateOptions struct {
	Width          int
	Height         int
	AspectRatio    string
	NegativePrompt string
	Model          string
}

// Option is a functional option for configuring GenerateOptions.
type Option func(*GenerateOptions)

func WithSize(width, height int) Option {
	return func(o *GenerateOptions) {
		o.Width = width
		o.Height = height
	}
}

func WithAspectRatio(ratio string) Option {
	return func(o *GenerateOptions) {
		o.AspectRatio = ratio
	}
}

func WithNegativePrompt(prompt string) Option {
	return func(o *GenerateOptions) {
		o.NegativePrompt = prompt
	}
}

func WithModel(model string) Option {
	return func(o *GenerateOptions) {
		o.Model = model
	}
}
