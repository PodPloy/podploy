{
  description = "Go Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        nativeBuildInputs = [ pkgs.pkg-config ];
        buildInputs = with pkgs; [
          # --- Go Language & Compiler ---
          go
          gcc

          # --- Task Runners & Hot Reload ---
          go-task
          lefthook
          air

          # --- Code Quality & Testing ---
          go-mockery_2
          golangci-lint
          govulncheck

          # --- Database & Migrations ---
          atlas
          sqlite

          # --- Protobuf & gRPC ---
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc

          # --- Containerization & Orchestration ---
          podman

          # --- Networking ---
          slirp4netns
          runc
          conmon
          skopeo
          slirp4netns
          fuse-overlayfs
        ];

        shellHook = ''
          export GOPATH=$HOME/go
          export PATH=$GOPATH/bin:$PATH

          export CGO_ENABLED=1

          echo "🚀 PodPloy DevShell Ready"
          task setup
          echo "📦 Go $(go version) | Podman $(podman --version)"
        '';
      };
    };
}
