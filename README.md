# simple-jwt-auth

## Getting started
### Generate ECDSA-512 key pair

```sh
# private key
openssl ecparam -genkey -name secp521r1 -noout -out ecdsa-p521-private.pem
# public key
openssl ec -in ecdsa-p521-private.pem -pubout -out ecdsa-p521-public.pem 
```

## API
### POST `/auth/login`

Expect a JSON object:

```json
{
    "email": "info@leberkleber.io",
    "password": "s3cr3t"
}
```

This endpoint will check the email/password combination and will set the respond with an jwtauthToken if correct:

```json
{
    "access_token":"<jwt>"
}
```
