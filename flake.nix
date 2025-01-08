{
  description = "bitmagnet dev shell";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
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
            go
            go-task
            golangci-lint
            jekyll
            nodejs_23
            nodePackages.prettier
            protobuf
            protoc-gen-go
            ruby
          ] ++ (if stdenv.isLinux then [
            chromium
          ] else []);
        };
      };
    });
}
