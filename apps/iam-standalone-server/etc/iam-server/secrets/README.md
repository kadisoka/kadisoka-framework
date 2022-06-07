This folder contains secrets for the IAM services. Do not commit anything
contained in here.

IAM service will attempt to load these secrets:

- `clients.csv` -- as we haven't implemented client management, we use this
  file to store registered clients. As the file contains clients' secret,
  it needs to be placed here. Create yours by copying `clients.csv.example`
  as `clients.csv` and start adding rows.
- `jwt.key` or `jwt_ed25519` -- a private key for use to sign JWT tokens.
  Use the command
  `openssl genpkey -algorithm ed25519 -outform PEM -out jwt_ed25519.key` to
  generate yours.
- `config.env` -- configuration for the service. For local execution, it's
  referenced by `docker-compose.yaml`.
