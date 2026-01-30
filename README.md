# Podploy 🐳

> **The Lightweight, Single-Binary, Self-Hosted Container Orchestrator.**

Podploy is a minimalist alternative to Kubernetes designed for developers who want to deploy applications on their own servers (VPS/Bare Metal) without the complexity of distributed clusters. It provides a "Heroku-like" experience on your own hardware.

Built with **Clean Architecture** in Go and featuring native **NixOS** support.

---

## ⚡ Why Podploy?

We built Podploy to cure "Kubernetes Fatigue".

* **Single Binary:** No complex cluster setup (etcd, kubelet, cni). Just one static binary.
* **Zero Dependencies:** Runs on any Linux distro. Native integration with Docker and Podman (default: podman).
* **GitOps Native:** Push to your repo, and Podploy handles the build and deployment.
* **Clean Architecture:** Built with strict Hexagonal Architecture, ensuring stability and modularity.
* **NixOS First:** First-class citizen support for NixOS via Flakes and Modules.

## 🏗️ Tech Stack

Podploy is engineered for performance and maintainability:

* **Core:** Go 1.23+
* **Architecture:** Hexagonal (Ports & Adapters)
* **Communication:** gRPC (Protobuf)
* **Data Layer:** [Ent](https://entgo.io/) (SQLite/PostgreSQL)
* **Frontend:** SvelteKit (Embedded)

## 🚀 Getting Started

To install Podploy on your server or try the demo:

*(Coming soon: Installation instructions for End Users)*

## 🤝 Contributing

Are you a developer? Do you want to build the future of self-hosted orchestration?

We have a standardized development environment using **Task**, **Air**, and **Nix**.

👉 **[Read our Developer Guide (CONTRIBUTING.md)](docs/CONTRIBUTING.md)** to set up your environment and start coding.
