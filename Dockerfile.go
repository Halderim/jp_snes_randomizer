# Dockerfile für Go Randomizer API

FROM golang:1.22 AS builder

# Installiere notwendige Tools
RUN apt-get update && apt-get install -y \
    git

# Arbeitsverzeichnis
WORKDIR /app

# Kopiere go mod files
COPY go.mod ./
RUN go mod download

# Kopiere den gesamten Code
COPY . .

# Baue die Anwendung
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o randomizer-api ./cmd/web

# Finales Image
FROM ubuntu:latest

# Installiere ca-certificates für HTTPS-Requests
RUN apt-get update && apt-get install -y ca-certificates libgtk-3-0 libgomp1

WORKDIR /app

# Kopiere die Binaries für rncpropack und flips
COPY --from=builder /app/internal/tools/rncpropack/rnc64 ./internal/tools/rncpropack/rnc64
COPY --from=builder /app/internal/tools/rncpropack/rnc32 ./internal/tools/rncpropack/rnc32
COPY --from=builder /app/internal/tools/flips/flips ./internal/tools/flips/flips

RUN chmod +x ./internal/tools/rncpropack/rnc64
RUN chmod +x ./internal/tools/rncpropack/rnc32
RUN chmod +x ./internal/tools/flips/flips

# Kopiere die kompilierte Anwendung
COPY --from=builder /app/randomizer-api .

# Kopiere notwendige Verzeichnisse
COPY internal/uncompressed ./internal/uncompressed
COPY internal/patches ./internal/patches
COPY internal/rom/unmodified ./internal/rom/unmodified

# Erstelle notwendige Verzeichnisse
RUN mkdir -p internal/logs internal/modbin internal/rom/modified

# Exponiere Port 8080
EXPOSE 8080

# Starte die Anwendung
CMD ["./randomizer-api"]

