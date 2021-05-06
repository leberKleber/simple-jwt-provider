[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Go](https://github.com/leberKleber/simple-jwt-provider/workflows/Go/badge.svg?branch=master)](https://github.com/leberKleber/simple-jwt-provider/actions?query=workflow%3Ago)
[![Go Report Card](https://goreportcard.com/badge/github.com/leberKleber/simple-jwt-provider)](https://goreportcard.com/report/github.com/leberKleber/simple-jwt-provider)
[![codecov](https://codecov.io/gh/leberKleber/simple-jwt-provider/branch/master/graph/badge.svg)](https://codecov.io/gh/leberKleber/simple-jwt-provider)

# simple-jwt-provider

Simple and lightweight JWT-Provider written in go (golang). It exhibits JWT for the in postgres persisted user, which
can be managed via api. Also, a password-reset flow via mail verification is available. User specific custom-claims also
available for jwt-generation and mail rendering.

dockerized: https://hub.docker.com/r/leberkleber/simple-jwt-provider

build it yourself:

```shell script
# as docker-image
docker build . -t leberkleber/simple-jwt-provider

# as binary
go build -o simple-jwt-provider ./cmd/provider/
```

# Table of contents

- [Try it](#try-it)
- [Getting started](#getting-started)
    - [Generate ECDSA-512 key pair](#generate-ecdsa-512-key-pair)
    - [Configuration](#configuration)
- [API](#api)
    - [POST `/v1/auth/login`](#post-v1authlogin)
    - [POST `/v1/auth/refresh`](#post-v1authrefresh)
    - [POST `/v1/auth/password-reset-request`](#post-v1authpassword-reset-request)
    - [POST `/v1/auth/password-reset`](#post-v1authpassword-reset)
    - [POST `/v1/admin/users`](#post-v1adminusers)
    - [PUT `/v1/admin/users/{email}`](#put-v1adminusersemail)
    - [DELETE `/v1/admin/users/{email}`](#delete-v1adminusersemail)
- [Mail](#mail)
    - [Password reset request](#password-reset-request)
- [Development](#development)
    - [mocks](#mocks)
    - [component tests](#component-tests)

## Try it

```shell script
git clone git@github.com:leberKleber/simple-jwt-provider.git
docker-compose -f example/docker-compose.yml up

# create user via admin-api
./example/create-user.sh test.test@test.test password {}

# login with created user
./example/login.sh test.tscest@test.test password

# reset password
# 1) create password reset request
#    - mail with reset token would be send
# 2) reset password with received token
# 3) do crud operations on user

# 1) create password reset request 
./example/create-password-reset-request.sh test.test@test.test
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

| Environment variable              | Description                                                                           | Required                            | Default               |
| --------------------------------- |:-------------------------------------------------------------------------------------:| -----------------------------------:|----------------------:|
| SJP_LOG_LEVEL                     | Log-Level can be TRACE DEBUG INFO WARN ERROR FATAL or PANIC                           | no                                  | INFO                  |
| SJP_SERVER_ADDRESS                | Server-address network-interface to bind on e.g.: '127.0.0.1:8080'                    | no                                  | 0.0.0.0:80            |
| SJP_JWT_LIFETIME                  | Lifetime of JWT                                                                       | no                                  | 4h                    |
| SJP_JWT_PRIVATE_KEY               | JWT PrivateKey ECDSA512                                                               | yes                                 | -                     |
| SJP_JWT_AUDIENCE                  | Audience private claim which will be applied in each JWT                              | no                                  | -                     |
| SJP_JWT_ISSUER                    | Issuer private claim which will be applied in each JWT                                | no                                  | -                     |
| SJP_JWT_SUBJECT                   | Subject private claim which will be applied in each JWT                               | no                                  | -                     |
| SJP_DSN                           | Data Source Name for persistence                                                      | yes                                 | -                     |
| SJP_MIGRATIONS_FOLDER_PATH        | Database Migrations Folder Path                                                       | no                                  | /db-migrations        |
| SJP_ADMIN_API_ENABLE              | Enable admin API to manage stored users (true / false)                                | no                                  | false                 |
| SJP_ADMIN_API_USERNAME            | Basic Auth Username if enable-admin-api = true                                        | yes, when enable-admin-api = true   | -                     |
| SJP_ADMIN_API_PASSWORD            | Basic Auth Password if enable-admin-api = true when is bcrypted prefix with 'bcrypt:' | yes, when enable-admin-api = true   | -                     |
| SJP_MAIL_TEMPLATES_FOLDER_PATH    | Path to mail-templates folder                                                         | no                                  | /mail-templates       |
| SJP_MAIL_SMTP_HOST                | SMTP host to connect to                                                               | yes                                 | -                     |
| SJP_MAIL_SMTP_PORT                | SMTP port to connect to                                                               | no                                  | 587                   |
| SJP_MAIL_SMTP_USERNAME            | SMTP username to authorize with                                                       | yes                                 | -                     |
| SJP_MAIL_SMTP_PASSWORD            | SMTP password to authorize with                                                       | yes                                 | -                     |
| SJP_MAIL_TLS_INSECURE_SKIP_VERIFY | true if certificates should not be verified                                           | no                                  | false                 |
| SJP_MAIL_TLS_SERVER_NAME          | name of the server who expose the certificate                                         | no                                  | -                     |

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
  "access_token": "<access-jwt>",
  "refresh_token": "<refresh-jwt>"
}
```

### POST `/v1/auth/refresh`

This endpoint will return a new access and refresh token. The submitted refresh-token will no longer be valid.

Request body:
```json
{
  "refresh_token": "<refresh_jwt>"
}
```

Response body (200 - OK):
```json
{
  "access_token": "<new-access-jwt>",
  "refresh_token": "<new-refresh-jwt>"
}
```
### POST `/v1/auth/password-reset-request`

This endpoint will trigger a password reset request. The user gets a token per mail. With this token, the password can
be reset via POST@`/v1/auth/password-reset`.

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

Response (204 - NO CONTENT)

### POST `/v1/admin/users`

This endpoint will create a new user if admin api auth was successfully:

Request body:
```json
{
  "email": "info@leberkleber.io",
  "password": "s3cr3t",
  "claims": {
    "myCustomClaim": "custom claims for jwt and mail templates"
  }
}
```

Response body (201 - CREATED)

### PUT `/v1/admin/users/{email}`

This endpoint will update the given properties (excluding email) of the user with the given email when the admin api
auth was successfully:

Request body:
```json
{
  "password": "n3wS3cr3t",
  "claims": {
    "updatedClaim": "now updated"
  }
}
```

Response body (200 - NO CONTENT)

```json
{
  "email": "info@leberkleber.io",
  "password": "**********",
  "claims": {
    "updatedClaim": "now updated"
  }
}
```

### DELETE `/v1/admin/users/{email}`

This endpoint will delete the user with the given email when there are no tokens which referred to this user, and the
admin api auth was successfully:

Response body (201 - NO CONTENT)

## Mail

Mails will be generated based on a set of templates which should be prepared for productive usage.

- `<mailType>.html` represents the html body of the mail and can be templated with `html.template` syntax
  (https://golang.org/pkg/html/template/). Available templating arguments listed in detailed template type description.
- `<mailType>.txt` represents the text body of the mail and can be templated with `text.template` syntax
  (https://golang.org/pkg/text/template/). Available templating arguments listed in detailed template type description.
- `<mailType>.yml` represents the header of the mail. In this template headers e.g. `From`, `To` or `Subject`
  can be set `text.template` syntax (https://golang.org/pkg/text/template/). Available templating arguments listed in
  detailed template type description.

### Password reset request

An example of this mail type can be found in `/mail-templates/password-reset-request.*`. Available template arguments:

| Argument           | Content                                                | Example usage                       |
|--------------------|--------------------------------------------------------|-------------------------------------|
| Recipient          | Users email address                                    | `{{.Recipient}}`                    |
| PasswordResetToken | The token which is required to reset the password      | `{{.PasswordResetToken}}`           |
| Claims             | All custom-claims which stored in relation to the user | `{{if index .Claims "first_name"}}` |

## Development

### mocks

Mocks will be generated with github.com/matryer/moq. Execute the following for generation:

```shell script
go get github.com/matryer/moq
go generate ./...
```

### component tests

Component tests can be executed locally with:

```shell script
# build simple-jwt-provider from source code
# setup infrastructure
# run all test file with build-tag component in /cmd/provider 
./component-tests.sh
```