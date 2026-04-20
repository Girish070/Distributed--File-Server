# 🚀 Go Distributed File Server (DFS)

A decentralized, peer-to-peer, content-addressable file storage network built natively in Go.

This project aims to bypass the scalability bottlenecks of traditional Distributed File Systems (like the HDFS NameNode memory limit) by utilizing a decentralized Hash Ring for metadata-free peer routing, and a deeply nested Content-Addressable Storage (CAS) engine.

## 🧱 Architecture & Features

### 1. Peer-to-Peer Networking

- **Custom TCP Transport Layer:** Server-side `Listen`, `Accept`, and active outbound `Dial` capabilities.
- **Message Framing:** Custom `LengthPrefixDecoder` to handle raw TCP byte streams reliably.
- **Wire Protocol:** Structured binary message encoding (`MessagePayload`) using Go's `encoding/gob`, supporting `PUT`, `GET`, and `DELETE` commands.

### 2. Decentralized Routing & Replication

- **Consistent Hash Ring:** Uses SHA-1 and binary tree search (`sort.Search`) to map files to nodes mathematically, completely eliminating the need for a central metadata database.
- **Dynamic Peer Discovery:** Nodes dynamically join the ring upon TCP handshake.
- **Automated Replication:** Files uploaded to any node are automatically broadcast to all connected peers, with local CAS lookups to prevent infinite replication loops.

### 3. Local Storage Engine (CAS)

- **Content Addressable Storage:** Files are stored based on the SHA-1 hash of their key, ensuring deduplication.
- **Optimized Disk I/O:** Hashes are split into deeply nested directory structures (e.g., `a1/b2/c3/...`) to prevent OS-level directory limitations during massive file ingestion.
- **Deep Clean:** The `DELETE` command not only removes the file but recursively prunes empty parent directories.

### 4. Cluster Orchestration

- **Containerized Testing:** Fully containerized using a multi-stage `Dockerfile` (Alpine Linux) and `docker-compose.yml`.
- Spins up an isolated virtual network with one seed node and multiple dynamic peers to simulate real-world cluster topologies.

---

## 🛠️ How to Run the Cluster

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.
- Go 1.25+ (if running the test client locally).

### Booting the Network

You can spin up a 3-node cluster with a single command. Open your terminal in the project root and run:

`bash
docker-compose up --build
`
_Node 1 will act as the seed node. Nodes 2 and 3 will automatically boot, read their Environment Variables, and dial Node 1 to join the Hash Ring._
