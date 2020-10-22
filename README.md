![Go](https://github.com/leberKleber/simple-jwt-provider/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/leberKleber/simple-jwt-provider)](https://goreportcard.com/report/github.com/leberKleber/simple-jwt-provider)

# simple-jwt-provider
Simple and lightweight JWT-Provider written in go (golang). It exhibits JWT for the in postgres
persisted user, which can be managed via api. Also, a password-reset flow via mail verification is available.
User specific custom-claims also available for jwt-generation and mail rendering.

dockerized: https://hub.docker.com/r/leberkleber/simple-jwt-provider

## Try it
```bash
git clone git@github.com:leberKleber/simple-jwt-provider.git
docker-compose -f example/docker-compose.yml up

# create user via admin-api
./example/create_user.sh test.test@test.test password

# login with created user
./example/login.sh test.tscest@test.test password

# reset password
# 1) create password reset request
#    - mail with reset token would be send
# 2) reset password with received token
# 3) do crud operations on user

# 1) create password reset request 
./example/create_password-reset-request.sh test.test@test.test
# 1.1) open browser at http://127.0.0.1:8025/ and copy reset token (token only not the url)
# 2) reset password with received token
./example/reset-password.sh test.test@test.test newPassword {reset-token}
# verify new password
./example/login.sh test.test@test.test newPassword

# 3) do crud operations on user
# see ./example/*.sh

```

## Getting started
### Generate ECDSA-512 key pair

```sh
# private key
openssl ecparam -genkey -name secp521r1 -noout -out ecdsa-p521-private.pem
# public key
openssl ec -in ecdsa-p521-private.pem -pubout -out ecdsa-p521-public.pem 
```

### Configuration
| Environment variable              | Description                                                         | Required                            | Default               |
| --------------------------------- |:-------------------------------------------------------------------:| -----------------------------------:|----------------------:|
| SJP_SERVER_ADDRESS                | Server-address network-interface to bind on e.g.: '127.0.0.1:8080'  | no                                  | 0.0.0.0:80            |
| SJP_JWT_PRIVATE_KEY               | JWT PrivateKey ECDSA512                                             | yes                                 | -                     |
| SJP_JWT_AUDIENCE                  | Audience private claim which will be applied in each JWT            | no                                  | -                     |
| SJP_JWT_ISSUER                    | Issuer private claim which will be applied in each JWT              | no                                  | -                     |
| SJP_JWT_SUBJECT                   | Subject private claim which will be applied in each JWT             | no                                  | -                     |
| SJP_DB_HOST                       | Database-Host (postgres)                                            | yes                                 | -                     |
| SJP_DB_PORT                       | Database-Port                                                       | no                                  | 5432                  |
| SJP_DB_NAME                       | Database-Name                                                       | no                                  | simple-jwt-provider   |
| SJP_DB_USERNAME                   | Database-Username                                                   | no                                  | -                     |
| SJP_DB_PASSWORD                   | Database-Password                                                   | no                                  | -                     |
| SJP_MIGRATIONS_FOLDER_PATH        | Database Migrations Folder Path                                     | no                                  | /db-migrations        |
| SJP_ADMIN_API_ENABLE              | Enable admin API to manage stored users (true / false)              | no                                  | false                 |
| SJP_ADMIN_API_USERNAME            | Basic Auth Username if enable-admin-api = true                      | yes, when enable-admin-api = true   | -                     |
| SJP_ADMIN_API_PASSWORD            | Basic Auth Password if enable-admin-api = true                      | yes, when enable-admin-api = true   | -                     |
| SJP_MAIL_TEMPLATES_FOLDER_PATH    | Path to mail-templates folder                                       | no                                  | /mail-templates       |
| SJP_MAIL_SMTP_HOST                | SMTP host to connect to                                             | yes                                 | -                     |
| SJP_MAIL_SMTP_PORT                | SMTP port to connect to                                             | no                                  | 587                   |
| SJP_MAIL_SMTP_USERNAME            | SMTP username to authorize with                                     | yes                                 | -                     |
| SJP_MAIL_SMTP_PASSWORD            | SMTP password to authorize with                                     | yes                                 | -                     |
| SJP_MAIL_TLS_INSECURE_SKIP_VERIFY | true if certificates should not be verified                         | no                                  | false                 |
| SJP_MAIL_TLS_SERVER_NAME          | name of the server who expose the certificate                       | no                                  | -                     |

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
This endpoint will create a new user if admin api auth was successfully:

Request body:
```json
{
    "email": "info@leberkleber.io",
    "password": "s3cr3t",
    "claims":  {
        "myCustomClaim": "custom claims for jwt and mail templates"
    }
}
```

Response body (201 - CREATED)

### PUT `/v1/admin/users/{email}`
This endpoint will update the given properties (excluding email) of the user with the given email when the admin api auth was successfully:

Request body:
```json
{
    "password": "n3wS3cr3t",
    "claims":  {
        "updatedClaim": "now updated"
    }
}
```

Response body (200 - NO CONTENT)
```json
{
    "email": "info@leberkleber.io",
    "password": "**********",
    "claims":  {
        "updatedClaim": "now updated"
    }
}
```

### DELETE `/v1/admin/users/{email}`
This endpoint will delete the user with the given email when there are no tokens which referred to this user and the admin api auth was successfully:

Response body (201 - NO CONTENT)

