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
- Custom **TCP transport layer**
- Server-side `Listen` and `Accept` loop
- Incoming connection handling
- Manual testing using **Telnet**
- Clean separation of transport logic

### 🔜 In Progress / Planned
- Message framing (length-prefixed protocol)
- Wire protocol (PUT / GET / DELETE)
- Local file storage engine
- Peer discovery & membership
- Data replication & fault tolerance

---

## 🏗️ Architecture (Early Stage)

