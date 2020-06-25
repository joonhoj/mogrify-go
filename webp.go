package mogrify

import (
	"bytes"
	"fmt"
	"io"
)

// Webp image that can be transformed.
type Webp struct {
	// Embed GdImage and all it's methods
	GdImage
}

// DecodeWebp decodes a WEBP image from a reader.
func DecodeWebp(reader io.Reader) (Image, error) {
	var image Webp

	image.gd = gdCreateFromWebp(drain(reader))
	if image.gd == nil {
		return nil, fmt.Errorf("couldn't create Webp decoder")
	}

	return &image, nil
}

// EncodeWebp encodes the image onto the writer as a WEBP.
func EncodeWebp(w io.Writer, img Image) (int64, error) {
	slice, err := img.image().gdImageWebp()
	if err != nil {
		return 0, err
	}

	return bytes.NewBuffer(slice).WriteTo(w)
}
