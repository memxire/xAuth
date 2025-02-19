# SSO (Single Sign-On) service xAuth

SSO (Single Sign-On) service with gRPC API called xAuth.

## Features

- Register user
- Login user
- Check if user is admin

## How to use

1. Clone repository
2. Run `go run cmd/sso/main.go --config=./config/local.yaml` in project root
3. Use gRPC client to interact with SSO service