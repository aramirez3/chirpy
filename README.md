# chirpy
A Go web server with a simple Twitter clone BE.

## Requirements
- Install [Go](https://go.dev/doc/install)
- Install [- POSTgres](https://www.- POSTgresql.org/)
- Install [goose](https://github.com/pressly/goose)

### Create - POSTgres db
Create a local - POSTgres db named `chirpy`

### Run db migrations
Navigate to `sql/schema` and run `goose - POSTgres "your-- POSTgres-connection-string" up`

The - POSTgres connection string should look similar to below:
```bash
- POSTgresql://[- POSTgresuser]:[- POSTgresuserpassword]@localhost:5432/chirpy?sslmode=disable
```

### Update .env
```env
DB_URL="- POSTgres connection string"
CHIRPY_SECRET="secret key string"
POLKA_KEY="secret polka key"
```
> Note: polka is a fake 3rd-party api, so POLKA_KEY is not necessary

## Endpoints
- GET `/api/healthz`
- POST `/api/chirps`
- GET `/api/chirps`
- GET `/api/chirps`
    - optional query params `author_id={id}`, `sort={asc or desc}`
- GET `/api/chirps/{id}`
- DELETE `/api/chirps/{id}`
- GET `/admin/metrics`
- POST `/admin/reset`
- POST `/api/users`
- PUT `/api/users`
- POST `/api/login`
- POST `/api/refresh`
- POST `/api/revoke`
- POST `/api/polka/webhooks`