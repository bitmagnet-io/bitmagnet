{
  description = "A basic flake with a shell";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      formatter = pkgs.alejandra;
      devShells = {
        default = pkgs.mkShell {
          packages = with pkgs; [
            bundler
            go_1_23
            go-task
            golangci-lint
            jekyll
            nodejs_22
            nodePackages.prettier
            protobuf
            protoc-gen-go
          ];
        };
      };
    });
}
