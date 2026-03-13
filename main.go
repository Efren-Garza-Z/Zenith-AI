package main

import (
	"Zenith-AI/backends/zig_simd"
	"Zenith-AI/core"
	"fmt"
	"os"
	"path/filepath"
	"strconv" // Para convertir texto a número
	"strings"
	"time"
)

func main() {
	// 1. Validar argumentos (Ruta obligatoria, Umbral opcional)
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: ./ZenithEngine.exe <archivo_o_carpeta> [umbral: 0.0-1.0]")
		fmt.Println("Ejemplo: ./ZenithEngine.exe firma.jpg 0.25")
		return
	}

	targetPath := os.Args[1]

	// 2. Definir umbral (por defecto 0.2 si el usuario no pone nada)
	thresholdValue := 0.2
	if len(os.Args) >= 3 {
		parsed, err := strconv.ParseFloat(os.Args[2], 64)
		if err == nil && parsed >= 0 && parsed <= 1 {
			thresholdValue = parsed
		} else {
			fmt.Println("⚠️  Umbral inválido. Usando valor por defecto: 0.2")
		}
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("❌ Error: Ruta no encontrada: %v\n", err)
		return
	}

	outputDir := "zenith_results"
	os.MkdirAll(outputDir, os.ModePerm)

	fmt.Printf("🚀 Zenith-AI | Umbral configurado: %.2f\n", thresholdValue)

	if info.IsDir() {
		processDirectory(targetPath, outputDir, thresholdValue)
	} else {
		processSingleFile(targetPath, outputDir, thresholdValue)
	}
}

func processSingleFile(path, outDir string, threshold float64) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return
	}

	input, err := core.LoadImageToTensor(path)
	if err != nil {
		fmt.Printf("❌ Error al cargar %s: %v\n", path, err)
		return
	}

	// PIPELINE
	clean := zig_simd.GaussianBlur(input)
	edges := zig_simd.ConvolveWithZig(clean, zig_simd.SobelX)
	final := zig_simd.Threshold(edges, threshold) // Usamos el valor dinámico

	pct := zig_simd.GetActivePixelPercentage(final)

	newName := "proc_" + filepath.Base(path)
	outputPath := filepath.Join(outDir, newName)
	core.SaveTensorToImage(outputPath, final)

	fmt.Printf("📸 %-20s | Contenido: %6.2f%% | Salida: %s\n", filepath.Base(path), pct, newName)
}

func processDirectory(dirPath, outDir string, threshold float64) {
	files, _ := os.ReadDir(dirPath)
	start := time.Now()
	count := 0

	for _, f := range files {
		if !f.IsDir() {
			processSingleFile(filepath.Join(dirPath, f.Name()), outDir, threshold)
			count++
		}
	}
	fmt.Printf("\n✨ Finalizado. %d archivos procesados en %s\n", count, time.Since(start))
}
