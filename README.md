# Codebase

## Design Architecture

This project uses hexagonal architecture. Hexagonal architecture is another form of applying the clean architecture principle. 
Basically, hexagonal architecture uses the basic principles of clean architecture. For further explanation 
regarding this, you can go to the following article [hexagonal](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/)

In this project there are 3 parts, namely the core, port and adapter. The port is a collection of interfaces that 
will later connect the core part of the code (business logic) with the implementation of the port, namely the adapter.
The adapter itself usually consists of an external implementation of the code that does not affect the code at all,
examples include repository, apicall and redis.

![img.png](img.png)

- **app** package contains cmd and deps. cmd contains commands to run the server and scheduler while deps contains initialize the required dependencies.
- **config** package contains initialization configuration
- **core** package contains ports and use cases
- **docs** package contains api documentation via swagger
- **infrastructure** package contains implementations of ports called adapters
- **interface**  package contains the initialization of the API route
- **pkg** package contains helpers and utilities

## Getting Started

- The system must be ensured to have MongoDB and Redis
- Install dependencies ```go mod download``` or ```go mod tidy```.
- Run the server, open terminal and run ```go run main.go```.
- Run the scheduler, open new terminal and run ```go run main.go cron```.
- Make sure the logo appear.
- Server run at ```localhost:8084```.
- Visit swagger at ```localhost:8084/api/swagger/index.html#/```.
- **Register user to access all endpoints**.
- Run unit test ```go test -v ./...```


## Stacks

- Golang 1.21
- MongoDB
- JWT
- Redis
