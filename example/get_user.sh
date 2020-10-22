#!/usr/bin/env sh

if [ "$#" -ne  "1" ]; then
	echo "One argument must be set e.g. ./get_user.sh email"
   exit 1
fi

DIR="$(cd "$(dirname "$0")" && pwd)"
urlencodedEMail=$("$DIR"/urlencode.sh "$1")

curl -X GET "username:password@localhost:8080/v1/admin/users/$urlencodedEMail" -v
