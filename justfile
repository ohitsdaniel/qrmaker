# QR Code Art Generator

# Default recipe - show help
default:
    @just --list

# Build the qrmaker binary
build:
    go build -o qrmaker .

# Run with an image (usage: just run image.png)
run image url="https://example.com":
    go run . -image {{image}} -url "{{url}}"

# Generate QR with custom output name
generate image url output="qrcode.png":
    go run . -image {{image}} -url "{{url}}" -output {{output}}

# Generate with dithering (better for photos)
dither image url="https://example.com":
    go run . -image {{image}} -url "{{url}}" -dither -version 8

# Generate with random artistic style
artistic image url="https://example.com":
    go run . -image {{image}} -url "{{url}}" -rand -dither -version 8

# Generate high-res QR code
highres image url="https://example.com":
    go run . -image {{image}} -url "{{url}}" -scale 16 -version 8

# Overlay logo on QR code (instead of embedding)
overlay image url="https://example.com":
    go run . -overlay -image {{image}} -url "{{url}}"

# Overlay with custom resolution
overlay-hd image url="https://example.com":
    go run . -overlay -image {{image}} -url "{{url}}" -size 1024

# Overlay with full customization
overlay-custom image url="https://example.com" size="512" logo="25":
    go run . -overlay -image {{image}} -url "{{url}}" -size {{size}} -overlay-size {{logo}}

# Show help
help:
    go run .

# Clean build artifacts
clean:
    rm -f qrmaker *.png

# Install globally
install:
    go install .

# Run tests with the sample image
test:
    @if [ -f extitutional.png ]; then \
        go run . -image extitutional.png -url "https://github.com" -output test_output.png && \
        echo "Test passed! Generated test_output.png"; \
    else \
        echo "No test image found. Place an image named extitutional.png in the directory."; \
    fi
