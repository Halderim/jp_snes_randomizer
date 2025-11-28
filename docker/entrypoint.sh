#!/bin/sh

cd /var/www/html

# Installiere Composer-Abhängigkeiten falls nicht vorhanden
if [ ! -d "vendor" ]; then
    echo "Installing Composer dependencies..."
    composer install --no-interaction --prefer-dist --optimize-autoloader || true
fi

# Setze Berechtigungen
chown -R www-data:www-data /var/www/html
chmod -R 755 /var/www/html/storage
chmod -R 755 /var/www/html/bootstrap/cache

# Erstelle .env falls nicht vorhanden
if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    touch .env
    # Füge minimale .env Einstellungen hinzu
    cat >> .env <<EOF
APP_NAME=Laravel
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_URL=http://localhost

DB_CONNECTION=sqlite

RANDOMIZER_API_URL=http://randomizer-api:8080
EOF
fi

# Stelle sicher, dass APP_CIPHER gesetzt ist (für Laravel 10+)
if ! grep -q "^APP_CIPHER=" .env 2>/dev/null; then
    echo "APP_CIPHER=AES-256-CBC" >> .env
fi

# Prüfe ob APP_KEY gesetzt ist (muss base64: Präfix haben)
if ! grep -q "^APP_KEY=base64:" .env 2>/dev/null || grep -q "^APP_KEY=$" .env 2>/dev/null; then
    echo "Generating application key..."
    # Entferne alte APP_KEY Zeile falls vorhanden
    sed -i '/^APP_KEY=/d' .env 2>/dev/null || true
    # Generiere neuen Key
    php artisan key:generate --force --show
    # Stelle sicher, dass der Key in .env steht
    if ! grep -q "^APP_KEY=base64:" .env 2>/dev/null; then
        echo "Warning: Key generation may have failed, trying again..."
        php artisan key:generate --force
    fi
fi

# Cache leeren für frische Konfiguration
php artisan config:clear 2>/dev/null || true
php artisan cache:clear 2>/dev/null || true

exec "$@"

