package watermark

import (
	"errors"
	"fmt"
)

// Position represents the placement position of a watermark on the base image
type Position string

// Position constants for watermark placement
const (
	PositionTopLeft      Position = "top-left"
	PositionTopCenter    Position = "top-center"
	PositionTopRight     Position = "top-right"
	PositionLeftCenter   Position = "left-center"
	PositionCenter       Position = "center"
	PositionRightCenter  Position = "right-center"
	PositionBottomLeft   Position = "bottom-left"
	PositionBottomCenter Position = "bottom-center"
	PositionBottomRight  Position = "bottom-right"
)

// Config contains all configuration parameters for watermarking
type Config struct {
	// Type selection (one must be set)
	Text  string // Text to use as watermark
	Image string // Path to image file to use as watermark

	// Position and spacing
	Position Position // Position of watermark on image
	Margin   int      // Margin from edge in pixels

	// Opacity (0.0 = fully transparent, 1.0 = fully opaque)
	Opacity float64

	// Text-specific options
	TextSize  int    // Font size in pixels
	TextColor string // Hex color string (e.g., "#FFFFFF")

	// Image-specific options
	Scale float64 // Scale factor (0.1-1.0, as percentage of base image width)
}

// Validation errors
var (
	ErrNoWatermark       = errors.New("no watermark specified (neither text nor image)")
	ErrBothWatermarks    = errors.New("cannot specify both text and image watermarks")
	ErrInvalidPosition   = errors.New("invalid watermark position")
	ErrInvalidOpacity    = errors.New("opacity must be between 0.0 and 1.0")
	ErrInvalidScale      = errors.New("scale must be between 0.1 and 1.0")
	ErrInvalidTextSize   = errors.New("text size must be positive")
	ErrInvalidMargin     = errors.New("margin must be non-negative")
	ErrEmptyImage        = errors.New("base image is empty")
	ErrUnsupportedFormat = errors.New("unsupported image format")
)

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check that exactly one watermark type is specified
	if c.Text == "" && c.Image == "" {
		return ErrNoWatermark
	}
	if c.Text != "" && c.Image != "" {
		return ErrBothWatermarks
	}

	// Validate position
	validPositions := []Position{
		PositionTopLeft, PositionTopCenter, PositionTopRight,
		PositionLeftCenter, PositionCenter, PositionRightCenter,
		PositionBottomLeft, PositionBottomCenter, PositionBottomRight,
	}
	positionValid := false
	for _, validPos := range validPositions {
		if c.Position == validPos {
			positionValid = true
			break
		}
	}
	if !positionValid {
		return fmt.Errorf("%w: %s", ErrInvalidPosition, c.Position)
	}

	// Validate opacity
	if c.Opacity < 0.0 || c.Opacity > 1.0 {
		return fmt.Errorf("%w: %f", ErrInvalidOpacity, c.Opacity)
	}

	// Validate scale
	if c.Scale < 0.1 || c.Scale > 1.0 {
		return fmt.Errorf("%w: %f", ErrInvalidScale, c.Scale)
	}

	// Validate text size
	if c.TextSize <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidTextSize, c.TextSize)
	}

	// Validate margin
	if c.Margin < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidMargin, c.Margin)
	}

	return nil
}

// IsTextWatermark returns true if this is a text watermark configuration
func (c *Config) IsTextWatermark() bool {
	return c.Text != ""
}

// IsImageWatermark returns true if this is an image watermark configuration
func (c *Config) IsImageWatermark() bool {
	return c.Image != ""
}
