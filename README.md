# qrmaker

A CLI tool for generating QR codes with embedded or overlaid images.

## Features

- **Embed Mode**: Encode image patterns directly into QR code data bits, creating artistic QR codes where the image becomes part of the code itself
- **Overlay Mode**: Place a logo or image in the center of a standard QR code
- Supports PNG, JPEG, and GIF input images
- Configurable QR code version, scale, and output size
- Dithering option for better photo representation
- Reproducible results with seed control

## Installation

```bash
go install github.com/ohitsdaniel/qrmaker@latest
```

Or build from source:

```bash
git clone https://github.com/ohitsdaniel/qrmaker.git
cd qrmaker
go build -o qrmaker .
```

## Usage

```bash
qrmaker -image <path> [options]
```

### Embed Mode (default)

Creates artistic QR codes where the image pattern is encoded into the QR data bits:

```bash
# Basic embed
qrmaker -image photo.png -url 'https://example.com'

# With dithering (better for photos)
qrmaker -image photo.png -url 'https://example.com' -dither -version 8

# High resolution
qrmaker -image photo.png -url 'https://example.com' -scale 16 -version 8
```

### Overlay Mode

Places an image on top of the QR code center (relies on QR error correction):

```bash
# Basic overlay
qrmaker -overlay -image logo.png -url 'https://example.com'

# Custom size
qrmaker -overlay -image logo.png -url 'https://example.com' -size 1024

# Adjust logo size (percentage of QR code)
qrmaker -overlay -image logo.png -url 'https://example.com' -overlay-size 30
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-image` | (required) | Path to the image to embed/overlay |
| `-url` | `https://example.com` | URL or text to encode |
| `-output` | `qrcode.png` | Output file path |
| `-overlay` | `false` | Use overlay mode instead of embed |
| `-version` | `6` | QR code version (1-8 for embed) |
| `-scale` | `8` | Size of each QR module in pixels |
| `-size` | `512` | Output image size in pixels (overlay mode) |
| `-overlay-size` | `25` | Logo size as percentage of QR code (10-40) |
| `-dither` | `false` | Enable dithering (embed mode) |
| `-rand` | `false` | Random control for pixel placement (embed mode) |
| `-seed` | `0` | Random seed (0 = use current time) |
| `-mask` | `2` | QR mask pattern 0-7 (embed mode) |
| `-dx`, `-dy` | `4` | Image positioning offset (embed mode) |

## Tips

- **Embed mode**: Use `-dither` for photographs, increase `-version` for longer URLs
- **Overlay mode**: Keep `-overlay-size` under 30% for reliable scanning
- Always test that your generated QR codes scan correctly

## License

MIT
