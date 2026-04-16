# Distributed File Server (DFS)

A work-in-progress **distributed file storage system** built from first principles in Go.  
This project focuses on understanding **networking, systems design, and distributed systems internals** by building each layer manually instead of relying on frameworks.

---

## 🎯 Project Goals

- Learn distributed systems by **building**, not just reading
- Implement a fault-tolerant, scalable file storage system
- Gain deep understanding of:
  - TCP networking
  - Message framing & protocols
  - Data replication
  - Node coordination

---

## 🧱 Current Status

### ✅ Implemented
- **Custom TCP Transport Layer:** Server-side `Listen`, `Accept`, and active `Dial` capabilities.
- **Message Framing:** Custom `LengthPrefixDecoder` to handle raw TCP byte streams reliably.
- **Wire Protocol:** Structured binary message encoding (`MessagePayload`) using Go's `encoding/gob` supporting PUT, GET, and DELETE commands.
- **Decoupled Architecture:** `FileServer` orchestrator that consumes messages from the networking layer via Go channels.
- **Local Storage Engine:** Content Addressable Storage (CAS) implementation utilizing SHA-1 hashing to create optimized, deeply nested directory structures, with full read/write/delete support.
- **Peer Discovery:** Bootstrap node configuration allowing servers to automatically connect and form a P2P network.

### 🔜 In Progress / Planned
- Data replication & fault tolerance (broadcasting files to peers)
- Graceful server shutdown & resource cleanup
- File chunking for handling extremely large files

---

## 🏗️ Architecture 
- **Transport Layer:** Handles raw TCP sockets and byte framing.
- **Orchestrator (FileServer):** Parses RPCs and routes commands.
- **Storage Engine:** Handles disk I/O and Content Addressable folder creation.