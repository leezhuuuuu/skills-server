# Skills Registry Server

[![Go](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-green.svg)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A lightweight, high-performance registry server designed to distribute, index, and visualize AI Agent Skills.

This server acts as the **"App Store" backend** for the AI Skills ecosystem. It serves three audiences simultaneously:
1.  **AI Agents (via MCP):** Provides a RESTful API for searching and downloading skills.
2.  **Humans (via Web UI):** Provides a modern web interface to browse documentation and file structures.
3.  **LLMs (via Context):** Provides raw Markdown endpoints optimized for direct context injection.

---

## ✨ Features

- **Live Indexing:** Automatically monitors your local `data/` directory. Drop a skill folder in, and it's instantly available via API and Web UI. No database required.
- **Dual Interface:**
  - **Web Mode:** A clean, dark-mode-first SPA built with Vue 3 + Tailwind CSS.
  - **API Mode:** Standard JSON API for `skills-mcp` client integration.
- **LLM-Friendly:** Special endpoints (`/skill.md`, `/skill/:name.md`) return pure Markdown context, perfect for pasting into Claude or ChatGPT.
- **Zero Dependency:** Compiles into a single binary with all frontend assets embedded.
- **Docker Ready:** Multi-stage Dockerfile produces a tiny Alpine-based image (~20MB).

---

## 🚀 Quick Start

### Running Locally

Prerequisites: Go 1.22+, Node.js 18+.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/leezhuuuuu/skills-server.git
    cd skills-server
    ```

2.  **Prepare Frontend:**
    ```bash
    cd web
    npm install
    npm run build  # Compiles to ../cmd/server/web_dist
    cd ..
    ```

3.  **Run Server:**
    ```bash
    # Point to the included sample data
    SKILLS_DATA_DIR=./data go run cmd/server/main.go
    ```

4.  **Access:**
    - Web UI: http://localhost:8080
    - API: http://localhost:8080/api/v1/skills
    - LLM Guide: http://localhost:8080/skill.md

### Running via Docker

```bash
docker-compose up -d --build
```
This will mount the `./data` directory into the container.

---

## 📂 Data Structure

To add a new skill, simply create a folder in your data directory. A skill **must** have a `SKILL.md` file with YAML frontmatter.

**Example Structure:**
```text
data/
└── pdf-toolkit/           <-- Folder Name (Skill ID)
    ├── SKILL.md           <-- Required: Metadata & Docs
    ├── requirements.txt   <-- Optional: Python dependencies
    └── scripts/           <-- Optional: Tool scripts
        └── main.py
```

**Example `SKILL.md`:**
```markdown
---
name: pdf-toolkit
description: Tools for manipulating PDF files (merge, split, fill forms).
author: Lee Zhu
tags: [pdf, office, utility]
version: 1.0.0
---

# PDF Toolkit
This skill provides tools to...
```

---

## 🔌 API Reference

### For MCP Clients

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/skills` | List all skills (supports `?q=query` for fuzzy search). |
| `GET` | `/api/v1/skills/:name` | Get detailed metadata and file tree. |
| `GET` | `/api/v1/download/:name` | Download skill as a `.zip` archive. |

### For LLMs (Direct Context)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/skill.md` | Returns system guide and installation instructions. |
| `GET` | `/skill/:name.md` | Returns complete skill context (metadata + file tree + docs). |

---

## 🛠️ Development

### Architecture

The project follows the standard Go layout:

*   `cmd/server`: Application entry point.
*   `internal/indexer`: Core logic for scanning filesystem and parsing Frontmatter. Uses `fsnotify` for real-time updates.
*   `internal/handlers`: HTTP handlers for API and Web routes.
*   `web`: Vue 3 frontend source code.

### Frontend Development

To work on the frontend with hot-reload:

1.  Run the Go backend (it serves the API on port 8080).
2.  In a separate terminal, run `cd web && npm run dev`.
3.  Access http://localhost:5173 (Vite proxies API requests to port 8080).

---

## 🤝 Contributing

Pull requests are welcome!

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/amazing-feature`).
3.  Commit your changes.
4.  Push to the branch.
5.  Open a Pull Request.
