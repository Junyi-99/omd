# Development Tools Preferences

This document records preferred tools and technologies for development tasks.

## Package Managers

### JavaScript/TypeScript Projects
- **Preferred**: Bun
- **Avoid**: npm, yarn, pnpm (unless project already uses them)

**Rationale**: Bun is faster and provides a better developer experience.

**Commands**:
```bash
# Install dependencies
bun install

# Run scripts
bun run <script>

# Add packages
bun add <package>

# Development
bun run dev
```

## Container Tools

### Container Orchestration
- **Preferred**: Docker
- **Also uses**: Kubernetes (k3s, k8s) for infrastructure projects

**Common workflows**:
```bash
# Build image
docker build -t <name> .

# Run container
docker run -d -p <port>:<port> <name>

# Compose
docker-compose up -d
```

## Build Tools

### General
- Language-specific build tools as appropriate
- Docker for containerization
- Kubernetes for orchestration (especially in Infrastructure projects)

## Programming Languages in Use

The user actively works with:
- **Python** - Research, trading, data analysis
- **TypeScript** - Web applications, tools
- **JavaScript** - Web development
- **Go** - Backend services, infrastructure
- **C++** - Performance-critical code
- **Swift** - iOS/macOS projects

## Infrastructure Tools

Based on user's project patterns:
- **k3s/k8s** - Kubernetes clusters
- **Docker** - Containerization
- **Configuration management** - Ansible/playbooks (xtra-playbooks)

## Notes

- When starting a new JavaScript/TypeScript project, initialize with Bun
- When working with existing projects, respect their existing package manager
- Always check for `package.json` scripts first before assuming commands
