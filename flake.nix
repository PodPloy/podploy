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
        buildInputs = with pkgs; [
          go
          go-task
          lefthook
          air
          go-mockery_2
          golangci-lint
          atlas
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
        ];

        shellHook = ''
          export GOPATH=$HOME/go
          export PATH=$GOPATH/bin:$PATH
          echo "🚀 Welcome to the Go & Protobuf Flake shell"
        '';
      };
    };
}
