# sqlc-trial

## Description

This is a trial project to test the sqlc and sql-migrate and air tools.

## How to run

1. install task

```bash
go install github.com/go-task/task/v3/cmd/task@latest
echo 'PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zlogin
source ~/.zlogin
```

2. run task to setup

```bash
task setup
```

3. setup database

```bash
task setup-db
```

4. run task to generate sqlc

```bash
sqlc generate
```

5. run server

```bash
air -c .air.toml
```
