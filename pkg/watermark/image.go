package watermark

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

// LoadWatermarkImage loads an image from a file path
// Supports PNG and JPEG formats
func LoadWatermarkImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open watermark image: %w", err)
	}
	defer file.Close()

	// Try to decode the image (auto-detects format)
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode watermark image: %w", err)
	}

	// Verify format is supported
	if format != "png" && format != "jpeg" {
		return nil, fmt.Errorf("%w: %s (only PNG and JPEG are supported)", ErrUnsupportedFormat, format)
	}

	return img, nil
}

// ScaleImage scales an image to a specific width while maintaining aspect ratio
// The scale parameter is a factor (0.1-1.0) of the base image width
func ScaleImage(img image.Image, scale float64, baseWidth int) image.Image {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Calculate new dimensions
	newWidth := int(float64(baseWidth) * scale)
	if newWidth <= 0 {
		newWidth = 1 // Minimum 1 pixel
	}

	// Calculate new height maintaining aspect ratio
	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := int(float64(newWidth) * aspectRatio)
	if newHeight <= 0 {
		newHeight = 1
	}

	// Create a new image with the target dimensions
	scaled := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Scale using nearest neighbor (simple but fast)
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map scaled coordinates back to original coordinates
			srcX := (x * originalWidth) / newWidth
			srcY := (y * originalHeight) / newHeight

			// Get the color from the original image
			c := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY)
			scaled.Set(x, y, c)
		}
	}

	return scaled
}

// ApplyOpacity applies an opacity/alpha level to an image
// opacity should be between 0.0 (fully transparent) and 1.0 (fully opaque)
func ApplyOpacity(img image.Image, opacity float64) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// Apply opacity to each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Apply opacity to the alpha channel
			// RGBA() returns values in range [0, 65535], so we need to scale
			newAlpha := uint16(float64(a) * opacity)

			// Create new color with modified alpha
			newColor := color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: newAlpha,
			}

			result.Set(x, y, newColor)
		}
	}

	return result
}

// decodeImage decodes image bytes into an image.Image
func decodeImage(data []byte) (image.Image, string, error) {
	// Create a reader from bytes
	reader := &bytesReader{data: data}

	// Decode the image
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	return img, format, nil
}

// encodeImage encodes an image to bytes in the specified format
func encodeImage(img image.Image, format string) ([]byte, error) {
	buf := &bytesBuffer{}

	switch format {
	case "png":
		err := png.Encode(buf, img)
		if err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	case "jpeg":
		// Use high quality for JPEG to minimize quality loss
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
		if err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}

	return buf.data, nil
}

// bytesReader is a simple bytes reader that implements io.Reader
type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// bytesBuffer is a simple bytes buffer that implements io.Writer
type bytesBuffer struct {
	data []byte
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}
