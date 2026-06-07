#!/bin/sh

set -eou pipefail

cd "$(dirname "$0")"

OUTPUT_DIR="${OUTPUT_DIR:-/data/public}"
OUTPUT_FILE="${OUTPUT_DIR}/dashboard.jpg"

if [ -z "${SYNC_CLOUD_STORAGE}" ]; then echo "error: SYNC_CLOUD_STORAGE is not set"; exit 1; fi
if [ "${SYNC_CLOUD_STORAGE}" ]; then
  if [ -z "${RCLONE_DESTINATION}" ]; then echo "error: RCLONE_DESTINATION is not set"; exit 1; fi
  if [ -z "${RCLONE_CONFIG}" ]; then echo "error: RCLONE_CONFIG is not set"; exit 1; fi
fi

mkdir -p "${OUTPUT_DIR}"

# generate dashboard images (light + derived dark variant)
/app/m5paper-dashboard --output "${OUTPUT_FILE}"

if [ -n "${SYNC_CLOUD_STORAGE}" ]; then
  # setup rclone (write to /tmp since the working dir /app is not writable by the non-root user)
  RCLONE_CONF="$(mktemp)"
  echo "$RCLONE_CONFIG" | base64 -d > "${RCLONE_CONF}"

  # upload
  /usr/local/bin/rclone -v sync --config "${RCLONE_CONF}" "${OUTPUT_DIR}" "${RCLONE_DESTINATION}"
fi
