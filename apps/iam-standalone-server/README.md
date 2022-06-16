# Kadisoka IAM Standalone Server

## Server Configurations and Credentials

The server application will look at `./etc/iam-server/secrets` for
the configuration and credential files.

### General configuration

The file is `./etc/iam-server/secrets/config.env`. This file will be loaded
by `docker-compose` instead of directly by the application.

The application loads its configuration from the environment variables. The
documentation or template, including all available fields, for the
configuration can be generated from the application itself. This way, the
documentation or template will always match the implementation.

To generate the template, use this command from the top directory:

```shell
$ go run pkg/iam/iam-server-app/main.go env-file-template > apps/iam-standalone-server/etc/iam-server/secrets/config.env
```

You would need to fill in some fields before the application could be started
properly.

To make it up and running **locally**, use these values. In production, these
fields must be set up properly:

```
IAM_DB_URL=postgres://iwm:hdig8g4g49htuhe@iwm-db/iwm?sslmode=disable
IAM_EAV_EMAIL_DELIVERY_SERVICE=null
IAM_MEDIA_STORE_SERVICE=local
IAM_MEDIA_LOCAL_SERVER_SERVE_PATH=/iam
IAM_MEDIA_LOCAL_SERVER_SERVE_PORT=10080
IAM_PNV_SMS_DELIVERY_SERVICE=null
```

Note that above configuration makes the application not sending the OTPs.
You'd need to access the database to get the OTPs.

For the field `IAM_DB_URL`, the value was constructed from `iam-db` credentials
found in the `docker-compose.yaml` file.

Note. To reduce noise in the log, add this field into the end of `config.env`:

```
LOG_TRIM_PKG_PREFIXES=github.com/kadisoka/kadisoka-framework/
```

### JWT Signer Key

The file is `./etc/iam-server/secrets/jwt.key` or
`./etc/iam-server/secrets/jwt_ed25519.key`.

To generate, use the command:

```shell
$ openssl genpkey -algorithm ed25519 -outform PEM -out apps/iam-standalone-server/etc/iam-server/secrets/jwt_ed25519.key
```

JWT signer key is used to sign all JWT tokens issued by the server.

### Client Applications Registry Table

The file is `./etc/iam-server/secrets/clients.csv`.

As we haven't implemented client management, we use this file to store
registered client application credentials.

Create yours by copying `clients.csv.example` as `clients.csv` and start
adding rows by running from the top directory:

```shell
$ go run pkg/iam/tools/app-id-gen/main.go
```

It will generate the ID and its secret. Insert them into the CSV file.

## Running the Application Locally

NOTE: the application requires some configuration and credentials as described
in the previous section.

To start the application, run this command from the top directory:

```shell
$ docker-compose -f apps/iam-standalone-server/docker-compose.yaml up --build
```
