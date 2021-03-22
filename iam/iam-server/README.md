# IAM Server Application

To ensure that env-file template is in-sync with the executable and its
built-in features, executable was made to be able to generate the template.
This will ensure that the template and, optionally, default values are in-sync
with the executable.

To generate env-file template, use the command (in the root directory):

```shell
$ go run iam/iam-server/main.go env_file_template
```

We can also use prebuilt executable (in the directory where the executable
resides):

```shell
$ ./iam-server env_file_template
```

To obtain the template as a file, use something like:

```shell
$ ./iam-server env_file_template > config.env
```
