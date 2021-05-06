#!/usr/bin/env sh

if [ "$#" -ne "1" ]; then
  echo "One argument must be set e.g. ./refresh.sh refreshToken"
  exit 1
fi

curl -X POST --data "{\"refresh_token\":\"$1\"}" "username:password@localhost:8080/v1/auth/refresh" -v
