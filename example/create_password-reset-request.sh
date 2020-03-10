#!/usr/bin/env sh

curl -X POST --data "{\"email\":\"$1\"}" "localhost:8080/v1/auth/password-reset-request" -v
