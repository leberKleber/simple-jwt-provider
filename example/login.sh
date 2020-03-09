#!/usr/bin/env sh

if [ "$#" -ne  "2" ]; then
   echo "Two arguments must be set e.g. ./login.sh email password"
   exit 1
fi

curl -X POST --data "{\"email\":\"$1\", \"password\":\"$2\"}" "username:password@localhost:8080/v1/auth/login" -v
