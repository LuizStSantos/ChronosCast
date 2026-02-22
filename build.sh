#!/bin/bash
VERSION="v1.1.0"
NAME="chronoscast-linux-amd64"
OUTPUT="${NAME}-${VERSION}"

echo "🔨 Compilando o Chronoscast..."
go build -ldflags="-s -w" -o $OUTPUT cmd/server/main.go

echo "--------------------------------------------------"
echo "✅ Sucesso! build completa $OUTPUT."