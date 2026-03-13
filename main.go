//package main
//
//import (
//	"github.com/efren-garza-z/zenith-ai/backends/zig_simd"
//	"github.com/efren-garza-z/zenith-ai/core"
//	"fmt"
//	"os"
//	"path/filepath"
//	"strconv" // Para convertir texto a número
//	"strings"
//	"time"
//)
//
//func main() {
//	// 1. Validar argumentos (Ruta obligatoria, Umbral opcional)
//	if len(os.Args) < 2 {
//		fmt.Println("❌ Uso: ./ZenithEngine.exe <archivo_o_carpeta> [umbral: 0.0-1.0]")
//		fmt.Println("Ejemplo: ./ZenithEngine.exe firma.jpg 0.25")
//		return
//	}
//
//	targetPath := os.Args[1]
//
//	// 2. Definir umbral (por defecto 0.2 si el usuario no pone nada)
//	thresholdValue := 0.2
//	if len(os.Args) >= 3 {
//		parsed, err := strconv.ParseFloat(os.Args[2], 64)
//		if err == nil && parsed >= 0 && parsed <= 1 {
//			thresholdValue = parsed
//		} else {
//			fmt.Println("⚠️  Umbral inválido. Usando valor por defecto: 0.2")
//		}
//	}
//
//	info, err := os.Stat(targetPath)
//	if err != nil {
//		fmt.Printf("❌ Error: Ruta no encontrada: %v\n", err)
//		return
//	}
//
//	outputDir := "zenith_results"
//	os.MkdirAll(outputDir, os.ModePerm)
//
//	fmt.Printf("🚀 Zenith-AI | Umbral configurado: %.2f\n", thresholdValue)
//
//	if info.IsDir() {
//		processDirectory(targetPath, outputDir, thresholdValue)
//	} else {
//		processSingleFile(targetPath, outputDir, thresholdValue)
//	}
//}
//
//func processSingleFile(path, outDir string, threshold float64) {
//	ext := strings.ToLower(filepath.Ext(path))
//	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
//		return
//	}
//
//	input, err := core.LoadImageToTensor(path)
//	if err != nil {
//		fmt.Printf("❌ Error al cargar %s: %v\n", path, err)
//		return
//	}
//
//	// PIPELINE
//	clean := zig_simd.GaussianBlur(input)
//	edges := zig_simd.ConvolveWithZig(clean, zig_simd.SobelX)
//	final := zig_simd.Threshold(edges, threshold) // Usamos el valor dinámico
//
//	pct := zig_simd.GetActivePixelPercentage(final)
//	com := zig_simd.GetCenterOfMass(final)
//
//	newName := "proc_" + filepath.Base(path)
//	outputPath := filepath.Join(outDir, newName)
//	core.SaveTensorToImage(outputPath, final)
//
//	// Imprimimos el reporte completo
//	fmt.Printf("📸 %-15s | Tinta: %5.2f%% | CoM: (X:%.1f, Y:%.1f) | 📁 %s\n",
//		filepath.Base(path), pct, com.X, com.Y, newName)
//}
//
//func processDirectory(dirPath, outDir string, threshold float64) {
//	files, _ := os.ReadDir(dirPath)
//	start := time.Now()
//	count := 0
//
//	for _, f := range files {
//		if !f.IsDir() {
//			processSingleFile(filepath.Join(dirPath, f.Name()), outDir, threshold)
//			count++
//		}
//	}
//	fmt.Printf("\n✨ Finalizado. %d archivos procesados en %s\n", count, time.Since(start))
//}

package main

import (
	"fmt"
	"os"

	"github.com/efren-garza-z/zenith-ai/backends/zig_simd"
	"github.com/efren-garza-z/zenith-ai/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: ./ZenithEngine.exe <imagen_de_firma>")
		return
	}

	imgPath := os.Args[1]
	fmt.Printf("🚀 Analizando firma: %s\n", imgPath)

	// 1. Cargar imagen
	input, err := core.LoadImageToTensor(imgPath)
	if err != nil {
		fmt.Printf("❌ Error al cargar: %v\n", err)
		return
	}

	// 2. Pipeline de Procesamiento
	clean := zig_simd.GaussianBlur(input)
	edges := zig_simd.ConvolveWithZig(clean, zig_simd.SobelX)
	binarized := zig_simd.Threshold(edges, 0.2)

	// 3. ANÁLISIS ANTES DEL CENTRADO
	comOriginal := zig_simd.GetCenterOfMass(binarized)
	fmt.Printf("📍 CoM Original: (%.1f, %.1f)\n", comOriginal.X, comOriginal.Y)

	// 4. ALINEACIÓN (La nueva magia)
	centered := zig_simd.CenterSignature(binarized)

	// 5. ANÁLISIS DESPUÉS DEL CENTRADO
	comFinal := zig_simd.GetCenterOfMass(centered)
	fmt.Printf("🎯 CoM Centrado: (%.1f, %.1f) [Objetivo: %.1f, %.1f]\n",
		comFinal.X, comFinal.Y, float64(centered.Width)/2, float64(centered.Height)/2)

	// 6. Guardar resultados para comparar
	os.MkdirAll("debug_output", os.ModePerm)
	core.SaveTensorToImage("debug_output/1_binarizada.png", binarized)
	core.SaveTensorToImage("debug_output/2_centrada.png", centered)

	fmt.Println("✅ Proceso terminado. Revisa la carpeta 'debug_output'")
}
