{
  description = "DevEnvironment Kubernetes Controller";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go toolchain
            go_1_25
            gopls
            gotools
            go-tools

            # Kubernetes development
            kubebuilder
            kubectl
            kubernetes-controller-tools

            # Testing tools (ginkgo/gomega)
            ginkgo

            # Useful utilities
            kind # For local K8s cluster testing
            kustomize
            gnumake
            setup-envtest
            just
          ];

          shellHook = ''
            echo "DevEnvironment Controller Development Shell"
            echo "Go version: $(go version)"
            echo "Kubebuilder version: $(kubebuilder version)"

            # Setup envtest automatically
            export KUBEBUILDER_ASSETS=$(setup-envtest use -p path --use-env)

            # Add local bin to PATH for controller-runtime tools
            export PATH="$PWD/bin:$PATH"

            # Optional: Setup a local KIND cluster
            # kind create cluster --name devenv-controller || true

            just
          '';

          # Environment variables for development
          KUBEBUILDER_ASSETS = ""; # Will be set in shellHook
        };

        # Package the controller (for later deployment)
        packages.default = pkgs.buildGoModule {
          pname = "devenv-controller";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # Update this after first build

          meta = with pkgs.lib; {
            description = "Kubernetes controller for managing development environments";
            homepage = "https://github.com/machinology/devenv-controller";
            license = licenses.mit;
          };
        };
      }
    );
}
