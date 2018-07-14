# PostgreSQL batch operations in Go

Testing transactions with multiple operations sent in batch.

## Dependencies

Docker and Docker Compose

## Install

* git clone https://github.com/andreiavrammsd/go-postgresql-batch-operations
* go get
* docker-compose up -d

## Start

go run *.go

## Usage

```
curl -X POST \
  http://localhost:8608/users \
  -H 'content-type: application/json' \
  -d '{
	"username": "johndoe",
	"profile":{  
         "firstname":"John",
         "lastname":"Doe"
      },
    "score": 5432
  }'
```

```
curl -X GET \
  http://localhost:8608/users \
  -H 'content-type: application/json'
```

```
curl -X GET \
  http://localhost:8608/actions \
  -H 'content-type: application/json'
```

## Database admin

* http://localhost:54321/
* Use [database credentials](./docker-compose.yml) from environment variables.
