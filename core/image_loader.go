package core

import (
	"image"
	"image/draw"
	_ "image/gif"  // Registro de decodificador GIF
	_ "image/jpeg" // Registro de decodificador JPEG
	_ "image/png"  // Registro de decodificador PNG
	"os"

	_ "golang.org/x/image/webp" // Registro de decodificador WebP
)

func LoadImageToTensor(path string) (*Tensor4D, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	tensor := NewTensor4D(1, 3, height, width)

	// Optimizamos: Convertimos la imagen a RGBA si no lo es para acceder a sus bytes directamente
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	}

	// Acceso directo a los píxeles (esto es miles de veces más rápido que At(x,y))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := rgba.PixOffset(x, y)
			tensor.Set(0, 0, y, x, float64(rgba.Pix[idx])/255.0)   // R
			tensor.Set(0, 1, y, x, float64(rgba.Pix[idx+1])/255.0) // G
			tensor.Set(0, 2, y, x, float64(rgba.Pix[idx+2])/255.0) // B
		}
	}
	return tensor, nil
}
