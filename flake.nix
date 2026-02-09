{
  description = "EVENTS";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShell = pkgs.mkShell {
          name = "EVENTS";

          buildInputs = with pkgs; [
            go
            nodejs_22
            golangci-lint
            pnpm
            libwebp
          ];
        };
      });
}