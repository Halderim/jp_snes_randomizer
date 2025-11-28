# Docker Setup für Jurassic Park SNES Randomizer

## Schnellstart

```bash
# Services starten
docker-compose up -d

# Logs anzeigen
docker-compose logs -f

# Services stoppen
docker-compose down
```

## Zugriff

- **Laravel Frontend**: http://localhost:8000
- **Go API**: http://localhost:8080

## Voraussetzungen

- Docker Desktop (Windows/Mac) oder Docker Engine + Docker Compose (Linux)
- Die Binaries müssen vorhanden sein:
  - `internal/tools/rncpropack/rnc64` (oder `rnc32` für 32-bit)
  - `internal/tools/flips/flips`

## Services

### randomizer-api (Go)
- Port: 8080
- Baut die Go-Anwendung und startet den Web-Server
- Verwendet Volumes für `internal/` Verzeichnis

### laravel (PHP)
- Port: 8000
- Laravel Frontend mit Nginx und PHP-FPM
- Automatische Installation von Composer-Abhängigkeiten
- Automatische Generierung des App-Keys

## Wichtige Befehle

### Logs anzeigen
```bash
# Alle Services
docker-compose logs -f

# Nur Laravel
docker-compose logs -f laravel

# Nur Go API
docker-compose logs -f randomizer-api
```

### Container neu bauen
```bash
docker-compose build --no-cache
docker-compose up -d
```

### Laravel-Befehle ausführen
```bash
# Composer installieren
docker-compose exec laravel composer install

# App-Key generieren
docker-compose exec laravel php artisan key:generate

# Cache leeren
docker-compose exec laravel php artisan cache:clear
docker-compose exec laravel php artisan config:clear
docker-compose exec laravel php artisan view:clear
```

### In Container einloggen
```bash
# Laravel Container
docker-compose exec laravel sh

# Go API Container
docker-compose exec randomizer-api sh
```

## Troubleshooting

### Laravel zeigt Fehler 500
1. Prüfe die Logs: `docker-compose logs laravel`
2. Generiere App-Key: `docker-compose exec laravel php artisan key:generate`
3. Setze Berechtigungen: `docker-compose exec laravel chmod -R 755 storage bootstrap/cache`

### Go API startet nicht
1. Prüfe die Logs: `docker-compose logs randomizer-api`
2. Stelle sicher, dass die Binaries vorhanden sind
3. Prüfe ob Port 8080 frei ist

### Services können nicht kommunizieren
- Prüfe ob beide Services im gleichen Netzwerk sind: `docker network ls`
- Prüfe die Umgebungsvariable `RANDOMIZER_API_URL` in `docker-compose.yml`

## Entwicklung

Für die Entwicklung werden Volumes verwendet, sodass Code-Änderungen sofort sichtbar sind:

- `./laravel` → `/var/www/html` (Laravel)
- `./internal` → `/app/internal` (Go API)

## Produktion

Für Produktion sollten:
- `APP_DEBUG=false` gesetzt werden
- `APP_ENV=production` gesetzt werden
- Volumes entfernt oder read-only gemacht werden
- Secrets über Environment-Variablen oder Docker Secrets verwaltet werden

