# Platform

Event-driven card generation platform using NATS JetStream + Wombat.

## Quick Start

```bash
# Install dependencies
task platform:deps

# Start platform
task platform:start

# Test connectivity
task platform:test:ping

# Stop platform
task platform:stop
```

## Services

- **NATS Server** - Message broker with JetStream
- **Wombat Runtime** - Stream processor for Go binaries
- **Goreman** - Process orchestration

## Architecture

```
NATS JetStream → Wombat → Go Binaries → NATS Results
```

## Binaries

All platform binaries installed to `./.bin/` with `platform-` prefix:
- `platform-nats-server` - NATS with JetStream
- `platform-nats` - NATS CLI
- `platform-wombat` - Synadia Connect Runtime
- `platform-connect` - Synadia Connect CLI
- `platform-goreman` - Process manager

## Configuration

- **Wombat**: `./.config/wombat.yaml`
- **Services**: `./.config/services/`
- **Procfile**: Auto-generated from Taskfile

## Versions

Pinned versions in `Taskfile-platform.yml`:
- Connect: `latest`
- Wombat: `v1.0.7-rc2`
- NATS: `latest`

## Tasks

```bash
task platform:start     # Start all services
task platform:stop      # Stop all services
task platform:status    # Check service status
task platform:monitor   # Monitor all messages
task platform:deps      # Install dependencies
```

## Links

- [Synadia Connect](https://github.com/synadia-io/connect)
- [Connect Runtime Wombat](https://github.com/synadia-io/connect-runtime-wombat)

### Benefits of This Stack

**🔄 Reactive by Design**
- Changes propagate automatically through the system
- Event-driven workflows reduce manual intervention
- Real-time updates to connected clients/games

**⚡ High Performance** 
- Both tools built in Go for speed and efficiency
- Designed for high-throughput, low-latency scenarios
- Horizontal scaling capabilities

**🛠️ Low-Code Friendly**
- YAML configuration for complex workflows
- Visual pipeline builders possible
- Declarative rather than imperative approach

**🔌 Integration Ready**
- Extensive connector ecosystems
- HTTP APIs for external integration
- WebSocket support for real-time web interfaces

**📈 Production Ready**
- Battle-tested in enterprise environments
- Monitoring and observability built-in
- Clustering and high-availability support

---

## Getting Started

### Development Approach
1. **Start Small**: Use Task runner for development workflows
2. **Add Benthos**: Begin with simple file transformation pipelines  
3. **Integrate NATS**: Add event streaming for reactive updates
4. **Scale Up**: Build visual workflow builders and admin interfaces

### Next Steps
- [ ] Experiment with Benthos for asset pipeline automation
- [ ] Set up NATS JetStream for event-driven workflows
- [ ] Create proof-of-concept reactive card generation
- [ ] Design low-code interface for workflow configuration
- [ ] Build multiplayer game state synchronization

---

*This architecture provides the foundation for a scalable, reactive platform that can grow from simple asset generation to complex, distributed game and graphics workflows.*
