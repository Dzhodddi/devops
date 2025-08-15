#!bin/bash

# Load .env file
#!/bin/bash

# Load .env file safely
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

migrate -path internal/db/migrations -database "$DB_ADDR_LOCAL" up


