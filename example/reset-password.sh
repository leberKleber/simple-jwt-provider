#!/usr/bin/env sh

if [ "$#" -ne  "3" ]; then
   echo "Two arguments must be set e.g. ./reset-password.sh email password reset-token"
   exit 1
fi

curl -X POST --data "{\"email\":\"$1\", \"password\":\"$2\", \"reset_token\": \"$3\"}"  localhost:8080/v1/auth/password-reset -v
