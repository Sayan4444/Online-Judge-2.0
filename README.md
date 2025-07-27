# Online-Judge-2.0

A full-stack **Online Judge (OJ)** platform where users can participate in contests, solve problems, and submit code securely inside a sandboxed environment.

---

## Tech Stack

- **Frontend**: [Next.js](https://nextjs.org/) (React-based UI framework)
- **Backend**: [Go](https://golang.org/) (Echo framework)
- **Worker**: [Go](https://golang.org/) (Worker pool for code execution)
- **Code Execution**: [isolate](https://github.com/ioi/isolate) (Secure sandboxing of untrusted code)
- **Message Broker**: [RabbitMQ](https://www.rabbitmq.com/) (Message broker for worker pool)

---

## Features

- User registration and authentication using JWT
- Contest and problem management
- Secure code execution via isolated worker pool
- RESTful API design with Echo
- PostgreSQL (optional) support via GORM ORM

---

## Prerequisites

- **Go** `v1.23+` – [Install Go](https://go.dev/dl/)
- **Node.js** `v20+` – for frontend
- **Docker** (optional) – for local containerized setup
- **isolate** – installed on the host (Linux only)

---

## Development guide

Make sure all dependencies are installed.

Clone the repo:

```bash
git clone https://github.com/lugnitdgp/Online-Judge-2.0.git
```

### Frontend

- Use node v20+ (you can use [nvm](https://github.com/nvm-sh/nvm) to change node version)

```bash
nvm use 20
```

- Set environment variables (refer to `.env.example`)

- Get inside frontend directory, install dependencies and run frontend

```bash
cd frontend
npm install
npm run dev
```

### Backend

- Use Go v1.23+

- Set environment variables (refer to `.env.example`)

- Get inside backend directory, run backend

```bash
cd backend
go mod tidy
go run server.go
```

### Worker

- Install Isolate locally (Linux machine)

- Remove `sudo` temporarily by adding this line to your `visudo`

```bash
<username> ALL=(ALL) NOPASSWD: ALL
```

- Get inside worker directory and run worker

```bash
cd worker
go mod tidy
go run worker.go
```

## Contribution Guide

Always use atomic commit messages (monorepo best-practice) create PR against the main branch.

- e.g. [Area]: [Component/file]: Context in one line

&copy; GNU/Linux User's Group, NIT Durgapur
