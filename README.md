# simple-jwt-auth

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
