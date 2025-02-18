# Go SDK TCK server

This is a server that implements the [SDK TCK specification](https://github.com/hiero-ledger/hiero-sdk-tck/) for the Go SDK.

# Server start up guide for Go 🛠️

This guide will help you set up, start, and test the TCK server using Docker and Go. Follow the steps below to ensure a smooth setup.

## 🚀 Start the TCK Server

Run the following commands to build and start the server:

```bash
# From the tck directory
go mod tidy
go run cmd/server.go
```

This will start the server on port **80**. You can change the port by setting the `TCK_PORT` environment variable or by adding a .env file with the same variable.

Once started, your TCK server will be up and running! 🚦

# Start all TCK tests with Docker 🐳

This guide will help you set up and start the TCK server, local node and run all TCK tests using Docker. Follow these steps to ensure all dependencies are installed and the server runs smoothly.

## Prerequisites

Before you begin, ensure you have the following installed:

-   **Go**: Version 1.20 or higher
-   **Docker**: Latest version
-   **Docker Compose**: Latest version
-   **Task**: Latest version

### Installing Task

You can install Task using one of these methods:

```bash
# Using Homebrew (macOS)
brew install go-task

# Using Go
go install github.com/go-task/task/v3/cmd/task@latest
```

## 🔧 Setup Instructions

### 1. Check Go Version

Verify that Go is installed and meets the version requirements:

```bash
go version
```

### 2. Install Hedera Local Node CLI

If not already installed, run the following command:

```bash
npm install @hashgraph/hedera-local -g
```

### 3. Start the Local Hedera Network

Run the following command to start the local Hedera network:

```bash
task start-local-node
```

### 4. Build the Docker Image

Build the Docker image for the TCK Go server:

```bash
task build-tck-go-server
```

### 5. Run a specific test

```bash
task run-specific-test -- TEST=AccountCreate
```

This will:

-   Spin up the TCK server
-   Start required containers
-   Run only the **AccountCreate** tests

### 6. Start All Services

Now, let's fire up all the services using Docker Compose:

```bash
task start-all-tests
```

This will:

-   Spin up the TCK server
-   Start required containers
-   Run all tests automatically

Sit back and let Docker do the magic!

### 🎉 All Done!

Your Go TCK server is now running inside Docker! 🚀 You can now execute tests and validate the system.
