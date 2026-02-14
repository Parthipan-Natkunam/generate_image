package watermark

import (
	"fmt"
	"image"
	"image/draw"
)

// Apply applies a watermark to a base image and returns the watermarked image bytes
//
// Parameters:
//   - baseImageData: the original image bytes (PNG or JPEG)
//   - cfg: watermark configuration
//
// Returns:
//   - watermarked image bytes in the same format as input
//   - error if watermarking fails
func Apply(baseImageData []byte, cfg Config) ([]byte, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid watermark configuration: %w", err)
	}

	// Validate base image is not empty
	if len(baseImageData) == 0 {
		return nil, ErrEmptyImage
	}

	// Decode the base image
	baseImg, format, err := decodeImage(baseImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base image: %w", err)
	}

	// Create or load the watermark image
	var watermarkImg image.Image
	if cfg.IsTextWatermark() {
		// Render text watermark
		watermarkImg, err = RenderText(cfg.Text, cfg.TextSize, cfg.TextColor)
		if err != nil {
			return nil, fmt.Errorf("failed to render text watermark: %w", err)
		}
	} else if cfg.IsImageWatermark() {
		// Load image watermark
		watermarkImg, err = LoadWatermarkImage(cfg.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to load image watermark: %w", err)
		}

		// Scale the watermark
		baseWidth := baseImg.Bounds().Dx()
		watermarkImg = ScaleImage(watermarkImg, cfg.Scale, baseWidth)
	}

	// Apply opacity to the watermark
	watermarkImg = ApplyOpacity(watermarkImg, cfg.Opacity)

	// Calculate watermark position
	baseWidth := baseImg.Bounds().Dx()
	baseHeight := baseImg.Bounds().Dy()
	wmWidth := watermarkImg.Bounds().Dx()
	wmHeight := watermarkImg.Bounds().Dy()

	x, y := CalculatePosition(baseWidth, baseHeight, wmWidth, wmHeight, cfg.Position, cfg.Margin)

	// Create a new image for the result (copy of base)
	resultImg := image.NewRGBA(baseImg.Bounds())
	draw.Draw(resultImg, resultImg.Bounds(), baseImg, image.Point{}, draw.Src)

	// Composite the watermark onto the base image
	wmRect := image.Rect(x, y, x+wmWidth, y+wmHeight)
	draw.Draw(resultImg, wmRect, watermarkImg, image.Point{}, draw.Over)

	// Encode the result back to the original format
	resultData, err := encodeImage(resultImg, format)
	if err != nil {
		return nil, fmt.Errorf("failed to encode watermarked image: %w", err)
	}

	return resultData, nil
}
