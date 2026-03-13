# 🚀 Zenith-AI: High-Performance Image Engine

Zenith-AI es un motor de pre-procesamiento de imágenes de alto rendimiento diseñado para flujos de trabajo de Inteligencia Artificial y Visión Artificial. Utiliza una arquitectura híbrida que delega el cálculo matemático pesado a Zig (SIMD) mientras mantiene una lógica de control ágil en Go.

---

## ✨ Características

- **Velocidad Extrema:** Filtros de convolución optimizados con instrucciones SIMD en Zig.
- **Multipropósito:** Soporta JPEG, PNG, GIF y WebP (Google Next-Gen).
- **Pipeline de Análisis:** Pipeline integrado: Gaussian Blur → Sobel Filter → Thresholding.
- **Métricas en Tiempo Real:** Cálculo automático del porcentaje de contenido activo (tinta/bordes).
- **CLI Versátil:** Procesa imágenes individuales o carpetas completas de datasets.

---

## 🛠️ Requisitos del Sistema

Esta librería utiliza **CGO** y **Zig** para alcanzar un rendimiento superior mediante instrucciones SIMD. Para usarla, el sistema debe tener instalado:

1. **Zig Compiler (v0.11.0 o superior):** [Descargar aquí](https://ziglang.org/download/).
2. **Go (v1.20 o superior).**
3. **Librerías de Go:**
```bash
go get golang.org/x/image/webp
```

---

## 📦 Estructura del Proyecto

- `/core` — Manejo de Tensores 4D y carga de imágenes optimizada.
- `/backends/zig_simd` — Código fuente en Zig y puente CGO.
- `/zenith_results` — Carpeta generada automáticamente con los resultados.

---

## 🏗️ Instalación y Configuración Nativa

Debido a que Go no compila código Zig automáticamente, debes generar el binario nativo antes de compilar tu proyecto:

### 1. Descargar la librería
```bash
go get github.com/efren-garza-z/Zenith-AI
```

### 2. Compilar el Motor Nativo (Zig)

Navega a la carpeta donde se instaló la librería (o dentro de tu `vendor`) y ejecuta:
```bash
cd backends/zig_simd
zig build-lib conv.zig -O ReleaseFast -mcpu=native

# Windows:
mv conv.lib ../../libconv.a

# Linux/macOS:
mv libconv.a ../../
```

### 3. Compilar tu Proyecto

Para compilar cualquier aplicación que use esta librería, debes indicar a Go que use Zig como compilador de C:
```powershell
# Windows (PowerShell)
$env:CGO_ENABLED="1"
$env:CC="zig cc -target x86_64-windows-gnu"
go build .
```
```bash
# Linux/macOS
CGO_ENABLED=1 CC="zig cc" go build .
```

---

## 🚀 Uso y Ejecución

El motor detecta automáticamente si la ruta es un archivo o una carpeta.

**Procesar una sola imagen:**
```powershell
./ZenithEngine.exe mi_firma.png
```

**Procesar un dataset completo:**
```powershell
./ZenithEngine.exe ./dataset_entrenamiento
```

> **Nota:** Todos los resultados se guardarán en la carpeta `zenith_results/` con el prefijo `proc_`.

---

## 🧠 Pipeline Técnico

1. **Load** — Decodificación rápida mediante acceso directo a memoria (RGBA Pixels).
2. **Clean** — Gaussian Blur 5×5 para eliminación de ruido.
3. **Extract** — Sobel Operator para detección de bordes y rasgos.
4. **Binarize** — Thresholding para convertir la imagen en datos binarios (0 y 1).
5. **Analyze** — Cálculo de densidad de píxeles activos para filtrado de datos.

---

## ⚖️ Licencia

Este proyecto es una herramienta de código abierto para la comunidad de desarrolladores de IA.