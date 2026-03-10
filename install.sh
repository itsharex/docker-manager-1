#!/bin/bash

set -e

APP_DIR="/opt/docker-manager"

BACKEND_IMAGE="ngthanhvu/docker-manager-backend:latest"
FRONTEND_IMAGE="ngthanhvu/docker-manager-frontend:latest"

DEFAULT_PORT=8088

if [ "$(id -u)" -eq 0 ]; then
    SUDO=""
else
    if ! command -v sudo >/dev/null 2>&1; then
        echo "This installer needs root privileges to write to $APP_DIR and manage Docker."
        echo "Please run it with sudo, for example:"
        echo "sudo bash <(curl -Ls https://raw.githubusercontent.com/ngthanhvu/docker-manager/refs/heads/main/install.sh)"
        exit 1
    fi
    SUDO="sudo"
fi

echo "================================="
echo " Docker Manager Installer"
echo "================================="
echo ""
echo "1) Install Docker Manager"
echo "2) Uninstall Docker Manager"
echo ""

read -p "Choose option [1-2]: " OPTION

# -------------------------
# INSTALL
# -------------------------

install_app() {

echo "Installing Docker Manager..."

# check docker
if ! command -v docker &> /dev/null
then
    echo "Docker not found. Installing..."
    curl -fsSL https://get.docker.com | $SUDO sh
fi

# check compose
if ! docker compose version &> /dev/null
then
    echo "Docker Compose v2 is required."
    exit 1
fi

# choose port
read -p "Enter port for Docker Manager [${DEFAULT_PORT}]: " PORT
PORT=${PORT:-$DEFAULT_PORT}

# check port
check_port() {
    if lsof -i:$1 >/dev/null 2>&1; then
        return 1
    else
        return 0
    fi
}

while ! check_port $PORT; do
    echo "Port $PORT is already in use."
    read -p "Choose another port: " PORT
done

echo "Using port: $PORT"

# create install dir
$SUDO mkdir -p "$APP_DIR"
cd "$APP_DIR"

# create compose file
$SUDO tee "$APP_DIR/docker-compose.yml" >/dev/null <<EOF
version: "3.8"

services:

  backend:
    image: $BACKEND_IMAGE
    container_name: docker-manager-backend
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - "127.0.0.1:8080:8080"

  frontend:
    image: $FRONTEND_IMAGE
    container_name: docker-manager-frontend
    restart: unless-stopped
    depends_on:
      - backend
    ports:
      - "$PORT:80"
EOF

echo "Pulling images..."
$SUDO docker compose -f "$APP_DIR/docker-compose.yml" pull

echo "Starting Docker Manager..."
$SUDO docker compose -f "$APP_DIR/docker-compose.yml" up -d

IP=$(hostname -I | awk '{print $1}')

echo ""
echo "================================="
echo " Docker Manager Installed"
echo "================================="
echo ""
echo "Access URL:"
echo "http://$IP:$PORT"
echo ""
echo "Install directory:"
echo "$APP_DIR"
echo ""

}

# -------------------------
# UNINSTALL
# -------------------------

uninstall_app() {

if [ ! -d "$APP_DIR" ]; then
    echo "Docker Manager is not installed."
    exit 0
fi

read -p "Are you sure you want to uninstall Docker Manager? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ]; then
    echo "Cancelled."
    exit 0
fi

cd "$APP_DIR"

echo "Stopping containers..."
$SUDO docker compose -f "$APP_DIR/docker-compose.yml" down

echo "Removing images..."
$SUDO docker image rm "$BACKEND_IMAGE" 2>/dev/null || true
$SUDO docker image rm "$FRONTEND_IMAGE" 2>/dev/null || true

echo "Removing install directory..."
$SUDO rm -rf "$APP_DIR"

echo ""
echo "Docker Manager removed successfully."
echo ""

}

# -------------------------
# MENU
# -------------------------

case $OPTION in
1)
    install_app
    ;;
2)
    uninstall_app
    ;;
*)
    echo "Invalid option."
    ;;
esac
