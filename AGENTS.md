# Agent Instructions for A2A Golang Project

## 1. Project Overview
This project implements the **A2A (Agent-to-Agent) Protocol** demonstration in Golang.
It showcases two agents, **Agent A (Assistant)** and **Agent B (Finance)**, collaborating via JSON-RPC 2.0 and Server-Sent Events (SSE).

**Key Features:**
- **Context-Aware**: Agents maintain conversational state.
- **Standardized Interface**: Language-agnostic A2A protocol.
- **Real-Time Visibility**: SSE for streaming long-running task progress.
- **Dependency-Free**: Pure Go standard library implementation (preferred).

## 2. Environment & Tooling
**System Requirements:**
- **OS**: Linux (NixOS preferred).
- **Go**: Latest stable version (managed via Nix).
- **Build System**: Standard `go` toolchain.
- **Environment**: Use `flake.nix` (if available) or `nix-shell -p go golangci-lint`.

**Nix Integration:**
- Always prefer Nix-based environments.
- Do not suggest `apt-get` or global `pip` installs.
- If `flake.nix` is missing, suggest creating one using the reference below.

## 3. Build, Run, and Test Commands

### Setup
Ensure the Go module is initialized. If `go.mod` is missing:
```bash
go mod init a2a
go mod tidy
```

### Running Agents
This project follows the `cmd/<app_name>/main.go` pattern.

**Agent B (Finance Server):**
```bash
# Must be started first to listen on port 8080
go run cmd/agent_b/main.go
# Expected output: 🚀 Agent B (財務專員) 已啟動...
```

**Agent A (Assistant Client):**
```bash
# Run in a separate terminal to interact with Agent B
go run cmd/agent_a/main.go
# Expected output: Interaction logs and SSE stream
```

### Testing
Run all tests in the repository:
```bash
go test ./... -v
```

**Running a Single Test:**
To run a specific test function (e.g., `TestProtocolHandshake`) in a specific package:
```bash
go test -v -run TestProtocolHandshake ./path/to/package
```

### Linting
Use `golangci-lint` for code quality checks:
```bash
golangci-lint run ./...
```
*Note: Ensure `golangci-lint` is installed via Nix.*

## 4. Code Style & Conventions

### General Principles
- **Idiomatic Go**: Follow `Effective Go` and standard Go proverbs.
- **Zero Dependencies**: Prefer `stdlib` over external packages unless absolutely necessary.
- **Simplicity**: Keep logic simple and readable (KISS).
- **Type Safety**: Use strong typing and interfaces.

### Formatting & Imports
- **Format**: Always run `gofmt -s -w .` before committing.
- **Imports**: Group imports into three sections separated by blank lines:
  1. Standard library (`"fmt"`, `"net/http"`)
  2. Third-party packages (if any)
  3. Local packages (`"a2a/models"`)

### Naming Conventions
- **Packages**: Short, lowercase, singular (e.g., `models`, `server`).
- **Interfaces**: Method name + "er" (e.g., `Reader`, `Worker`) or explicit (e.g., `AgentHandler`).
- **Variables**: `camelCase` for local, `CamelCase` for exported.
- **Acronyms**: Keep acronyms consistent case (e.g., `ServeHTTP`, `ID`, `URL`).

### Comments & Documentation
- **Language**: **Traditional Chinese (Taiwan)** for all comments and logic explanations.
  - Example: `// 1. 定義 Agent B 的身分卡`
- **Docstrings**: Public functions must have JSDoc-style or Go-style comments.
- **Tone**: Professional and clear.

### Error Handling
- **Explicit Checks**: Handle errors immediately.
  ```go
  if err != nil {
      return fmt.Errorf("context: %w", err) // Use %w for wrapping
  }
  ```
- **Context**: Wrap errors to provide context using `%w`.
- **Panic**: Avoid `panic` in library code; use it only in `main` for unrecoverable startup errors.

## 5. Project Structure
```
.
├── cmd/
│   ├── agent_a/    # Client Entrypoint (Assistant)
│   └── agent_b/    # Server Entrypoint (Finance)
├── models/         # Shared Data Structures (Task, Message, Artifact)
├── server/         # A2A Server Implementation (HTTP/JSON-RPC)
├── go.mod          # Module Definition
└── README.md       # Project Documentation
```

## 6. A2A Protocol Implementation Details
The codebase relies on specific structures defined in `models/`:
- **AgentCard**: Describes the agent's identity, capabilities, and skills.
- **Task**: Represents a unit of work with state (`TaskStateWorking`, `TaskStateCompleted`).
- **Message**: JSON-RPC payload containing user input or agent response.
- **Artifact**: Represents tangible outputs (e.g., reports, code) or streaming text updates.

**Server Logic (`server/`):**
- **TaskUpdateFunc**: A callback for streaming updates (SSE) to the client.
- **Handlers**: Map JSON-RPC methods (`message/send`, `message/stream`) to Go functions.

## 7. Development Workflow & Git

### Commits
- Use **Conventional Commits** format:
  - `feat: add SSE support for task streaming`
  - `fix: resolve JSON unmarshal error in handshake`
  - `docs: update AGENTS.md with build instructions`

### Adding New Features
1. **Analyze**: Understand the requirement and existing patterns.
2. **Plan**: Define the data models in `models/`.
3. **Implement**: Add logic in `server/` or `cmd/`.
4. **Verify**: Add unit tests (even if simple).
5. **Run**: Verify end-to-end with `agent_a` and `agent_b`.

## 8. Nix Configuration Reference (Setup)
If `flake.nix` is missing, use this template to set up the environment:

```nix
{
  description = "A2A Golang Dev Environment";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          golangci-lint
          gotools
        ];
      };
    };
}
```

## 9. Specific User Preferences (Context)
- **Role**: Senior Software Engineer (Taiwan-based).
- **Editor**: Neovim (Keyboard-centric).
- **Communication**: Explain in Traditional Chinese, technical terms in English.
- **Pure Functions**: Prefer functional patterns where applicable.
- **Immutable Infrastructure**: Changes to system config go to `/etc/nixos/` or flakes.

---
*Generated by Agentic Assistant on Fri Feb 13 2026*
