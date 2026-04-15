set -e

KEY_DIR="/home/martin-wachira/Martin/storeforge/usermanagement/internal/keys"

# Ensure directory exists
mkdir -p "$KEY_DIR"
chmod 700 "$KEY_DIR"

if [ ! -f "$KEY_DIR/jwt_private.pem" ]; then
  echo "Generating JWT keys..."
  openssl genpkey \
    -algorithm RSA \
    -out "$KEY_DIR/jwt_private.pem" \
    -pkeyopt rsa_keygen_bits:2048

  openssl pkey \
    -in "$KEY_DIR/jwt_private.pem" \
    -pubout \
    -out "$KEY_DIR/jwt_public.pem"

  chmod 600 "$KEY_DIR/jwt_private.pem"
  chmod 644 "$KEY_DIR/jwt_public.pem"
else
  echo "JWT keys already exist, skipping generation."
fi

exec "$@"