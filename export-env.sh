#!/bin/sh

if [ -f .env ]; then
  # Load environment variables from .env file
  source .env
  echo "Environment variables loaded from .env file."
else
  echo ".env file not found."
  exit 1
fi