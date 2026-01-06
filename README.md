# Agent

Agent is a lightweight, deployable system collector designed to run on any server or workstation.
Its job is to gather hardware information, resource usage, and system-level metrics, and expose them to upstream services such as your Command Center.

The binary is intentionally small, dependency-minimal, and easy to deploy over SSH or via automated provisioning tools.

## Features

* Collects hardware information (CPU, memory, disks, network)
* Monitors real-time resource usage
* Exposes information via gRPC endpoints
* Lightweight footprint, suitable for cloud VPS, on-prem servers, and embedded Linux systems
* Simple build and run workflow

## Cloning

Inside the project we use a git submodule that refers to www.github.com/paradigm-ehb/agent-resources for the C-Library.

```sh
git clone --recursive <agent-repo-url>
# Or if already cloned:
git submodule update --init --recursive
```

## Requirements

Before building the project or generating protobuf files, ensure the following tools are available in your environment:

* Go (compatible with your project version)
* `protoc` — Protocol Buffers compiler
* `protoc-gen-go` — Go protobuf generator
* `protoc-gen-go-grpc` — gRPC code generator for Go

Verify availability using:

```sh
protoc-gen-go --version
protoc-gen-go-grpc --version
```

If they’re not detected in your PATH, install them using your preferred method (Go install, Homebrew, package managers, etc.).

## Building Protobuf Files

The project includes a helper script:

```sh
./build.sh
```

This script regenerates all `.pb.go` and gRPC service definitions.
Use it after editing any `.proto` file under `internal/connection`.

## Running the Agent

The agent can be run directly during development:

```sh
go run cmd/main.go
```

The default configuration starts the gRPC server on port 50051.
You can override this using command-line flags or environment variables depending on your setup.

## Project Structure

A minimal overview of the source layout:

```
cmd/
  main.go          Entry point for the agent  
internal/
  connection/      Protobuf definition and generated gRPC files
  system/          System info collectors (platform-specific)
  handlers/        gRPC service implementations
build.sh           Regenerate protobuf files
README.md          You are here
```

This structure keeps public API code clean while isolating platform-specific logic inside `internal/`.

## Usage Example (Command Line)

Tools such as `grpcurl` allow quick manual interaction with the agent:

```sh
grpcurl -plaintext localhost:50051 agent.Greeter/SayHello -d '{"name":"Appie"}'
```

Useful during debugging or when testing connectivity between Agent and Command Center.


