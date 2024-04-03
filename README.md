# sqlc-trial

## Description

This is a trial project to test the sqlc tool.

1. install task

```bash
go install github.com/go-task/task/v3/cmd/task@latest
echo 'PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zlogin
source ~/.zlogin
```

1. run task to setup

```bash
task setup
```

1. setup database

```bash
task setup-db
```

1. run task to generate sqlc

```bash
sqlc generate
```

1. run server

```bash
air -c .air.toml
```
