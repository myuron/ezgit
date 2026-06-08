{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            govulncheck
            golangci-lint
          ];
        };
        apps = {
          fmt = {
            type = "app";
            program = toString (pkgs.writeShellScript "fmt" ''
              set -e
              echo "==> Running formatter..."
              go fmt ./...
            '');
          };
          lint = {
            type = "app";
            program = toString (pkgs.writeShellScript "lint" ''
              set -e
              echo "==> Running linter..."
              go vet ./...
              golangci-lint run ./...
            '');
          };
          test = {
            type = "app";
            program = toString (pkgs.writeShellScript "test" ''
              set -e
              echo "==> Running test..."
              go test -v -cover ./... -coverprofile=c.out
              go tool cover -html=c.out -o ./c.html
            '');
          };
          build = {
            type = "app";
            program = toString (pkgs.writeShellScript "build" ''
              set -e
              echo "==> Running build..."
              go build .
            '');
          };
          security = {
            type = "app";
            program = toString (pkgs.writeShellScript "security" ''
              set -e
              echo "==> Running vulnerability check..."
              govulncheck ./...
            '');
          };
          ci = {
            type = "app";
            program = toString (pkgs.writeShellScript "check" ''
              set -e
              echo "==> Running formatter..."
              go fmt ./...
              echo "==> Running linter..."
              go vet ./...
              golangci-lint run ./...
              echo "==> Running test..."
              go test -cover ./...
              echo "==> Running build..."
              go build .
            '');
          };
        };
      }
    );
}
