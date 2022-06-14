# Kadisoka IAM Standalone Server

## Configurations

This application loads its configuration from the environment variables. The
documentation or template, including all available fields, for the
configuration can be generated from the application itself. This way, the
documentation or template will always match the implementation.

To generate the template, use this command from the top directory:

```shell
$ go run pkg/iam/iam-server-app/main.go env_file_template
```

The template will be printed out in the stdout. For printing into a file, use:

```shell
$ go run pkg/iam/iam-server-app/main.go env_file_template > config.env
```
