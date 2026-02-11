{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [ pkgs.pkg-config ];
  buildInputs = with pkgs; [
    go
    gcc           

    # --- Task Runners & Hot Reload ---
    go-task       
    lefthook      
    air           

    # --- Code Quality ---
    go-mockery    
    golangci-lint 
    govulncheck   

    # --- Database ---
    atlas         
    sqlite    

    # --- Protobuf & gRPC ---
    protobuf
    protoc-gen-go
    protoc-gen-go-grpc

    # --- Containerization & Orchestration ---
    podman
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
    export CONTAINERS_POLICY_PATH=${pkgs.skopeo.src}/default-policy.json
    export REGISTRIES_CONFIG_PATH="$PWD/.nix-podman-registries.conf"
    export PODMAN_IGNORE_CGROUPSV1_WARNING=1

    NEWUIDMAP=$(type -P newuidmap || echo "/usr/bin/newuidmap")
    NEWGIDMAP=$(type -P newgidmap || echo "/usr/bin/newgidmap")

    if [ ! -f "$NEWUIDMAP" ]; then
        echo "🚨 CRITICAL: Host system does not have newuidmap installed."
        echo "   Please install 'uidmap' or 'shadow-utils' on your host OS."
    fi

    mkdir -p .config/containers
    cat <<EOF > .config/containers/containers.conf
    [engine]
    helper_binaries_dir = ["$(dirname $NEWUIDMAP)", "/usr/libexec/podman", "/usr/local/libexec/podman"]
    EOF

    if [ ! -f "$REGISTRIES_CONFIG_PATH" ]; then
      echo 'unqualified-search-registries = ["docker.io"]' > "$REGISTRIES_CONFIG_PATH"
    fi
    export CONTAINERS_REGISTRIES_CONF="$REGISTRIES_CONFIG_PATH"
    export CONTAINERS_CONF="$PWD/.config/containers/containers.conf"

    echo "🚀 PodPloy Legacy Shell Ready"
    echo "📦 Go $(go version)"
    echo "🐳 Podman $(podman --version) (Configured for Rootless)"
  '';
}