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

This will start the server on port **8054**. You can change the port by setting the `TCK_PORT` environment variable or by adding a .env file with the same variable.

Once started, your TCK server will be up and running! 🚦

# Start all TCK tests with Docker 🐳

This section covers setting up and running TCK tests using Docker.

## Prerequisites

Before you begin, ensure you have the following installed:

-   **Go**: Version 1.20 or higher
-   **Docker**: Latest version
-   **Docker Compose**: Latest version
-   **Task**: Latest version

## 🔹 Run a specific test

```bash
task run-specific-test TEST=AccountCreate
```

This will:

-   Verifies prerequisites

-   Starts the TCK server

-   Launches required containers

-   Run only the `AccountCreate` tests

## 🔹 Run all tests

To run all tests:

```bash
task start-all-tests
```

This will:

-   Verifies prerequisites

-   Starts the TCK server

-   Launches required containers

-   Run all tests automatically

Sit back and let Docker do the work! 🚀

### ⚙️ Running Tests Against Hiero Testnet

To run tests against the Hiero Testnet, use the following command:

```bash
task run-specific-test \
  NETWORK=testnet \
  OPERATOR_ACCOUNT_ID=your-account-id \
  OPERATOR_ACCOUNT_PRIVATE_KEY=your-private-key \
  MIRROR_NODE_REST_URL=https://testnet.mirrornode.hedera.com \
  MIRROR_NODE_REST_JAVA_URL=https://testnet.mirrornode.hedera.com \
  # Run specific test
  TEST=AccountCreate
```

### 🎉 All Done!

Your TCK server is now running inside Docker! 🚀 You can now execute tests and validate the system.

Need help? Reach out to the team! 💬👨‍💻

Happy coding! 💻✨
