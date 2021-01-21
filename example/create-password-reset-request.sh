#!/usr/bin/env sh

if [ "$#" -ne "1" ]; then
  echo "One argument must be set e.g. ./create-password-reset-request.sh email"
  exit 1
fi
curl -X POST --data "{\"email\":\"$1\"}" "localhost:8080/v1/auth/password-reset-request" -v
