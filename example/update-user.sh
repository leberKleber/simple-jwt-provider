#!/usr/bin/env sh

if [ "$#" -ne "2" ]; then
  echo "Two arguments must be set e.g. ./update-user.sh email user(json)"
  exit 1
fi

curl -X PUT --data "$2" "username:password@localhost:8080/v1/admin/users/$1" -v
