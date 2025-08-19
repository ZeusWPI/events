// Package image provides various image related functions
package image

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Add support for decoding jpg's
	_ "image/png"  // Add support for decoding png's

	_ "golang.org/x/image/webp" // Add support for decoding webp's

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func ToWebp(imageBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("converting bytes to image.Image: %w", err)
	}

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetPicture, 60)
	if err != nil {
		return nil, fmt.Errorf("creating lossless webp encoder options: %w", err)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, fmt.Errorf("encode bytes to webp: %w", err)
	}

	return buf.Bytes(), nil
}
