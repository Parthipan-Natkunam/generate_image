package watermark

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// RenderText renders text as an image with the specified size and color
func RenderText(text string, size int, colorHex string) (image.Image, error) {
	// Parse the color
	textColor, err := ParseHexColor(colorHex)
	if err != nil {
		return nil, fmt.Errorf("invalid text color: %w", err)
	}

	// Use basicfont for MVP (simple built-in font)
	// Note: basicfont is fixed-size, so we'll scale the result if needed
	face := basicfont.Face7x13

	// Calculate text dimensions
	textWidth := font.MeasureString(face, text).Ceil()
	textHeight := face.Metrics().Height.Ceil()

	// Create an image for the text
	textImg := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))

	// Fill with transparent background
	draw.Draw(textImg, textImg.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Create a drawer for rendering text
	drawer := &font.Drawer{
		Dst:  textImg,
		Src:  &image.Uniform{textColor},
		Face: face,
		Dot:  fixed.Point26_6{X: 0, Y: face.Metrics().Ascent},
	}

	// Draw the text
	drawer.DrawString(text)

	// Scale the text to the desired size if needed
	// basicfont is roughly 13 pixels tall, so we scale proportionally
	baseFontHeight := 13.0
	scaleFactor := float64(size) / baseFontHeight

	if scaleFactor != 1.0 {
		scaledWidth := int(float64(textWidth) * scaleFactor)
		scaledHeight := int(float64(textHeight) * scaleFactor)

		if scaledWidth <= 0 {
			scaledWidth = 1
		}
		if scaledHeight <= 0 {
			scaledHeight = 1
		}

		// Create scaled image
		scaled := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))

		// Scale using nearest neighbor
		for y := 0; y < scaledHeight; y++ {
			for x := 0; x < scaledWidth; x++ {
				srcX := int(float64(x) / scaleFactor)
				srcY := int(float64(y) / scaleFactor)
				if srcX >= textWidth {
					srcX = textWidth - 1
				}
				if srcY >= textHeight {
					srcY = textHeight - 1
				}
				c := textImg.At(srcX, srcY)
				scaled.Set(x, y, c)
			}
		}

		return scaled, nil
	}

	return textImg, nil
}

// ParseHexColor parses a hex color string (e.g., "#FFFFFF" or "FFFFFF")
// and returns a color.Color
func ParseHexColor(hex string) (color.Color, error) {
	// Remove the "#" prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Validate length
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex color format: %s (expected format: #RRGGBB)", hex)
	}

	// Parse red component
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid red component in hex color: %w", err)
	}

	// Parse green component
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid green component in hex color: %w", err)
	}

	// Parse blue component
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid blue component in hex color: %w", err)
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255, // Fully opaque (opacity will be applied separately)
	}, nil
}
