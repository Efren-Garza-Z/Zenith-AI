package main

import (
	"Zenith-AI/backends/zig_simd"
	"Zenith-AI/core"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// 1. Validar que el usuario pasó una ruta
	if len(os.Args) < 2 {
		fmt.Println("❌ Uso: ./ZenithEngine.exe <ruta_archivo_o_carpeta>")
		return
	}

	targetPath := os.Args[1]
	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("❌ Error: No se pudo encontrar la ruta: %v\n", err)
		return
	}

	// 2. Crear carpeta de salida si no existe
	outputDir := "zenith_results"
	os.MkdirAll(outputDir, os.ModePerm)

	// 3. Decidir si procesar uno o muchos
	if info.IsDir() {
		fmt.Printf("📂 Procesando carpeta: %s\n", targetPath)
		processDirectory(targetPath, outputDir)
	} else {
		processSingleFile(targetPath, outputDir)
	}
}

func processSingleFile(path, outDir string) {
	// Filtrar extensiones válidas
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return
	}

	fmt.Printf("📸 Procesando: %s... ", filepath.Base(path))

	input, err := core.LoadImageToTensor(path)
	if err != nil {
		fmt.Printf("Error al cargar: %v\n", err)
		return
	}

	// --- PIPELINE ---
	clean := zig_simd.GaussianBlur(input)
	edges := zig_simd.ConvolveWithZig(clean, zig_simd.SobelX)
	final := zig_simd.Threshold(edges, 0.2)
	// ----------------

	// Análisis
	pct := zig_simd.GetActivePixelPercentage(final)

	// Guardar con nombre prefijado en la carpeta de salida
	newName := "proc_" + filepath.Base(path)
	outputPath := filepath.Join(outDir, newName)
	core.SaveTensorToImage(outputPath, final)

	fmt.Printf("Done! (Contenido: %.2f%%) -> %s\n", pct, outputPath)
}

func processDirectory(dirPath, outDir string) {
	files, _ := os.ReadDir(dirPath)
	start := time.Now()
	count := 0

	for _, f := range files {
		if !f.IsDir() {
			fullPath := filepath.Join(dirPath, f.Name())
			processSingleFile(fullPath, outDir)
			count++
		}
	}

	fmt.Printf("\n✨ Finalizado. %d imágenes procesadas en %s\n", count, time.Since(start))
}
