# Docker Setup für Jurassic Park SNES Randomizer

## Voraussetzungen

- Docker
- Docker Compose

## Starten der Anwendung

1. Stelle sicher, dass die Go-Binaries vorhanden sind:
   - `internal/tools/rncpropack/rnc64` oder `rnc32`
   - `internal/tools/flips/flips`

2. Starte die Services:
```bash
docker-compose up -d
```

3. Die Anwendung ist verfügbar unter:
   - Laravel Frontend: http://localhost:8000
   - Go API: http://localhost:8080

## Services

- **laravel**: Laravel Frontend (Port 8000)
- **randomizer-api**: Go Randomizer API (Port 8080)

## Logs anzeigen

```bash
# Alle Logs
docker-compose logs -f

# Nur Laravel
docker-compose logs -f laravel

# Nur Go API
docker-compose logs -f randomizer-api
```

## Services stoppen

```bash
docker-compose down
```

## Services neu bauen

```bash
docker-compose build --no-cache
docker-compose up -d
```

## Laravel-spezifische Befehle

```bash
# Composer installieren
docker-compose exec laravel composer install

# App-Key generieren
docker-compose exec laravel php artisan key:generate

# Cache leeren
docker-compose exec laravel php artisan cache:clear
docker-compose exec laravel php artisan config:clear
```

