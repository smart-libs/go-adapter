# go-adapter/sdk/lib

## Overview

This Golang library provides an SDK that helps the development of new adapters based on the `go-adapter/interfaces` library.
It provides adapter.UseCaseHandler implementations that are able to invoke any Golang function whose
input and output parameters are tagged to identify how the adapter has to set the input and collect the output.
This reduces the coupling between the adapter and the application functionality.

## Tags available

- default - sets a default value to the input if the parameter is not set.