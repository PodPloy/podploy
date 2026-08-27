#!/bin/sh

set -euo pipefail

CONTAINER_NAME="podploy_db"
VOLUME_NAME="podploy_data"
IMAGE_NAME="docker.io/library/postgres:16-alpine"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="postgres"

if podman ps --format "{{.Names}}" | grep -Eq "^${CONTAINER_NAME}\$"; then
    echo "[INFO] Container '${CONTAINER_NAME}' is already running."
    exit 0
fi

if podman ps -a --format "{{.Names}}" | grep -Eq "^${CONTAINER_NAME}\$"; then
    echo "[INFO] Container '${CONTAINER_NAME}' exists but is stopped. Starting..."
    podman start "${CONTAINER_NAME}"
    echo "[SUCCESS] Container '${CONTAINER_NAME}' has been started."
    exit 0
fi

echo "[INFO] Container '${CONTAINER_NAME}' not found. Initializing new deployment..."

if ! podman volume exists "${VOLUME_NAME}"; then
    echo "[INFO] Creating persistent volume '${VOLUME_NAME}'..."
    podman volume create "${VOLUME_NAME}"
fi

echo "[INFO] Launching container '${CONTAINER_NAME}'..."
podman run -d \
    --name "${CONTAINER_NAME}" \
    -e POSTGRES_USER="${DB_USER}" \
    -e POSTGRES_PASSWORD="${DB_PASSWORD}" \
    -e POSTGRES_DB="${DB_NAME}" \
    -p "${DB_PORT}:5432" \
    -v "${VOLUME_NAME}:/var/lib/postgresql/data:Z" \
    --restart unless-stopped \
    "${IMAGE_NAME}"

echo "[SUCCESS] PostgreSQL container '${CONTAINER_NAME}' created and running on port ${DB_PORT}."
