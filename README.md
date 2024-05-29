# gRPC User Service

This repository contains a simple gRPC-based user service implemented in Go. The service allows for creating, retrieving, and searching users. This document provides an overview of the project, the design patterns used, and best practices followed.

## Overview

The user service exposes four main operations:
1. **CreateUser**: Add a new user.
2. **GetUser**: Retrieve a user by their ID.
3. **GetUsers**: Retrieve multiple users by their IDs.
4. **SearchUser**: Search for users based on various fields.

## Project Structure

- **main**: The entry point of the application where the gRPC server is initialized and run.
- **proto**: Contains the Protocol Buffers (`.proto`) file and the generated Go code for gRPC.
- **models**: Defines the data structures used in the service.
- **repository**: Implements data access logic using the repository pattern.
- **service**: Contains the business logic and implements the gRPC service methods.

## Design Patterns Used

1. **Repository Pattern**:
   - Abstracts data access logic.
   - The `UserRepository` interface defines methods for CRUD operations and searching users.
   - `InMemoryUserRepository` is a concrete implementation that manages user data in memory.

2. **Service Layer Pattern**:
   - Encapsulates business logic.
   - `UserService` struct provides methods for user-related operations and interacts with the repository.

3. **Factory Pattern**:
   - Used to create instances of the service and repository.
   - `NewUserService` creates and returns a new instance of `UserService`.
   - `NewUserRepository` creates and returns a new instance of `InMemoryUserRepository`.

## Best Practices Followed

1. **Dependency Injection**:
   - The service layer is initialized with a repository dependency, making the code modular and testable.

2. **Concurrency Safety**:
   - `sync.Mutex` is used in `UserRepository` to ensure thread safety when accessing shared data.

3. **Separation of Concerns**:
   - The project is organized into distinct packages for different concerns (data access, business logic, etc.).

4. **Use of Interfaces**:
   - The `UserRepository` interface defines the contract for data access, allowing flexibility and easier testing.

5. **Error Handling**:
   - Consistent error checking and handling throughout the code to ensure robustness.


## Running the Project

### Prerequisites

- Go 1.16 or higher
- Protocol Buffers compiler (`protoc`)
- gRPC Go plugin for Protocol Buffers

```sh
git clone https://github.com/vkmaurya09/grpc-user-service.git
cd grpc-user-service
go mod tidy
protoc --go_out=. --go-grpc_out=. proto/user.proto
go run main.go
```

### Testing

To run the unit tests, use the following command:

```sh
cd test
go test ./...
```

### Docker

```sh

 - Build
    docker build -t user-service:latest .

 - Run
    docker run -p 50051:50051 user-service

