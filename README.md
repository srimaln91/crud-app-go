# CRUD App Go

This is a sample Go application that showcases CRUD functionalities using the Go programming language. The application utilizes a SQLite database for data storage and retrieval.

## How to run

The application will create an SQLite database automatically. The database path can be configured.

## Build the binary and run it

Build tasks are already configured in the Makefile. You can take a list of tasks by executing `make help` command

```bash
make build
./<build/version/binary> -config=<path/to/config/file.yaml>
```

### Use Docker

#### Run using publicly hosted image
```bash
docker run -dt --name crud-app-go -p 8080:8080 -v </absolute/path/to/config.yaml>:/app/config/config.yaml srimaln91/crud-app-go:latest
```

#### Buid and Run

```bash
cd /path/to/project
make image
docker run -dt --name crud-app-go -p 8080:8080 <image>
```

The application will be listening on http port 8080.

#### Use custom configuration

You can use the custom config by mounting config.yaml file on the /app/config/config.yml mount point in the container.

```bash
docker run -dt --name crud-app-go -p 8080:8080 -v </absolute/path/to/config.yaml>:/app/config/config.yaml <image>
```

Following is a sample config file

```yaml
# Note: All duration values are in milliseconds
http:
  port: 8080

logger:
  level: "DEBUG"

database:
  path: db/tasks.db

```

## Run tests

```bash
make test
```

## Host applications in Kubernetes

Kubernetes config files are placed in the `infra` directory. Use the following command to apply configs

```bash
kubectl apply -R -f infra/
```

Configs contain the following resouce types which are required to run the application
01. Deployment (infra/controllers/deployment.yaml)
02. ConfigMap (infra/config/configmap.yaml)
04. Service (infra/svc/appService.yaml)

## HTTP API

### Create Task
```bash
curl -x POST --location 'http://localhost:8080/api/tasks' \
--header 'Content-Type: application/json' \
--data '{
  "title": "Finish Homework",
  "description": "Complete math problems 1-10 and submit online.",
  "dueDate": "2023-07-05T10:15:00+00:00",
  "completed": false
}'
```

### Get Task
```bash
curl --location 'http://localhost:8080/api/tasks/727fa72b-503a-4799-b315-6dda9f145461'
```

### Get All Tasks
```bash
curl --location 'http://localhost:8080/api/tasks'
```

### Delete Task
```bash
curl -X DELETE --location 'http://localhost:8080/api/tasks/3b7e5edb-6d0c-44a5-94a4-ba47742c0cd3'
```