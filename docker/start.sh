#!/bin/bash

echo "🚀 Starting Jurassic Park SNES Randomizer..."

# Prüfe ob Docker läuft
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker ist nicht gestartet. Bitte starte Docker zuerst."
    exit 1
fi

# Baue und starte die Container
echo "📦 Building and starting containers..."
docker-compose up -d --build

# Warte auf Laravel
echo "⏳ Waiting for Laravel to be ready..."
sleep 5

# Installiere Laravel-Abhängigkeiten falls nötig
echo "📥 Installing Laravel dependencies..."
docker-compose exec -T laravel composer install --no-interaction || true

# Generiere App-Key falls nötig
echo "🔑 Generating Laravel app key..."
docker-compose exec -T laravel php artisan key:generate --force || true

# Setze Berechtigungen
echo "🔐 Setting permissions..."
docker-compose exec -T laravel chown -R www-data:www-data /var/www/html/storage /var/www/html/bootstrap/cache || true
docker-compose exec -T laravel chmod -R 755 /var/www/html/storage /var/www/html/bootstrap/cache || true

echo ""
echo "✅ Setup complete!"
echo ""
echo "🌐 Laravel Frontend: http://localhost:8000"
echo "🔧 Go API: http://localhost:8080"
echo ""
echo "📋 View logs: docker-compose logs -f"
echo "🛑 Stop services: docker-compose down"

