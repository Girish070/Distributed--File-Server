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
- **Custom TCP Transport Layer:** Server-side `Listen` and `Accept` loop with clean interface decoupling.
- **Message Framing:** Custom `LengthPrefixDecoder` to handle raw TCP byte streams reliably.
- **Wire Protocol:** Structured binary message encoding (`DataMessage`) using Go's `encoding/gob`.
- **Decoupled Architecture:** `FileServer` orchestrator that consumes messages from the networking layer via Go channels (`<-chan RPC`).
- **Local Storage Engine:** Content Addressable Storage (CAS) implementation utilizing SHA-1 hashing to create optimized, deeply nested directory structures.
- **Automated Testing:** In-memory network testing using `net.Pipe()`.

### 🔜 In Progress / Planned
- File retrieval logic (GET commands)
- File deletion logic (DELETE commands)
- Peer discovery & routing (knowing which node has which file)
- Data replication & fault tolerance
- Graceful server shutdown & resource cleanup

---

## 🏗️ Architecture 
- **Transport Layer:** Handles raw TCP sockets and byte framing.
- **Orchestrator (FileServer):** Parses RPCs and routes commands.
- **Storage Engine:** Handles disk I/O and Content Addressable folder creation.