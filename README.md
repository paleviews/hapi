# HAPI

HAPI is a restful api design tool utilizing protobuf as IDL
to generate golang code and OpenAPI specifications
that are closely aligned with each other.

## Overview

HAPI takes a different approach to building services by making it possible to
describe the design of the service API using protobuf. HAPI uses the description
to generate specialized service helper code and documentation.

The generated code and documentation are closely aligned with each other,
which takes out most (if not all) error-prone efforts needed to ensure the consistencies
between service implementation and api documentation.

## Install

### 1. Prerequisite

Make sure Go(golang) and protoc are installed in your environment.

HAPI codegen is implemented as a protobuf plugin, so protoc are needed for codegen.

Go to [protocol-compiler-installation](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)
and follow the instructions to install protoc.

### 2. protoc-gen-hapi

```shell
go install github.com/paleviews/hapi/cmd/protoc-gen-hapi@latest
```

## Usage

```shell
protoc --proto_path /path/to/hapi/annotations \
    --proto_path /path/to/your/api/definations \
    --hapi_opt paths=source_relative \
    --hapi_opt doc_file=/path/to/doc/file.yaml \
    --hapi_out /path/to/golang/package/dir \
    foo/foo.proto bar/bar/bar.proto baz.proto
```

## Example

An in-memory todo management service is provided as an example at
[example/todo](example/todo).

## Background

We use a lot of protobuf and gRPC at work. It's great.
We like the design-first approach in particular, and it's strong typed.
The rpc specifications in protobuf and gRPC has become the contract
among all rpc participants involved.

But when gRPC is not an option, which is too often unfortunately,
we are back to the restful world or even worse json-on-http world,
where swagger is the king.
But we still want the strong-typed design-first approach.

So we developed HAPI, a protobuf way of designing restful api.
And start from there, we went out our ways to make HAPI a quick and concise
way of developing restful api.

## FAQ

### 1. Why not just design in OpenAPI and generate code from it?
- Writing api documentation in OpenAPI is a hassle and error-prone.
  We want a more coding-like way when designing.
- There are a lot of check and lint tools in the protobuf ecosystem.

### 2. If you are into protobuf and restful, does grpc-gateway ring a bell?
grpc-gateway is more like a patch. If you have services in gRPC already,
grpc-gateway is a great way to proxy them as restful-ish services. 
But if you are designing restful api from the beginning,
grpc-gateway feels more like a detour.

### 3. What do you mean by "code and documentation are closely aligned with each other"?
For example, if a rpc is annotated as authenticated by bearer token in header,
HAPI will generate:
- code that gets token from header, passes the token to the validating function,
and blocks the request if the validating fails.
- operation in OpenAPI with security scheme of `type: http` and `scheme: bearer`.

So the generated documentation describes exactly what the generated code do,
therefore eliminate the possibility that the documentation says one thing,
and the code says another.
