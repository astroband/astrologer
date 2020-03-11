#!/bin/sh

set -ue

astrologer_create_indices() {
  local INDICES_INITIALIZED="/root/.indices-initialized"

  if [ -f "$INDICES_INITIALIZED" ]; then
    echo "ES indices have already been initialized."
    return 0
  fi

  echo "Initializing ES indices..."

  /root/astrologer create-index

  echo "Finished initializing ES indices"

  touch $INDICES_INITIALIZED
}

astrologer_create_indices

exec "$@"
