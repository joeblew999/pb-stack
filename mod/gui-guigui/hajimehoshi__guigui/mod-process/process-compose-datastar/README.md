# Process Compose Datastar GUI


A simple, reactive web interface for Process Compose built with Datastar.

https://github.com/starfederation/datastar

## 🌟 Features

- **Real-time Process Monitoring** - Live updates via SSE
- **Process Control** - Start, stop, restart processes with one click
- **Beautiful UI** - Clean, responsive design with Tailwind CSS
- **Live Events** - Real-time event log from Process Compose
- **Zero JavaScript** - Powered by Datastar's declarative approach

## 🚀 Quick Start

1. **Start Process Compose with Huma API**:
   ```bash
   cd ../
   task run  # Starts enhanced API on :8888
   ```

2. **Start Datastar GUI**:
   ```bash
   go mod tidy
   go run main.go
   ```

3. **Open Browser**:
   - GUI: http://localhost:3000
   - API: http://localhost:8888

## 🎯 How It Works

The Datastar GUI connects to the Process Compose Huma API and provides:

- **Dashboard View** - Overview of all processes with status indicators
- **Interactive Controls** - Click buttons to control processes
- **Auto-refresh** - Updates every 5 seconds automatically
- **SSE Integration** - Real-time events streamed from Process Compose
- **Responsive Design** - Works on desktop and mobile

## 🔧 Architecture

```
Browser (Datastar) ←→ Go Server (:3000) ←→ Process Compose API (:8888)
                                        ←→ Process Compose SSE (:8888/events)
```

The Go server acts as a proxy between the Datastar frontend and the Process Compose API, handling:
- API requests (GET/POST)
- SSE event streaming
- Static file serving

