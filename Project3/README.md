# Book CRUD API in Go

A simple CRUD API built with Go for managing books. The API connects to a MySQL database running in Docker and allows you to create, read, update, and delete book records.

## Features

- **Create** a new book
- **Read** all books or a specific book by ID
- **Update** book details
- **Delete** a book by ID

## Prerequisites

- Go 1.16 or higher installed
- Docker installed for MySQL

## Setting Up MySQL in Docker

1. Pull the MySQL Docker image:
   ```bash
   docker pull mysql:latest