# SSO (Single Sign-On) service xAuth

SSO (Single Sign-On) service with gRPC API called xAuth.

## Features

- Register user
- Login user
- Check if user is admin

## How to use

1. Clone repository

2. Run `go run cmd/sso/main.go --config=./config/local.yaml` or `task server`
   in project root

3. Use gRPC client to interact with SSO service

## Database Migrations

To complete database migrations, run the following command in the project root:

```bash
task migrate
```

File `sso.db` will be created in the `xauth/storage` directory.

## Project structure

```bash
.xauth
├── cmd.............. Commands for running application and utilities
│   ├── migrator.... Database Migration Utility
│   └── sso......... Main entry point to the SSO service
├── config........... Configuration yaml files
├── internal......... Project insides
│   ├── app.......... Code to launch various components of the application
│   │   └── grpc.... Starting gRPC server
│   ├── config....... Loading configuration
│   ├── domain
│   │   └── models.. Data structures and domain models
│   ├── grpc
│   │   └── auth.... gRPC handlers of the Auth service
│   ├── lib.......... General helper utilities and functions
│   ├── services..... Service layer (business logic)
│   │   ├── auth
│   │   └── permissions
│   └── storage...... Data processing layer
│       └── sqlite.. Implementation in SQLite
├── migrations....... Migrations for the database
├── storage.......... Storage files, such as SQLite databases
└── tests............ Functional tests
```
