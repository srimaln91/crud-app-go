# CRUD App Go

## How to run

### Use Docker

#### Run using publicly hosted image
```bash
docker run -it --name crud-app-go -p 8086:8086 srimaln91/crud-app-go:latest
```

#### Buid and Run

```bash
cd /path/to/project
docker build --tag crud-app-go:latest .
docker run -it --name crud-app-go -p 8086:8086 crud-app-go:latest
```

The application will be listening on http port 8080.

#### Use custom configuration

You can use the custom config by mounting config.yaml file on the /app/config.yml path in the container.

```bash
docker run -it --name crud-app-go -p 8086:8086 -v /app/config.yaml:</absolute/path/to/config.yaml> crud-app-go:latest
```

Following is a sample config file

```yaml
http:
  port: 8080

logger:
  level: "DEBUG"

database:
  host: localhost
  port: 5432
  name: dbname
  user: username
  password: *****
  pool_size: 5
  max_idle_connections: 2
  conn_max_lifetime: 300000

```

## TODO

01. Write unit tests and automate
02. Expose application matrices in prometheus format
03. improvements on logging
04. Write benchmarks for the hot code path
05. Automate docker builds on release tags