# go-adapter/sdk/lib

## Overview

This Golang library provides an SDK that helps the development of new adapters based on the `go-adapter/interfaces` library. It provides `adapter.UseCaseHandler` implementations that are able to invoke any Golang function whose input and output parameters are tagged to identify how the adapter has to set the input and collect the output. This reduces the coupling between the adapter and the application functionality.

The SDK uses reflection and struct tags to automatically map adapter inputs to function parameters and function return values to adapter outputs, eliminating the need for manual mapping code.

## Architecture

The SDK follows a layered architecture:

1. **Handler Layer** (`pkg/handler/`): Provides high-level handler interfaces and builders
2. **Tag-based Handler** (`pkg/handler/tagbased/`): Implements tag-based function invocation using reflection
3. **Use Case Handler** (`pkg/handler/usecase/`): Bridges the SDK handlers with the adapter interfaces
4. **Parameter Specification** (`pkg/param/`): Defines input/output parameter specifications and options
5. **Tag-based Parameter Factory** (`pkg/param/tagbased/`): Creates parameter specs from struct tags

## Core Concepts

### Handler

A `Handler[Input, Output]` is the main interface that adapters use to invoke application use cases. It takes adapter-specific Input and Output types and handles the conversion to/from the application's use case function.

```go
type Handler[Input any, Output any] interface {
    Invoke(ctx context.Context, input Input, output Output) error
}
```

### Input/Output Specifications

- **InputSpecs**: Maps `adapter.ParamRef` to `InputParamSpec[Input]` instances that know how to extract values from the adapter's Input type
- **OutputSpecs**: Maps `adapter.ParamRef` to `OutputParamSpec[Output]` instances that know how to set values into the adapter's Output type

### Tag-based Configuration

The SDK uses struct tags to automatically configure parameter mappings:

- **Input tags**: Define how to extract values from the adapter Input (e.g., from flags, environment variables, positional args)
- **Output tags**: Define how to set values into the adapter Output (e.g., to response fields, error fields)

## Package Structure

### `pkg/handler/`

Core handler interfaces and builder patterns:

- **`handler.go`**: Defines the `Handler[Input, Output]` interface
- **`builder.go`**: Provides builder interfaces for constructing handlers with input/output specs

### `pkg/handler/tagbased/`

Tag-based handler implementation using reflection:

- **`handler.go`**: Core tag-based handler that invokes functions using reflection
- **`handler_builder.go`**: Builder for creating tag-based handlers from functions
- **`arg_factory.go`**: Creates argument factories for function parameters
- **`arg_factory_structure.go`**: Handles nested struct parameters recursively
- **`arg_factory_func_inputs.go`**: Processes function input arguments
- **`output_builder_action.go`**: Actions for building output from function return values
- **`output_builder_action_factory.go`**: Factory for creating output builder actions

### `pkg/handler/usecase/`

Use case handler implementation:

- **`handler.go`**: Implements `adapter.UseCaseHandler` using SDK handlers
- **`input_accessor.go`**: Adapter input accessor implementation
- **`output_builder.go`**: Adapter output builder implementation

### `pkg/param/`

Parameter specification and options:

- **`spec.go`**: Base `Spec` interface with name and options
- **`spec_in.go`**: `InputParamSpec[Input]` interface and implementation
- **`spec_out.go`**: `OutputParamSpec[Output]` interface and implementation
- **`input_specs.go`**: Collection of input parameter specifications
- **`output_specs.go`**: Collection of output parameter specifications
- **`option.go`**: Option function type for parameter processing
- **`option_default.go`**: Default value option
- **`option_mandatory.go`**: Mandatory parameter validation
- **`option_not_default_value.go`**: Validation against default values
- **`option_not_blank_or_empty_string.go`**: String validation options
- **`option_converter.go`**: Type conversion options
- **`option_target_type.go`**: Target type specification

### `pkg/param/tagbased/`

Tag-based parameter specification factories:

- **`spec_in_factory.go`**: Factory interface for creating input specs from struct fields
- **`spec_out_factory.go`**: Factory interface for creating output specs from struct fields
- **`input_specs_builder.go`**: Builder for input specifications from struct tags
- **`output_specs_builder.go`**: Builder for output specifications from struct tags
- **`option_factory.go`**: Factory for creating options from struct tags
- **`option_factory_default.go`**: Default value option factory from `default` tag
- **`option_factory_mandatory.go`**: Mandatory/assertion options from `assert` tag
- **`spec_in_factory_impl.go`**: Registry-based input spec factory implementation
- **`spec_out_factory_impl.go`**: Registry-based output spec factory implementation
- **`spec_in_factory_registry.go`**: Registry interface for input spec factories
- **`spec_out_factory_registry.go`**: Registry interface for output spec factories
- **`tag_name.go`**: Tag name type definitions

### `pkg/async/`

Asynchronous task management (currently commented out/in development):

- **`manager.go`**: Task manager for managing concurrent tasks
- **`handler.go`**: Task handler implementation
- **`handler_list.go`**: Collection of task handlers

### `pkg/`

Utility packages:

- **`debug.go`**: Debug utilities for development
- **`logger.go`**: Logging utilities

## Usage

### Basic Tag-based Handler

The simplest way to use the SDK is with tag-based handlers:

```go
import (
    "context"
    sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
    tagbasedhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler/tagbased"
    tagbased "github.com/smart-libs/go-adapter/sdk/lib/pkg/param/tagbased"
)

// Define your handler function
func MyHandler(ctx context.Context, input MyInput) (MyOutput, error) {
    // Your business logic here
    return MyOutput{Result: input.Value}, nil
}

// Create input and output factories (provided by specific adapters)
inputFactory := tagbased.NewInputParamSpecFactory[MyAdapterInput](inputRegistry)
outputFactory := tagbased.NewOutputParamSpecFactory[MyAdapterOutput](outputRegistry)

// Build the handler
handler := tagbasedhandler.NewBuilderForFunc[MyAdapterInput, MyAdapterOutput](MyHandler).
    WithInTagBasedFactory(inputFactory).
    WithOutTagBasedFactory(outputFactory).
    Build()

// Use the handler
err := handler.Invoke(ctx, adapterInput, adapterOutput)
```

### Manual Specification Builder

For more control, you can manually specify input/output mappings:

```go
import (
    sdkhandler "github.com/smart-libs/go-adapter/sdk/lib/pkg/handler"
    sdkparam "github.com/smart-libs/go-adapter/sdk/lib/pkg/param"
    "github.com/smart-libs/go-adapter/interfaces/pkg/adapter"
)

// Create input spec
inputSpec := sdkparam.NewInputParamSpec[MyAdapterInput](
    "paramName",
    []sdkparam.Option{
        sdkparam.Mandatory(),
        sdkparam.Default("defaultValue"),
    },
    converters,
    func(input MyAdapterInput) (any, error) {
        return input.Field, nil
    },
)

// Create output spec
outputSpec := sdkparam.NewOutputParamSpec[MyAdapterOutput](
    sdkparam.NewSpec("result"),
    func(output MyAdapterOutput, value any) error {
        output.Result = value.(string)
        return nil
    },
)

// Build handler (using builder pattern)
handler := sdkhandler.NewBuilder[MyAdapterInput, MyAdapterOutput]().
    WithInParamSpec(adapter.StringParamRef("paramName"), inputSpec).
    Output().
    WithOutParamSpec(adapter.StringParamRef("result"), outputSpec).
    Build()
```

## Tag-based Configuration

### Input Tags

Input tags are processed by `InputParamSpecFactory` implementations. The SDK provides a registry-based factory that can be extended with custom factories.

**Available Tags:**

- **`default`**: Sets a default value if the parameter is not provided
  ```go
  type MyInput struct {
      Value string `default:"defaultValue"`
  }
  ```

- **`assert`**: Adds validation assertions (comma-separated)
  ```go
  type MyInput struct {
      Name string `assert:"mandatory,notBlank"`
      ID   int    `assert:"mandatory,notDefaultValue"`
  }
  ```

  Supported assertions:
  - `mandatory`: Parameter must be provided
  - `notDefaultValue`: Value cannot be the zero/default value for its type
  - `notBlank`: String cannot be blank (whitespace only)
  - `notEmpty`: String cannot be empty

### Output Tags

Output tags are processed by `OutputParamSpecFactory` implementations. The factory determines how function return values map to adapter output fields.

## Parameter Options

Options are functions that process parameter values. They can be chained together:

### Default Value

```go
sdkparam.Default(defaultValue)
sdkparam.Default(defaultValue, customNilCheck)
```

Sets a default value if the parameter is nil or not provided.

### Mandatory

```go
sdkparam.Mandatory()
sdkparam.Mandatory(customNilCheck)
```

Validates that the parameter is provided and not nil.

### Type Conversion

```go
sdkparam.ConverterTyped[From, To](converterFunc)
sdkparam.ConvertTo[TargetType](converters)
```

Converts parameter values to target types.

### String Validation

```go
sdkparam.NotBlankString()        // String cannot be blank
sdkparam.NotEmptyString()        // String cannot be empty
sdkparam.NotBlankOrEmptyString() // Both checks
```

### Value Validation

```go
sdkparam.NotDefaultValue[T]()           // Type-specific default check
sdkparam.NotDefaultValueReflection()    // Reflection-based default check
```

### Conditional Options

```go
sdkparam.IfNotNil(options...) // Apply options only if value is not nil
```

## Input/Output Specifications

### InputParamSpec

An `InputParamSpec[Input]` knows how to:
- Extract a value from the adapter Input type (`GetValue`)
- Copy a value from the adapter Input to a target (`CopyValue`)

### OutputParamSpec

An `OutputParamSpec[Output]` knows how to:
- Set a value into the adapter Output type (`SetValue`)

## Factory Pattern

The SDK uses factory patterns for extensibility:

### InputParamSpecFactory

Creates `InputParamSpec[Input]` instances from struct fields. Implementations can:
- Read struct tags
- Analyze field types
- Create appropriate parameter specifications

### OutputParamSpecFactory

Creates `OutputParamSpec[Output]` instances from struct fields. Implementations can:
- Read struct tags
- Analyze field types
- Create appropriate output mappings

### Registry Pattern

Factories can be registered in a registry to support multiple tag-based strategies:

```go
type InputParamSpecFactoryRegistry[Input any] interface {
    AsList() []InputParamSpecFactory[Input]
    // Add factories to the registry
}
```

## Function Invocation Flow

1. **Handler.Invoke** is called with adapter Input and Output
2. **InputAccessor** is created from Input and InputSpecs
3. **OutputBuilder** is created from Output and OutputSpecs
4. **Function arguments** are built using arg factories:
   - Context is passed through if first parameter
   - Struct arguments are created and populated from InputAccessor
   - Nested structs are handled recursively
5. **Function is invoked** using reflection (`reflect.Value.Call`)
6. **Return values** are processed using output builder actions
7. **Output is built** and returned

## Nested Structures

The SDK supports nested structures recursively:

```go
type NestedInput struct {
    User struct {
        Name string `assert:"mandatory"`
        Age  int    `default:"0"`
    }
}
```

The SDK will automatically:
- Create nested struct instances
- Map fields recursively
- Handle pointer vs value types

## Error Handling

- **Parameter extraction errors**: Returned as errors from `GetValue`/`CopyValue`
- **Validation errors**: Returned when options fail (e.g., mandatory check)
- **Function invocation panics**: Should be handled by the adapter
- **Output building errors**: Panic in `SetValue` (can be customized)

## Debugging

Enable debug output:

```go
import "github.com/smart-libs/go-adapter/sdk/lib/pkg"

sdk.DebugEnabled = true
sdk.Debug = customDebugFunction // Optional: customize debug output
```

## Dependencies

- `github.com/smart-libs/go-adapter/interfaces`: Adapter interface definitions
- `github.com/smart-libs/go-crosscutting/converter/lib`: Type conversion utilities

## Examples

See the test files in `cli/lib/test/` for usage examples.

## Extension Points

The SDK is designed for extension:

1. **Custom Option Factories**: Add to `tagbased.OptionFactories` to support new tags
2. **Custom Input/Output Factories**: Implement factory interfaces for custom tag processing
3. **Custom Converters**: Provide converter implementations for type conversions
4. **Custom Options**: Create new `Option` functions for parameter processing

## Thread Safety

- Handlers are stateless and can be used concurrently
- InputSpecs and OutputSpecs should be immutable after creation
- Factories should be thread-safe if shared

## Performance Considerations

- Reflection is used for function invocation - consider caching for hot paths
- Parameter specifications are created once during handler construction
- Options are chained and executed in sequence - minimize option count for performance
