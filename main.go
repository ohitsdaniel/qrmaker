package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/skip2/go-qrcode"
	"github.com/vitrun/qart"
)

func main() {
	// Command line flags
	imagePath := flag.String("image", "", "Path to the image to embed/overlay in the QR code")
	url := flag.String("url", "https://example.com", "URL or text to encode in the QR code")
	output := flag.String("output", "qrcode.png", "Output file path for the generated QR code")
	version := flag.Int("version", 6, "QR code version (1-8 for embed, auto for overlay)")
	scale := flag.Int("scale", 8, "Size of each QR module in pixels")
	mask := flag.Int("mask", 2, "QR mask pattern (0-7, embed mode only)")
	dx := flag.Int("dx", 4, "X offset for image positioning (embed mode only)")
	dy := flag.Int("dy", 4, "Y offset for image positioning (embed mode only)")
	randControl := flag.Bool("rand", false, "Use random control for pixel placement (embed mode only)")
	dither := flag.Bool("dither", false, "Enable dithering for better image representation (embed mode only)")
	onlyData := flag.Bool("only-data", false, "Only use data pixels for the image (embed mode only)")
	saveControl := flag.Bool("save-control", false, "Save control image instead of QR code (embed mode only)")
	seed := flag.Int64("seed", 0, "Random seed for reproducible results (0 = use current time)")

	// Overlay mode flags
	overlay := flag.Bool("overlay", false, "Overlay image on QR code instead of embedding")
	overlaySize := flag.Int("overlay-size", 25, "Size of overlay image as percentage of QR code (10-40)")
	size := flag.Int("size", 512, "Output image size in pixels (overlay mode)")

	flag.Parse()

	if *imagePath == "" {
		printHelp()
		os.Exit(0)
	}

	var qrData []byte
	var err error

	if *overlay {
		qrData, err = generateOverlayQR(*url, *imagePath, *size, *overlaySize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating overlay QR: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("QR code generated: %s\n", *output)
		fmt.Printf("  URL: %s\n", *url)
		fmt.Printf("  Mode: overlay, Size: %dx%d, Logo: %d%%\n", *size, *size, *overlaySize)
	} else {
		// Read the source image for embed mode
		imgData, err := os.ReadFile(*imagePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading image: %v\n", err)
			os.Exit(1)
		}

		// Set seed
		seedVal := *seed
		if seedVal == 0 {
			seedVal = time.Now().UnixNano()
		}

		// Generate the QR code with embedded image
		qrData = qart.Encode(
			*url,
			imgData,
			seedVal,
			*version,
			*scale,
			*mask,
			*dx,
			*dy,
			*randControl,
			*dither,
			*onlyData,
			*saveControl,
		)

		if qrData == nil {
			fmt.Fprintf(os.Stderr, "Error: failed to generate QR code\n")
			fmt.Fprintf(os.Stderr, "Tip: Try a higher version (-version 7 or 8) for longer URLs\n")
			os.Exit(1)
		}

		fmt.Printf("QR code generated: %s\n", *output)
		fmt.Printf("  URL: %s\n", *url)
		fmt.Printf("  Mode: embed, Version: %d, Scale: %d\n", *version, *scale)
		if *dither {
			fmt.Println("  Dithering: enabled")
		}
		if *randControl {
			fmt.Println("  Random control: enabled")
		}
	}

	// Save the output
	if err := os.WriteFile(*output, qrData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}
}

func generateOverlayQR(url, imagePath string, qrSize, overlaySizePct int) ([]byte, error) {
	// Clamp overlay size
	if overlaySizePct < 10 {
		overlaySizePct = 10
	}
	if overlaySizePct > 40 {
		overlaySizePct = 40
	}

	// Minimum QR size
	if qrSize < 128 {
		qrSize = 128
	}

	qrc, err := qrcode.New(url, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("creating QR code: %w", err)
	}

	qrImg := qrc.Image(qrSize)

	// Load overlay image
	overlayFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("opening overlay image: %w", err)
	}
	defer overlayFile.Close()

	overlayImg, _, err := image.Decode(overlayFile)
	if err != nil {
		return nil, fmt.Errorf("decoding overlay image: %w", err)
	}

	// Calculate overlay size and position
	overlayWidth := qrSize * overlaySizePct / 100
	overlayHeight := qrSize * overlaySizePct / 100

	// Resize overlay image to fit
	resizedOverlay := resizeImage(overlayImg, overlayWidth, overlayHeight)

	// Create output image
	bounds := qrImg.Bounds()
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, qrImg, image.Point{}, draw.Src)

	// Calculate center position for overlay
	offsetX := (qrSize - overlayWidth) / 2
	offsetY := (qrSize - overlayHeight) / 2

	// Draw overlay centered on QR code
	overlayRect := image.Rect(offsetX, offsetY, offsetX+overlayWidth, offsetY+overlayHeight)
	draw.Draw(result, overlayRect, resizedOverlay, image.Point{}, draw.Over)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, result); err != nil {
		return nil, fmt.Errorf("encoding PNG: %w", err)
	}

	return buf.Bytes(), nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcWidth / width
			srcY := y * srcHeight / height
			result.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return result
}

func printHelp() {
	fmt.Println("QR Code Art Generator")
	fmt.Println("=====================")
	fmt.Println()
	fmt.Println("Creates QR codes with embedded or overlaid images.")
	fmt.Println()
	fmt.Println("Usage: qrmaker -image <path> [options]")
	fmt.Println()
	fmt.Println("Modes:")
	fmt.Println("  Embed (default): Image pattern is encoded into QR data bits")
	fmt.Println("  Overlay (-overlay): Image is placed on top of QR code center")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Embed mode (artistic QR)")
	fmt.Println("  qrmaker -image photo.png -url 'https://mysite.com'")
	fmt.Println("  qrmaker -image photo.png -url 'https://mysite.com' -dither -version 8")
	fmt.Println()
	fmt.Println("  # Overlay mode (logo in center)")
	fmt.Println("  qrmaker -overlay -image logo.png -url 'https://mysite.com'")
	fmt.Println("  qrmaker -overlay -image logo.png -url 'https://mysite.com' -size 1024")
	fmt.Println("  qrmaker -overlay -image logo.png -url 'https://mysite.com' -size 1024 -overlay-size 30")
	fmt.Println()
	fmt.Println("Styling Tips:")
	fmt.Println("  - Embed: Use -dither for photos, higher -version for longer URLs")
	fmt.Println("  - Overlay: Keep -overlay-size under 30% for reliable scanning")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}
