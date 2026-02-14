package watermark

// CalculatePosition calculates the x,y coordinates for placing a watermark
// on a base image based on the desired position and margin.
//
// Parameters:
//   - baseWidth, baseHeight: dimensions of the base image
//   - wmWidth, wmHeight: dimensions of the watermark
//   - pos: desired position (e.g., PositionBottomRight)
//   - margin: distance in pixels from the edge
//
// Returns:
//   - x, y: top-left coordinates where the watermark should be placed
func CalculatePosition(baseWidth, baseHeight, wmWidth, wmHeight int, pos Position, margin int) (x, y int) {
	switch pos {
	case PositionTopLeft:
		return margin, margin

	case PositionTopCenter:
		return (baseWidth-wmWidth)/2, margin

	case PositionTopRight:
		return baseWidth - wmWidth - margin, margin

	case PositionLeftCenter:
		return margin, (baseHeight-wmHeight)/2

	case PositionCenter:
		return (baseWidth-wmWidth)/2, (baseHeight-wmHeight)/2

	case PositionRightCenter:
		return baseWidth - wmWidth - margin, (baseHeight-wmHeight)/2

	case PositionBottomLeft:
		return margin, baseHeight - wmHeight - margin

	case PositionBottomCenter:
		return (baseWidth-wmWidth)/2, baseHeight - wmHeight - margin

	case PositionBottomRight:
		return baseWidth - wmWidth - margin, baseHeight - wmHeight - margin

	default:
		// Default to bottom-right if position is somehow invalid
		return baseWidth - wmWidth - margin, baseHeight - wmHeight - margin
	}
}
