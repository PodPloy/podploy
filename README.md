# PodPloy: The Sovereign PaaS

![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-apache-v2.0-green)
![Status](https://img.shields.io/badge/status-active_development-orange)
![Architecture](https://img.shields.io/badge/arch-podman_native-purple)

**PodPloy** is an open-source, reactive container orchestration platform designed for simplicity, security, and sovereignty. Built on **Podman (Rootless)** and **Go**, it eliminates the complexity of Kubernetes and the security risks of traditional Docker daemons. It provides a "Heroku-like" experience on your own hardware.

> "Your infrastructure should be as resilient as your code, but easier to manage."

Built with **Clean Architecture** in Go and featuring native **NixOS** support.
---

## ⚡ Why PodPloy?

While tools like Coolify or Dokploy are great, PodPloy takes a different architectural approach focused on **Security-by-Design** and **Edge Capabilities**:

* **Daemonless & Rootless:** Uses Podman to run containers without root privileges, minimizing the attack surface.
* **Reactive & Autonomous:** An event-driven architecture (NATS) that self-heals via the **Lazarus Protocol**.
* **IoT & Edge Ready:** A "Headless" agent mode designed for low-resource devices (Raspberry Pi).
* **True Sovereignty:** You own your data. Native support for **Local Backups**, XChaCha20 encryption, and self-hosted registries.

---

## 🏗️ Architecture: The Trinity

PodPloy is decoupled into three standalone binaries:

1.  🧠 **The Brain (`podploy-server`)**: The control plane. Manages state (SQLite/Postgres + `ent`), exposes the GraphQL/gRPC API, and serves the embedded SvelteKit UI.
2.  💪 **The Muscle (`podploy-agent`)**: The lightweight executor. Runs on every node, manages Podman sockets, and executes Zero-Touch provisioning.
3.  🎮 **The Remote (`podploy-cli`)**: The developer tool. For CI/CD integration, log streaming, and secret management.

---

## ✨ Key Features

### 👻 Headless Agent & Edge Computing
Designed for massive fleets and hostile environments.
* **Zero-Touch Provisioning:** Agents auto-configure via `cloud-init` or env vars (`PODPLOY_JOIN_TOKEN`).
* **Low Footprint:** Optimized GC and reduced telemetry for 4G/Satellite networks.
* **OTA Updates:** Agents self-update without stopping running containers (Hot Swap).
* **Local Lockdown:** Rejects physical console access; only obeys signed gRPC commands.

### ❤️‍🩹 Resilience (Lazarus Protocol)
Stop waking up at 3 AM.
* **Crash Detection:** Millisecond detection of dead processes.
* **Smart Restarts:** Exponential back-off strategies to prevent CPU trashing.
* **Liveness Probes:** HTTP/TCP health checks that auto-restart unresponsive services.

### 🚀 Zero Downtime Deployments
* **Rolling Updates:** Spawns v2, waits for health checks, drains traffic via Caddy, and gracefully shuts down v1.
* **Atomic Switch:** No 502 errors during deployment.

### 🔒 Enterprise Security
* **Zero-Config Mesh:** Automatic Wireguard mesh network between nodes with Cross-Node balancing.
* **Audit & Compliance:** Web Terminal sessions are recorded and logged.
* **Paranoid Encryption:** All backups and secrets are encrypted at rest using **XChaCha20-Poly1305**.

### 💾 Data Sovereignty
* **Time-Machine Backups:** Policy-based snapshots to Local Filesystem, S3, MinIO, or SFTP.
* **Atomic Rollbacks:** One-click revert of Code + Config + Secrets.

---

## 🛠️ Tech Stack

| Component | Technology | Role |
| :--- | :--- | :--- |
| **Language** | Go 1.25+ | Backend & Agent binaries |
| **Runtime** | Podman | OCI Container Engine (Rootless) |
| **Cache** | ristreto | OCI Container Engine (Rootless) |
| **Database** | SQLite | State Management (via `ent` ORM) |
| **Communication** | gRPC + mTLS | Secure Node-to-Node transport |
| **Event Bus** | NATS JetStream | Async tasks & Telemetry |
| **Proxy** | Caddy (Embedded) | Ingress & Auto SSL |
| **Frontend** | SvelteKit | Single Page Application (Embedded) |
| **DevOps** | Nix / Taskfile | Reproducible builds |

---

## 🚀 Getting Started

To install Podploy on your server or try the demo:

*(Coming soon: Installation instructions for End Users)*

## 📄 License
This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## 🤝 Contributing

Are you a developer? Do you want to build the future of self-hosted orchestration?

We have a standardized development environment using **Task**, **Air**, and **Nix**.

👉 **[Read our Developer Guide (CONTRIBUTING.md)](docs/CONTRIBUTING.md)** to set up your environment and start coding.

### Prerequisites
* Linux Server (Debian/Ubuntu/NixOS recommended)
* Podman 4.0+ installed

Made with ❤️ and Go in Venezuela.