# kit

**kit** is a modern, minimal Git implementation built from scratch in Go. It provides both a command-line interface and a RESTful API for version control, making it a great learning resource and a foundation for custom workflows.

---

## What is kit?

kit is a simplified version control system inspired by Git. It lets you create repositories, stage and commit files, manage branches, and inspect history—all using Go. You can interact with kit via CLI or HTTP endpoints.

---

## Core Concepts

- **Repository**: Each user gets a workspace with a `.kit` directory storing objects, refs, and index.
- **Objects**: Files and directories are stored as compressed objects, similar to Git blobs and trees.
- **Index**: Tracks staged files before committing.
- **Commits**: Snapshots of the project, with author, message, and parent info.
- **Branches**: Multiple lines of development, switchable and creatable.

---

## Features

- Initialize a new repository
- Add files and folders to the index
- Commit changes with messages
- View commit logs
- Check repository status
- Restore files to previous states
- Create and checkout branches
- Reset branch to previous commit
- List branches and directories
- Web API for all major operations

---

## Getting Started

### Prerequisites

- Go 1.18+ installed

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/kit.git
cd kit
```

Build the CLI:

```bash
go build -o kit ./cmd/kit
```

---

## Usage (CLI)

Initialize a new repository:

```bash
./kit init
```

Add files to the index:

```bash
./kit add <filename>
```

Commit changes:

```bash
./kit commit -m "Your commit message"
```

Check status:

```bash
./kit status
```

View logs:

```bash
./kit log
```

Create a branch:

```bash
./kit branch <branch_name>
```

Checkout a branch:

```bash
./kit checkout <branch_name>
```

---

## REST API Endpoints

All endpoints are prefixed with `/api/v1`.

| Endpoint    | Method | Description              | Example Payload / Params                           |
| ----------- | ------ | ------------------------ | -------------------------------------------------- |
| `/init`     | POST   | Initialize repository    | `{ "username": "alice" }`                          |
| `/add`      | POST   | Add files/folders        | `{ "username": "alice", "files": ["."] }`          |
| `/commit`   | POST   | Commit changes           | `{ "username": "alice", "message": "msg" }`        |
| `/log`      | POST   | Get commit logs          | `{ "username": "alice", "count": 5 }`              |
| `/status`   | POST   | Get repo status          | `{ "username": "alice" }`                          |
| `/restore`  | POST   | Restore files            | `{ "username": "alice", "paths": ["file.txt"] }`   |
| `/branch`   | POST   | Create branch            | `{ "username": "alice", "branch": "dev" }`         |
| `/checkout` | POST   | Checkout branch          | `{ "username": "alice", "branch": "dev" }`         |
| `/branches` | GET    | List branches            | `?username=alice`                                  |
| `/dir`      | GET    | List user directories    | None                                               |
| `/check`    | GET    | Check if user dir exists | `?username=alice`                                  |
| `/reset`    | POST   | Reset branch to commit   | `{ "username": "alice", "hash": "<commit_hash>" }` |

---

## Project Structure

```
cmd/kit/           # CLI entry point
cmd/handlers/      # HTTP handlers for REST API
internals/git/     # Core git logic (init, add, commit, log, status, restore, branch, reset)
internals/utils/   # Utility functions (index, tree, file ops)
pkg/               # Shared types and helpers
routes/            # API route registration
workspaces/        # User repositories (created at runtime)
```

---

## Example Workflow

1. **Initialize**: Create a new workspace for a user.
2. **Add**: Stage files or folders for commit.
3. **Commit**: Save changes with a message.
4. **Branch**: Create and switch between branches.
5. **Status/Log**: Inspect repository state and history.
6. **Restore/Reset**: Revert files or branch to previous versions.

---

## Contribution

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## License

MIT
