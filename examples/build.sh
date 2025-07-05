#!/bin/bash

set -e

# Ensure static directory exists
mkdir -p static

# 1. Generate Templ Go code
echo "Generating Templ Go code..."
templ generate -path templates

# 2. Generate GoCSS
echo "Generating GoCSS..."
# Use the local gocss binary
../cmd/gocss/main.go --input "./templates/*.templ" --output ./static/gocss.css

# 3. Run the Go application
echo "Running Go application..."
go run .
