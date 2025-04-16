# HTTP Listener

The tool is a simple HTTP listener that listens for incoming HTTP requests and prints the request details to the console. Useful for testing webhooks and other HTTP requests locally. Combined with ngrok, it can be used to expose a local server to the internet.

## Installation

```bash
go install github.com/PumpkinSeed/httplistener@latest
```

## Usage

```bash
httplistener [flags]
```

### Flags

- `host`: The host to listen on (default: `:8177`)
- `output`: The output format, options: `terminal`, `terminal-json`, `file`, `file-json` (default: `terminal`)
- `filepath`: The path to the file in case of file output

You can set the host also with an environment variable: 

```bash
HL_HOST=:8000 httplistener
```