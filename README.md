# kit

**kit** is a minimal Git implementation built from scratch in Go. It aims to help you understand the internals of Git by providing a simple, readable codebase and command-line interface.

## Features

- Initialize a new kit repository
- Add files to the index
- Create and write objects
- Compute file hashes
- Basic index management

## Getting Started

### Prerequisites

- Go 1.18+ installed

### Installation

Clone the repository:

```bash
git clone https://github.com/m21power/kit.git
cd kit
```

Build the CLI:

```bash
go build -o kit ./cmd/kit
```

### Usage

Initialize a new repository:

```bash
./kit init
```

Add a file to the index:

```bash
./kit add <filename>
```

## Project Structure

```
cmd/kit/         # CLI entry point
internals/git/   # Core git logic (add, hash, init)
internals/utils/ # Utility functions
pkg/             # Shared packages
```

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

MIT
