## API specification

This micro service provides two endpoints that accept and provide JSON.

_Clients are responsible for storing the id and AES key in order to get the secret later._

### Create a secret

#### Request

Route: `POST /v1/secrets`

Headers:

- `Content-Type: application/json`

Body:

| Attribute | Type   | Description               |
|:----------|:-------|:--------------------------|
| id        | string | identifier for the secret |
| data      | string | plain text to encrypt     |

Example request:
```bash
curl \
  --cacert acceptance_test/testdata/cert.pem \
  https://127.0.0.1:8080/v1/secrets \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"id":"some-id","data":"some plain text"}'
```

#### Response

Status code: `200 OK`

Headers:

- `Content-Type: application/json`

Body:

| Attribute | Type   | Description                                      |
|:----------|:-------|:-------------------------------------------------|
| key       | string | hexadecimal encoded AES key used to encrypt data |

Example response body:
```json
{
  "key":"1bc50ee2992feba6c1d9e384b3c8e9203dcfc0eed50c032dfc2821ca2aa0cfa5",
}
```

### Get a secret

#### Request

Route: `GET /v1/secrets/{id}`

Headers:

- `Content-Type: application/json`

Body:

| Attribute | Type   | Description                                       |
|:----------|:-------|:--------------------------------------------------|
| key       | string | hexadecimal encoded AES key to decrypt the secret |

Example request:
```bash
curl \
  --cacert acceptance_test/testdata/cert.pem \
  https://127.0.0.1:8080/v1/secrets/some-id \
  -X GET \
  -H 'Content-Type: application/json' \
  -d '{"key":"1bc50ee2992feba6c1d9e384b3c8e9203dcfc0eed50c032dfc2821ca2aa0cfa5"}'
```

### Response

Status code: `200 OK`

Headers:

- `Content-Type: application/json`

Body:

| Attribute | Type   | Description          |
|:----------|:-------|:---------------------|
| data      | string | decrypted plain text |

Example response body:
```json
{
  "data": "some plain text"
}
```

### Error responses

Status code: `4xx–5xx`

Headers:

- `Content-Type: application/json`

Body:

| Attribute | Type   | Description   |
|:----------|:-------|:--------------|
| error     | string | error message |

Example response body:
```json
{
  "error": "unsupported Content-Type"
}
```
