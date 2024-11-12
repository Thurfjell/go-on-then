# Hi!

## Installation

```bash
go mod tidy
go mod verify
```

## Usage

```bash
go run cmd/run-server/main.go
```

Open your browser and go to `http://localhost:1234/game`.

## Party mode

```bash
brew install ngrok/ngrok/ngrok
ngrok http http://localhost:1234
```
