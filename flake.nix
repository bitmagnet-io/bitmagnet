{
  description = "bitmagnet dev shell";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let

        lintOverlay = final: prev: {
          golangci-lint = prev.callPackage ./nix/golangci-lint.nix { inherit (system) ; };
        };
        pkgs = import nixpkgs {
          system = system;
          overlays = [
            lintOverlay
          ];
        };

      in
      {
        formatter = pkgs.alejandra;
        devShells = {
          default = pkgs.mkShell {
            packages =
              with pkgs;
              [
                bundler
                go
                go-task
                golangci-lint
                jekyll
                nodejs_22
                nodePackages.prettier
                protobuf
                protoc-gen-go
                ruby
              ]
              ++ (
                if stdenv.isLinux then
                  [
                    chromium
                  ]
                else
                  [ ]
              );
          };
        };
      }
    );
}
