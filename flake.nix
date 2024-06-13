{
  description = "terraform-provider-cala";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      nativeBuildInputs = with pkgs; [
        gnumake
        tfplugindocs
      ];
      buildInputs = with pkgs; [
        go
        gotools
      ];
      terraform-provider-cala = pkgs.buildGoModule {
        pname = "terraform-provider-cala";
        version = "0.1.0";

        src = ./.;

        vendorSha256 = null;
      };
    in
      with pkgs; {
        packages = {
          inherit terraform-provider-cala;
          default = terraform-provider-cala;
        };

        devShells.default = mkShell {
          inherit buildInputs nativeBuildInputs;
          packages = [
            alejandra
            opentofu
            gopls
            vendir
          ];
        };

        formatter = alejandra;
      });
}
