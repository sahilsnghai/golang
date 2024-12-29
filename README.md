# Go Projects Overview

This repository contains five Go-based projects:

- **Project1**: Renders a static HTML page.
- **Project2**: CRUD API for movies.
- **Project3**: CRUD API for a bookstore with MySQL.
- **Project4**: CRUD API for people with MongoDB.
- **Project5**: GraphQL CRUD API for job listings.
- **Project6**: SuperMarket Store Backend.


## Project1
```bash
cd Project1
go run main.go
```


## Project2
```bash
Copy code
cd Project2
go run main.go
```
Endpoints: GET /movies, POST /movies, PUT /movies/{id}, DELETE /movies/{id}


## Project3
```bash
Copy code
docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=bookdb -p 3306:3306 -d mysql
cd Project3
go run main.go
```
Endpoints: GET /books, POST /books, PUT /books/{id}, DELETE /books/{id}


## Project4
```bash
Copy code
docker run --name mongo -p 27017:27017 -d mongo
cd Project4
go run main.go
```
Endpoints: GET /people, POST /people, PUT /people/{id}, DELETE /people/{id}


## Project5
```bash
Copy code
cd Project5
go run main.go
```
GraphQL Playground: http://localhost:8080


## Project6
```bash
Copy code
cd Project6
go run main.go
```
SuperMarket Backend: http://localhost:8080/
POST /register, POST /login, GET /getproduct, POST cart/checkout ...


## Project7
```bash
Copy code
cd Project7
go run cmd/student-api/main.go
```
Student-API Backend: http://localhost:8080/
POST /api/students, GET /api/students/{id}, GET /api/students


## Project10
```bash
Copy code
cd Project10
go run . -commands
```
go run .|-add "Item" | -edit "id:'Updated Item'" |-toggle id | -del id
