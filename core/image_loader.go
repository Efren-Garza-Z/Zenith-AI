package core

import (
	"image"
	_ "image/gif"  // Registro de decodificador GIF
	_ "image/jpeg" // Registro de decodificador JPEG
	_ "image/png"  // Registro de decodificador PNG
	"os"

	"golang.org/x/image/draw"
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

// ResizeTensor redimensiona el tensor a un tamaño fijo para normalizar la comparación
func ResizeImage(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	// Utilizamos BiLinear para no perder calidad en los trazos de la firma
	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return newImg
}

func LoadAndNormalizeImage(path string, w, h int) (*Tensor4D, error) {
	file, _ := os.Open(path)
	defer file.Close()
	img, _, _ := image.Decode(file)

	// Forzamos el tamaño estándar
	resized := ResizeImage(img, w, h)

	bounds := resized.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	tensor := NewTensor4D(1, 3, height, width)

	rgba := resized.(*image.RGBA)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := rgba.PixOffset(x, y)
			tensor.Set(0, 0, y, x, float64(rgba.Pix[idx])/255.0)
			tensor.Set(0, 1, y, x, float64(rgba.Pix[idx+1])/255.0)
			tensor.Set(0, 2, y, x, float64(rgba.Pix[idx+2])/255.0)
		}
	}
	return tensor, nil
}
