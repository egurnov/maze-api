API server for storing and solving simple mazes

# Endpoints
GitHub: https://github.com/egurnov/maze-api  
Heroku: https://egurnov-maze-api.herokuapp.com  
Swagger: https://egurnov-maze-api.herokuapp.com/swagger/index.html  
Postman: https://github.com/egurnov/maze-api/blob/main/maze-api.postman_collection.json  

# How to run it

## Prerequisites
* make
* docker
* docker-compose
* go 1.19

### Optional
* mysql

## Local development with Go and MySQL installed
The following make targets can be used:

| Command | Description |
| --- | --- |
| lint | statically analyze the code for common mistakes |
| test | run unit tests for all code |
| e2e-test | run end-to-end tests |
| clean | remove generated binaries |
| build | generate a binary |
| all | all of the above combined |

## Docker
In case you don't have Go and MySQL installed locally, all the same steps can run on Docker:

| Command | Description |
| --- | --- |
| docker-lint | statically analyze the code for common mistakes |
| docker-build-dev | build development image for testing |
| docker-test | run unit ans end-to-end tests in docker |
| docker-build | build the final image |
| docker-all | all of the above combined |

## Docker-compose
This setup allows to properly interact with the system without installing MySQL or running the cod locally.
```
make docker-build up
```

# API Docs
When running locally go to http://localhost:8080/swagger/index.html


# Heroku
To deploy to Heroku:
```
heroku login
heroku container:login
heroku container:push web --app egurnov-maze-api
heroku container:release web --app egurnov-maze-api
```

# Open questions
1. Is this a valid maze? There are multiple open cells in the last row, but only one of them is directly accessible.
```
|_|_|_|_|
|_|_|_|_|
|_|X|X|X|
|_|_|_|_|
```