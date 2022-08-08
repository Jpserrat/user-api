# USER API

This is a simple user CRUD

## Run the app

    docker-compose up
    go run main.go

The api will be running on the port 8082

# REST API

## Get list of users

### Request

`GET /v1/users?limit=10&page=1`

### Response

    {
	    "body": [
		    {
			    "id": "62efb852a6f111e1ad00c90b",
			    "name": "test",
			    "email": "test@test.com",
			    "age": 20
		    }
	    ],
	    "status": 200
    }

## Create a new user

### Request

`POST /v1/auth/register`

### Response

    201 Created

## Create a new user

### Request

`POST /v1/auth/login`

### Response

    {
	    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjE2NTk5Njg3NDZ9.K0qDF325ADofnn4T8VMSRJfBAav4nZMxVOluj5IGeNE"
    }

## Get a specific user

### Request

`GET /v1/users/id`

### Response
    {
	    "id": "62efb852a6f111e1ad00c90b",
	    "name": "test",
	    "email": "test@test.com",
	    "age": 24
    }

## Update your user

### Request

`PUT /v1/users/id`


### Response

   204 No Content

## Delete your user

### Request

`DELETE /v1/users/id`

### Response

    204 No Content
