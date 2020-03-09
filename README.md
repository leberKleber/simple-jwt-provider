# simple-jwt-provider

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
