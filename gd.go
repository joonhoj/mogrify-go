package mogrify

// #cgo LDFLAGS: -lgd -lexif
// #include <gd.h>
/*
#include <stdlib.h>
#include <string.h>

typedef struct _YCbCr {
  unsigned char* ch[3];
} YCbCr;

void free_ycbcr(YCbCr* ycbcr);

int RGBToYCbCr(unsigned char r, unsigned char g, unsigned char b) {
	int r1 = (int)r;
	int g1 = (int)g;
	int b1 = (int)b;

	int yy = (19595*r1 + 38470*g1 + 7471*b1 + (1<<15)) >> 16;

	int cb = -11056*r1 - 21712*g1 + 32768*b1 + (257<<15);
	if ((cb & 0xff000000) == 0) {
		cb >>= 16;
	} else {
		cb = ~(cb >> 31);
	}

	int cr = 32768*r1 - 27440*g1 - 5328*b1 + (257<<15);
	if ((cr & 0xff000000) == 0) {
		cr >>= 16;
	} else {
		cr = ~(cr >> 31);
	}

	return ((yy & 0xff) << 16) | ((cb & 0xff) << 8) | (cr & 0xff);
}

YCbCr* get_ycbcr(gdImagePtr ptr, int* len) {
	YCbCr* ycbcr = NULL;
	int i = 0;

	ycbcr = (YCbCr*)calloc(1, sizeof(YCbCr*));
	if (ycbcr == NULL) {
		return NULL;
	}

	for (int i = 0; i < 3; i++) {
		ycbcr->ch[i] = (unsigned char*)malloc(ptr->sx * ptr->sy);
		if (ycbcr->ch[i] == NULL) {
			free_ycbcr(ycbcr);
			return NULL;
		}
	}

	for (int y = 0; y < ptr->sy; y++) {
		for (int x = 0; x < ptr->sx; x++) {
			int c = gdImageGetTrueColorPixel(ptr, x, y);
			int res = RGBToYCbCr(gdTrueColorGetRed(c), gdTrueColorGetGreen(c), gdTrueColorGetBlue(c));
			ycbcr->ch[0][i] = (unsigned char)(res >> 16);
			ycbcr->ch[1][i] = (unsigned char)(res >> 8);
			ycbcr->ch[2][i] = (unsigned char)res;
			i++;
		}
	}

	*len = ptr->sx * ptr->sy;
	return ycbcr;
}

unsigned char* get_pixels(gdImagePtr ptr, int* len) {
	int pitch = ptr->trueColor ? sizeof(int) * ptr->sx : ptr->sx;
	void* src = ptr->trueColor ? (void*) ptr->tpixels : (void*) ptr->pixels;
	void* dest = NULL;
	void* buf = NULL;

	buf = malloc(pitch * ptr->sy);
	if (buf == NULL) {
		return NULL;
	}

	dest = buf;
	for (int i = 0; i < ptr->sy; i++) {
		memcpy(dest, src, pitch);
		dest += pitch;
	}

	*len = pitch * ptr->sy;
	return (unsigned char*) buf;
}

unsigned char* get_rgb_pixels(gdImagePtr ptr, int* len) {
	int pitch = 3 * ptr->sx; // RGB
	unsigned char* buf = NULL;
	int i = 0;

	buf = (unsigned char*) malloc(pitch * ptr->sy);
	if (buf == NULL) {
		return NULL;
	}

	for (int y = 0; y < ptr->sy; y++) {
		for (int x = 0; x < ptr->sx; x++) {
			int c = gdImageGetTrueColorPixel(ptr, x, y);
			buf[i++] = gdTrueColorGetRed(c);
			buf[i++] = gdTrueColorGetGreen(c);
			buf[i++] = gdTrueColorGetBlue(c);
		}
	}

	*len = pitch * ptr->sy;
	return buf;
}

unsigned char* get_quantization_pixels(gdImagePtr ptr, int* len) {
	int pitch = 3 * ptr->sx; // RGB
	unsigned char* buf = NULL;
	int i = 0;

	buf = (unsigned char*) malloc(pitch * ptr->sy);
	if (buf == NULL) {
		return NULL;
	}

	for (int y = 0; y < ptr->sy; y++) {
		for (int x = 0; x < ptr->sx; x++) {
			int c = gdImageGetPixel(ptr, x, y);
			buf[i++] = gdTrueColorGetRed(c) >> 2;
			buf[i++] = gdTrueColorGetGreen(c) >> 2;
			buf[i++] = gdTrueColorGetBlue(c) >> 2;
		}
	}

	*len = pitch * ptr->sy;
	return buf;
}

void free_pixels(unsigned char* pixels) {
	if (pixels != NULL) free(pixels);
}

void free_ycbcr(YCbCr* ycbcr) {
	if (ycbcr != NULL) {
		for (int i = 0; i < 3; i++) {
			if (ycbcr->ch[i] != NULL) free(ycbcr->ch[i]);
		}
		free(ycbcr);
	}
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	imageError  = errors.New("[GD] image is nil")
	createError = errors.New("[GD] cannot create new image")
	writeError  = errors.New("[GD] image cannot be written")
)

type gdImage struct {
	img *C.gdImage
}

func img(img *C.gdImage) *gdImage {
	if img == nil {
		return nil
	}

	image := &gdImage{img}
	if isInvalid(image) {
		return nil
	}
	return image
}

func cbool(b bool) int {
	if b == true {
		return 1
	}
	return 0
}

func gdCreate(sx, sy int) *gdImage {
	img := img(C.gdImageCreateTrueColor(C.int(sx), C.int(sy)))

	if img == nil {
		return nil
	}

	C.gdImageAlphaBlending(img.img, C.int(cbool(false)))
	C.gdImageSaveAlpha(img.img, C.int(cbool(true)))

	return img
}

func (p *gdImage) gdAlphaBlending(b bool) {
	C.gdImageAlphaBlending(p.img, C.int(cbool(b)))
}

func (p *gdImage) gdSaveAlpha(b bool) {
	C.gdImageSaveAlpha(p.img, C.int(cbool(b)))
}

func gdCreateFromJpeg(buffer []byte) *gdImage {
	return img(C.gdImageCreateFromJpegPtr(C.int(len(buffer)), unsafe.Pointer(&buffer[0])))
}

func gdCreateFromGif(buffer []byte) *gdImage {
	return img(C.gdImageCreateFromGifPtr(C.int(len(buffer)), unsafe.Pointer(&buffer[0])))
}

func gdCreateFromPng(buffer []byte) *gdImage {
	return img(C.gdImageCreateFromPngPtr(C.int(len(buffer)), unsafe.Pointer(&buffer[0])))
}

func gdCreateFromWebp(buffer []byte) *gdImage {
	return img(C.gdImageCreateFromWebpPtr(C.int(len(buffer)), unsafe.Pointer(&buffer[0])))
}

func (p *gdImage) gdDestroy() {
	if p != nil && p.img != nil {
		C.gdImageDestroy(p.img)
	}
}

func isInvalid(p *gdImage) bool {
	return p == nil || p.img == nil
}

func (p *gdImage) width() int {
	if p == nil {
		panic(imageError)
	}
	return int((*p.img).sx)
}

func (p *gdImage) transparent() bool {
	if p == nil {
		panic(imageError)
	}
	return int((*p.img).transparent) == 1
}

func (p *gdImage) height() int {
	if p == nil {
		panic(imageError)
	}
	return int((*p.img).sy)
}

func (p *gdImage) gdCopy(dstX, dstY, srcX, srcY, dstW, dstH int) *gdImage {

	if p == nil || p.img == nil {
		panic(imageError)
	}

	dst := gdCreate(dstW, dstH)

	if dst == nil {
		return nil
	}

	C.gdImageCopy(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(dstW), C.int(dstH))

	if isInvalid(dst) {
		dst.gdDestroy()
		return nil
	}

	return dst
}

func (p *gdImage) gdCopyResampled(dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) *gdImage {

	if p == nil || p.img == nil {
		panic(imageError)
	}

	dst := gdCreate(dstW, dstH)

	if dst == nil {
		return nil
	}

	C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))

	if isInvalid(dst) {
		dst.gdDestroy()
		return nil
	}

	return dst
}

func (p *gdImage) gdCopyResized(dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) *gdImage {
	if p == nil || p.img == nil {
		panic(imageError)
	}

	dst := gdCreate(dstW, dstH)

	if dst == nil {
		return nil
	}

	transparency := C.gdImageColorAllocateAlpha(dst.img, 255, 255, 255, 127)
	C.gdImageFilledRectangle(dst.img, C.int(dstX), C.int(dstY), C.int(dstW), C.int(dstH), transparency)
	C.gdImageColorTransparent(dst.img, transparency)

	C.gdImageCopyResized(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))

	if isInvalid(dst) {
		dst.gdDestroy()
		return nil
	}

	return dst
}

func (p *gdImage) gdImagePng() ([]byte, error) {
	if p == nil {
		panic(imageError)
	}

	var size C.int

	data := C.gdImagePngPtr(p.img, &size)
	if data == nil || int(size) == 0 {
		return []byte{}, writeError
	}

	defer C.gdFree(unsafe.Pointer(data))

	return C.GoBytes(data, size), nil
}

func (p *gdImage) gdImageGif() ([]byte, error) {
	if p == nil {
		panic(imageError)
	}

	var size C.int

	data := C.gdImageGifPtr(p.img, &size)
	if data == nil || int(size) == 0 {
		return []byte{}, writeError
	}

	defer C.gdFree(unsafe.Pointer(data))

	return C.GoBytes(data, size), nil
}

func (p *gdImage) gdImageJpeg(quality int) ([]byte, error) {
	if p == nil {
		panic(imageError)
	}

	var size C.int

	// use -1 as quality, this will mean to use standard Jpeg quality
	data := C.gdImageJpegPtr(p.img, &size, C.int(quality))
	if data == nil || int(size) == 0 {
		return []byte{}, writeError
	}

	defer C.gdFree(unsafe.Pointer(data))

	return C.GoBytes(data, size), nil
}

func (p *gdImage) gdImageWebp() ([]byte, error) {
	if p == nil {
		panic(imageError)
	}

	var size C.int

	data := C.gdImageWebpPtr(p.img, &size)
	if data == nil || int(size) == 0 {
		return []byte{}, writeError
	}

	defer C.gdFree(unsafe.Pointer(data))

	return C.GoBytes(data, size), nil
}

func (p *gdImage) gdImagePixels() ([]byte, error) {
	var len C.int

	pixels := C.get_pixels(p.img, &len)
	if pixels == nil {
		return nil, errors.New("failed to get pixels")
	}
	defer C.free_pixels(pixels)

	bytes := C.GoBytes(unsafe.Pointer(pixels), len)

	return bytes, nil
}

func (p *gdImage) gdImageRGBPixels() ([]byte, error) {
	var len C.int

	pixels := C.get_rgb_pixels(p.img, &len)
	if pixels == nil {
		return nil, errors.New("failed to get pixels")
	}
	defer C.free_pixels(pixels)

	bytes := C.GoBytes(unsafe.Pointer(pixels), len)

	return bytes, nil
}

func (p *gdImage) gdImageQuantizationPixels() ([]byte, error) {
	var len C.int

	pixels := C.get_quantization_pixels(p.img, &len)
	if pixels == nil {
		return nil, errors.New("failed to get pixels")
	}
	defer C.free_pixels(pixels)

	bytes := C.GoBytes(unsafe.Pointer(pixels), len)

	return bytes, nil
}

func (p *gdImage) gdImageYCbCr() ([][]byte, error) {
	var len C.int

	ycbcr := C.get_ycbcr(p.img, &len)
	if ycbcr == nil {
		return nil, errors.New("failed to get ycbcr")
	}
	defer C.free_ycbcr(ycbcr)

	ch := [][]byte{
		C.GoBytes(unsafe.Pointer(ycbcr.ch[0]), len),
		C.GoBytes(unsafe.Pointer(ycbcr.ch[1]), len),
		C.GoBytes(unsafe.Pointer(ycbcr.ch[2]), len),
	}

	return ch, nil
}
