#!/usr/bin/env bash
# Build manual a WebAssembly para publicar el juego en un server propio.
# Para desarrollo local rápido, usá wasmserve en cambio (ver README).
set -euo pipefail

cd "$(dirname "$0")/.."

echo "Compilando a WASM..."
GOOS=js GOARCH=wasm go build -o web/f1game.wasm ./cmd/f1game

GOVERSION_MAJOR_MINOR=$(go env GOVERSION | sed 's/go//' | cut -d. -f1-2)
echo "Copiando wasm_exec.js (Go $GOVERSION_MAJOR_MINOR)..."

# Go 1.24+ movió wasm_exec.js de misc/wasm a lib/wasm.
if [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then
	cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js
else
	cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/wasm_exec.js
fi

echo "Listo. Serví la carpeta web/ con cualquier servidor estático, ej:"
echo "  cd web && python3 -m http.server 8080"
echo "y abrí http://localhost:8080"
