# go-adapter/http/lib

## Overview

This Golang library provides a framework for building HTTP adapters that bridge HTTP requests/responses with application use cases. It uses the `go-adapter/sdk/lib` to automatically map HTTP request parameters to function inputs and function outputs to HTTP responses using struct tags, reducing boilerplate code and coupling between HTTP handling and business logic.

The library provides a clean abstraction layer that allows you to:
- Define HTTP handlers using plain Go functions
- Automatically extract query parameters, headers, and other HTTP data
- Automatically map function return values to HTTP status codes and response bodies
- Handle errors and convert them to appropriate HTTP status codes

## Architecture

The HTTP adapter library follows a layered architecture:

1. **Adapter Interface** (`pkg/adapter.go`): Defines the `Adapter` interface for starting/stopping HTTP servers
2. **Configuration** (`pkg/config.go`): Defines the `Config` structure for adapter setup
3. **Binding System** (`pkg/binding.go`, `pkg/binding_builder.go`): Maps HTTP routes/methods to handlers
4. **Request/Response** (`pkg/request.go`, `pkg/response.go`): Defines HTTP request and response abstractions
5. **Parameter Specifications** (`pkg/param_in_*.go`, `pkg/param_out_*.go`): Handles input/output parameter mapping
6. **Converters** (`pkg/converter.go`): Converts between HTTP types and application types

## Core Concepts

### Adapter

An `Adapter` is the main interface that HTTP adapter implementations must provide:

```go
type Adapter interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

### Configuration

The `Config` structure defines how the HTTP adapter should be configured:

```go
type Config struct {
    Bindings Bindings  // Route-to-handler mappings
    Host     *string   // Optional host
    Port     *int      // Optional port
    Other    any       // Implementation-specific config
}
```

### Bindings

A `Binding` associates HTTP conditions (path, methods) with a handler function:

```go
type Binding struct {
    Condition Condition  // HTTP conditions (path, methods)
    Handler   Handler    // Handler to invoke
}
```

### Request and Response

- **`Request`**: Interface for accessing HTTP request data (query parameters, headers, etc.)
- **`Response`**: Structure for building HTTP responses (status code, body, headers)

## Usage

### Basic Example

```go
package main

import (
    "context"
    httpadpt "github.com/smart-libs/go-adapter/http/lib/pkg"
    "net/http"
)

// Define your handler input structure
type CreateUserInput struct {
    Name  string `query:"name"`
    Email string `query:"email"`
}

// Define your handler output structure
type CreateUserOutput struct {
    UserID int `statuscode:""`
}

// Your handler function
func createUser(input CreateUserInput) (*CreateUserOutput, error) {
    // Business logic here
    return &CreateUserOutput{UserID: 123}, nil
}

func main() {
    ctx := context.Background()
    
    // Create adapter configuration
    config := httpadpt.Config{
        Bindings: []httpadpt.Binding{
            httpadpt.NewBindingBuilderUsingPath("/api/users").
                WithMethods(http.MethodPost).
                WithHandlerFunc(createUser),
        },
        Port: pointers.To(8080),
    }
    
    // Create and start adapter (implementation-specific)
    // adapter := someImpl.NewAdapter(config)
    // adapter.Start(ctx)
}
```

### Query Parameters

Extract query parameters using the `query` tag:

```go
type SearchInput struct {
    Query  string `query:"q"`
    Limit  int    `query:"limit"`
    Offset int    `query:"offset"`
}
```

### Status Codes

Set HTTP status codes using the `statuscode` tag:

```go
type CreateResponse struct {
    StatusCode int `statuscode:""`
    Message    string
}
```

### Error Handling

Errors are automatically converted to appropriate HTTP status codes:

- `IllegalArgumentError` → `400 Bad Request`
- `NotFoundError` → `404 Not Found`
- `DuplicateError` → `409 Conflict`
- `TimeoutError` → `504 Gateway Timeout`
- `IllegalConfigError` → `400 Bad Request`
- Generic errors → `500 Internal Server Error`

## Package Structure

### Core Types

- **`pkg/adapter.go`**: Adapter interface definition
- **`pkg/config.go`**: Configuration structure
- **`pkg/binding.go`**: Binding and Condition types
- **`pkg/handler.go`**: Handler type alias
- **`pkg/request.go`**: Request interface and query parameter handling
- **`pkg/response.go`**: Response structure

### Builder Pattern

- **`pkg/binding_builder.go`**: Fluent builder for creating bindings

```go
// Create a binding with path and methods
binding := httpadpt.NewBindingBuilderUsingPath("/api/users").
    WithMethods(http.MethodGet, http.MethodPost).
    WithHandlerFunc(handlerFunc)

// Or start with methods
binding := httpadpt.NewBindingBuilderUsingMethods(http.MethodGet).
    WithPath("/api/users").
    WithHandlerFunc(handlerFunc)
```

### Parameter Specifications

- **`pkg/param_in_query.go`**: Query parameter extraction
- **`pkg/param_in_spec_factory.go`**: Input parameter spec factory
- **`pkg/param_out_status_code.go`**: Status code output mapping
- **`pkg/param_out_error.go`**: Error output handling
- **`pkg/param_out_spec_factory.go`**: Output parameter spec factory

### Converters

- **`pkg/converter.go`**: Type converters for HTTP-specific conversions
  - Error to HTTP status code conversion
  - String array to single string conversion (for query parameters)

### Utilities

- **`pkg/assertions.go`**: Error handling utilities

## Supported Tags

### Input Tags

- **`query:"name"`**: Extract value from query parameter `name`

### Output Tags

- **`statuscode:""`**: Set HTTP status code from this field
- Error return values are automatically handled and converted to status codes

## Type Conversions

The library includes automatic type conversions:

1. **Query Parameters**: `[]string` (from HTTP) → `string` (to handler)
2. **Errors**: `error` → `int` (HTTP status code)
3. **Standard conversions**: Via the converter library

## Testing

The package includes comprehensive unit tests. See the `pkg/*_test.go` files for examples.

A test suite is provided in `test/test_suite.go` that can be used by HTTP adapter implementations to verify compliance.

## Implementation

This library provides the **interface and framework** for HTTP adapters. Actual HTTP server implementations (e.g., using `net/http`, Gin, Echo, etc.) should be in separate packages like `http/impl/gonethttp`.

To create a new HTTP adapter implementation:

1. Implement the `Adapter` interface
2. Use the `Config` structure for configuration
3. Process `Bindings` to register routes
4. Use the `Handler.Invoke()` method to execute handlers
5. Convert `Request`/`Response` to/from your HTTP framework's types

## Dependencies

- `github.com/smart-libs/go-adapter/sdk/lib`: Core SDK for tag-based handlers
- `github.com/smart-libs/go-adapter/interfaces`: Adapter interfaces
- `github.com/smart-libs/go-crosscutting/assertions/lib`: Assertion utilities
- `github.com/smart-libs/go-crosscutting/converter/lib`: Type conversion
- `github.com/smart-libs/go-crosscutting/serror/lib`: Error handling

## Examples

See `test/test_suite.go` for a complete working example.

## License

[Add your license information here]


