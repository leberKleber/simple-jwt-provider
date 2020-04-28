![Go](https://github.com/leberKleber/simple-jwt-provider/workflows/Go/badge.svg?branch=master)

# simple-jwt-provider

## Try it

```bash
git clone git@github.com:leberKleber/simple-jwt-provider.git
docker-compose -f example/docker-compose.yml up

# create user via admin-api
./example/create_user.sh test.test@test.test password

# login with created user
./example/login.sh test.test@test.test password

# reset password
# 1) create password reset request
#    - mail with reset token would be send
# 2) reset password with received token

# 1) create password reset request 
./example/create_password-reset-request.sh test.test@test.test
# 1.1) open browser at http://127.0.0.1:8025/ and copy reset token (token only not the url)
# 2) reset password with received token
./example/reset-password.sh test.test@test.test newPassword {reset-token}
# verify new password
./example/login.sh test.test@test.test newPassword

```


## Getting started
### Generate ECDSA-512 key pair

```sh
# private key
openssl ecparam -genkey -name secp521r1 -noout -out ecdsa-p521-private.pem
# public key
openssl ec -in ecdsa-p521-private.pem -pubout -out ecdsa-p521-public.pem 
```

## API
### POST `/v1/auth/login`
This endpoint will check the email/password combination and will set the respond with an jwtauthToken if correct:

Request body:
```json
{
    "email": "info@leberkleber.io",
    "password": "s3cr3t"
}
```

Response body (200 - OK):
```json
{
    "access_token":"<jwt>"
}
```

### POST `/v1/auth/password-reset-request`
This endpoint will trigger a password reset request. The user gets a token per mail.
With this token, the password can be reset via POST@`/v1/auth/password-reset` .

Request body:
```json
{
    "email": "info@leberkleber.io"
}
```

Response (201 - CREATED)

### POST `/v1/auth/password-reset`
This endpoint will reset the password of the given user if the reset-token is valid and matches to the given email.

Request body:
```json
{
    "email": "info@leberkleber.io",
    "reset_token": "rAnDoMsHiT456",
    "password": "SeCReT"
}
```

Response (200 - OK)


### POST `/v1/admin/users`
This endpoint will create an new user if admin api auth was successfully:

Request body:
```json
{
    "email": "info@leberkleber.io",
    "password": "s3cr3t"
}
```

Response body (201 - CREATED)
