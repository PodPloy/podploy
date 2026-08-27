{
  description = "Go Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        nativeBuildInputs = [ pkgs.pkg-config ];
        buildInputs = with pkgs; [
          go
          gcc

          go-task
          lefthook
          air

          go-mockery_2
          golangci-lint
          govulncheck

          atlas
          sqlite

          protobuf
          protoc-gen-go
          protoc-gen-go-grpc

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
