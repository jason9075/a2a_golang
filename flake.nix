{
  description = "A2A Golang Dev Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            golangci-lint
            gotools
            delve # Debugger
            just  # Task runner
          ];

          shellHook = ''
            echo "🚀 Go Development Environment Loaded"
            echo "Go version: $(go version)"
          '';
        };
      }
    );
}
