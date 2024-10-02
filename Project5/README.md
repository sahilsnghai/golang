# Job Listing GraphQL API in Go

A GraphQL API built with Go using `gqlgen` for managing job listings. This API allows you to perform CRUD operations for job listings with a MySQL database as the backend, running in a Docker container.

## Features

- **Create** a new job listing
- **Read** all job listings or a specific job listing by ID
- **Update** job details
- **Delete** a job listing by ID

## Prerequisites

- Go 1.16 or higher installed
- Docker installed for MySQL
- `gqlgen` package installed

## Setting Up MySQL in Docker

1. Pull the MySQL Docker image:
   ```bash
   docker pull mysql:latest
