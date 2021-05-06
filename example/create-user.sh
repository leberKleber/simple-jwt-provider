#!/usr/bin/env sh

if [ "$#" -ne "3" ]; then
  echo "Three arguments must be set e.g. ./create-user.sh email password claim(json)"
  exit 1
fi

curl -X POST --data "{\"email\":\"$1\", \"password\":\"$2\", \"claims\":$3}" "username:password@localhost:8080/v1/admin/users" -v
