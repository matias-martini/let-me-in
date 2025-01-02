# 🚀 **Let Me In**

**A secure and user-friendly terminal access service for managing multiple terminal sessions via the web.**  

---

## 📚 **Overview**

**Let Me In** is a web-based terminal management service that provides **secure and persistent access** to terminal sessions through a web interface. Designed with simplicity and extensibility in mind, the initial version focuses on **user authentication**, **session persistence**, and **basic permission management**.

The project aims to provide **on-demand, ephemeral terminal sessions** that can recover gracefully from WebSocket disconnections while maintaining a secure and auditable environment.

---

## 🎯 **Core Goals for MVP**

1. **Authentication:**  
   - Simple login system with **username and password**.  
   - A single level of permissions: **All Access** if logged in.

2. **Web-Based Terminal Access:**  
   - WebSocket connection for **real-time interaction** with terminal sessions.  
   - Support for **multiple simultaneous sessions per user**.

3. **Session Management:**  
   - Persist terminal sessions even if the WebSocket connection is lost.  
   - Configurable timeout for session persistence (e.g., sessions persist for 5 minutes after disconnection).  
   - Ability to **reconnect to active sessions** seamlessly.

4. **Terminal Security:**  
   - Sessions are **isolated per user**.  
   - All communications are encrypted using **TLS**.

5. **Configuration:**  
   - Configurable **session persistence timeout**.  
   - Easily adjustable **environment variables** for security settings.

---

## 🛠️ **Architecture Plan**

### ✅ **MVP Architecture: Single Service Design**

**Components:**
- **User Authentication Module:** Handles user login and basic authorization.
- **WebSocket Terminal Module:** Manages terminal sessions over WebSocket connections.
- **Session Persistence Module:** Keeps sessions alive temporarily after disconnections.
- **Configuration Module:** Allows easy adjustment of session timeout and other parameters.
- **Frontend:** Allow a user friendly interface to manage and access terminal sessions.

**Workflow:**
1. User logs in with credentials.  
2. Upon successful login, the user is granted **full access** to terminal sessions.  
3. A WebSocket connection is established, and a **terminal session** is started.  
4. If the WebSocket connection drops, the session persists for a configurable timeout (e.g., 5 minutes).  
5. The user can reconnect and resume the active session within the timeout window.  
6. After the timeout, inactive sessions are **terminated automatically**.

---

## 🛤️ **Milestones**

### 🚀 **Milestone 1: Minimal Authentication Mechanism**
- [x] Implement a **JWT-based authentication system**.
- [x] Provide **login and token issuance** endpoints.
- [ ] Protect API routes using **middleware for JWT validation**.
- [x] Store user credentials securely (hashed passwords with salt and pepper).
- [ ] Create basic session management (e.g., token expiration and refresh).

**✅ Deliverables:**  
- `/auth/login` endpoint  
- JWT-based auth middleware  
- Secure user credential storage  

---

### 🚀 **Milestone 2: Terminal Session Management API**
- [ ] Design the **Terminal Session model**:
   - Unique session ID
   - Owner (user ID)
   - State (`active`, `disconnected`, `terminated`)
   - Timestamp tracking (created, last active, expires)
- [ ] Implement an **API to create, list, retrieve, and terminate terminal sessions**.
- [ ] Add basic validation and error handling (e.g., session ownership checks).

**✅ Deliverables:**  
- `/sessions` (list user sessions)  
- `/sessions/:id` (retrieve session details)  
- `/sessions` (create a session)  
- `/sessions/:id/terminate` (terminate a session)  

---

### 🚀 **Milestone 3: Terminal Process Controller**
**(Find a better name – e.g., Terminal Process Orchestrator)**  
- [ ] Create a **controller** to manage low-level terminal sessions:
   - Start and manage terminal processes (`bash`, `sh`, or configurable shell).
   - Attach `stdin`, `stdout`, and `stderr` streams to WebSocket connections.
   - Store **session buffers** to allow clients to retrieve recent history when reconnecting.
- [ ] Support **multiple concurrent terminal processes** per user.
- [ ] Implement graceful **session termination** and cleanup processes.
- [ ] Add support for **session timeouts** (configurable via environment variables).

**✅ Deliverables:**  
- Session creation (`/sessions/:id/start`)  
- Session attachment via WebSocket (`/sessions/:id/connect`)  
- Session termination (`/sessions/:id/terminate`)  

---

### 🚀 **Milestone 4: Web Interface**
- [ ] Create a **minimal frontend** for the MVP.
- [ ] Implement a **login page** with JWT-based authentication.
- [ ] Create a **session dashboard**:
   - List active sessions.
   - Create new terminal sessions.
   - Terminate sessions.
- [ ] Build a **terminal interface**:
   - Real-time interaction via WebSocket.
   - Display session buffers on reconnect.
- [ ] Ensure the interface handles **WebSocket disconnections gracefully**.

**✅ Deliverables:**  
- `/login` (user authentication interface)  
- `/dashboard` (manage terminal sessions)  
- `/terminal/:id` (interactive terminal interface)  

---

### 🚀 **Milestone 5: Session Persistence and Recovery**
- [ ] Implement **session persistence** when WebSocket connections drop.
- [ ] Store session metadata and buffers temporarily (in-memory or lightweight store).
- [ ] Add a **configurable session timeout** (e.g., persist sessions for 5 minutes).
- [ ] Allow users to **reconnect to active sessions** seamlessly.
- [ ] Automatically terminate expired sessions.

**✅ Deliverables:**  
- Configurable `SESSION_TIMEOUT` environment variable  
- Reconnection logic in terminal controller  
- Automatic cleanup of stale sessions  

---

### 🚀 **Milestone 6: Configuration and Deployment**
- [ ] Introduce **environment variables** for core configuration:
   - `SESSION_TIMEOUT`
   - `JWT_SECRET`
   - `TLS_CERT` and `TLS_KEY` paths
- [ ] Create a **Dockerfile** for the application.
- [ ] Write a **Kubernetes deployment manifest**.
- [ ] Add documentation for environment variable configurations.

**✅ Deliverables:**  
- `.env` file with default configuration  
- Docker and Kubernetes manifests  
- Clear deployment documentation  

---

### 🚀 **Milestone 7: Documentation**
- [ ] Write **API documentation** (OpenAPI/Swagger).
- [ ] Document **terminal session behavior and lifecycle**.
- [ ] Provide setup and deployment guides.
- [ ] Add a **contribution guide** for open-source contributions.

**✅ Deliverables:**  
- API documentation (Swagger/OpenAPI)  
- README with setup instructions  
- Contribution guide  

---

### 🚀 **Possible Future Enhancements (Post-MVP)**

- Implement **Role-Based Access Control (RBAC)** for granular permissions.
- Enable **advanced logging and auditing dashboards**.
- Add support for **external authentication providers (e.g., OAuth, LDAP)**.
- Introduce **multi-node architecture** for scalability.
- Improve security (tls certificates and security protocols).
- Create **terminal session activity dashboards**.

---

## ✅ **Milestone Roadmap Summary**

1. **Minimal Authentication Mechanism**
2. **Terminal Session Management API**
3. **Terminal Process Controller**
4. **Minimal Web Interface**
5. **Session Persistence and Recovery**
6. **Configuration and Deployment**
7. **Documentation**

---

## 🧠 **Next Steps**

- Start with **Milestone 1: Minimal Authentication Mechanism**.  
- Plan incremental releases, with clear milestones and testing between each phase.  
- Adjust architecture as needed based on feedback during each milestone.

---
