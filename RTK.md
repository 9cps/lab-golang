# RTK.md - Token Optimization Guide for lab-golang

RTK (Rust Token Killer) reduces token consumption by 60-90% on common dev commands.

## Installation

### Windows (Recommended: WSL)
```bash
# Option 1: WSL (best experience - full hook support)
wsl
curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh
rtk init -g

# Option 2: Native Windows (limited)
# 1. Download rtk-x86_64-pc-windows-msvc.zip from https://github.com/rtk-ai/rtk/releases
# 2. Extract and add rtk.exe to PATH
# 3. rtk init -g
```

### macOS / Linux
```bash
brew install rtk       # or:
cargo install --git https://github.com/rtk-ai/rtk

rtk init -g            # Setup for Claude Code
rtk --version
```

## Commands for lab-golang (Go + Gin + GORM + PostgreSQL)

### Testing
```bash
rtk go test ./...              # Failures only (-90%)
rtk go test ./test/...         # Unit tests compact
rtk test go test -v            # Generic test wrapper
```

### Linting
```bash
rtk golangci-lint run          # Go lint grouped by file (-85%)
rtk err golangci-lint run      # Errors only
```

### Git Operations
```bash
rtk git status                 # Compact status (-80%)
rtk git diff                   # Condensed diff (-75%)
rtk git log -n 10              # One-line commits (-80%)
rtk git add/commit/push        # Minimal confirmation
```

### File Operations
```bash
rtk ls -la                     # Token-optimized tree (-80%)
rtk read main.go               # Smart file reading
rtk read main.go -l aggressive # Signatures only
rtk grep "pattern" .           # Grouped search results (-80%)
rtk find "*.go" .              # Compact find results
```

### Database & Docker
```bash
rtk docker ps                  # Compact container list
rtk docker logs <container>    # Deduplicated logs
docker exec -it <container> psql -U user -d dbname  # (manual for DB queries)
```

### Code Analysis
```bash
rtk smart main.go              # 2-line heuristic summary
rtk deps                       # Dependency tree summary
rtk json go.mod                # Structure without values
```

### Token Savings Analytics
```bash
rtk gain                       # Summary stats
rtk gain --graph               # ASCII graph (last 30 days)
rtk gain --history             # Recent command history
rtk discover                   # Find missed savings opportunities
```

## Token Savings Estimate (lab-golang project)

| Command | Baseline | With RTK | Savings |
|---|---|---|---|
| `go test ./...` (on failure) | 6,000 tokens | 600 tokens | -90% |
| `golangci-lint run` | 3,000 tokens | 450 tokens | -85% |
| `git status` | 2,000 tokens | 400 tokens | -80% |
| `git diff` | 5,000 tokens | 1,250 tokens | -75% |
| `ls -la` | 800 tokens | 160 tokens | -80% |
| **30-min session total** | ~30,000 | ~5,000 | **-83%** |

## Configuration

RTK config: `~/.config/rtk/config.toml` (macOS: `~/Library/Application Support/rtk/config.toml`)

### Exclude commands (if needed)
```toml
[hooks]
exclude_commands = ["curl", "psql"]  # Skip rewrite for these

[tee]
enabled = true          # Save full output on failure
mode = "failures"       # "failures", "always", or "never"
```

## Auto-Rewrite Hook

After `rtk init -g`, commands are automatically rewritten:
```bash
go test ./...           # → rtk go test ./...
git status              # → rtk git status
golangci-lint run       # → rtk golangci-lint run
```

**Note**: Hook works on Bash. Claude Code built-in tools (Read, Grep, Glob) bypass the hook — use shell commands for RTK filtering:
```bash
cat main.go             # → Use `rtk read main.go` instead
grep "pattern" .        # → Use `rtk grep "pattern" .` instead
```

## Workflow Tips

1. **Initialize once**:
   ```bash
   rtk init -g                 # Setup hook + RTK.md
   ```

2. **Monitor savings**:
   ```bash
   rtk gain                    # Check how many tokens you've saved
   rtk discover                # Find commands that could save more tokens
   ```

3. **Per-project overrides** (optional):
   ```toml
   # .rtk/config.toml (project-scoped)
   [filters]
   go = { aggressive = true }  # More aggressive Go filtering
   ```

4. **Restart after install**:
   After `rtk init -g`, restart Claude Code for the hook to activate.

## Troubleshooting

- **Windows native (PowerShell)**: RTK filters work but auto-rewrite hook doesn't activate. Use explicit `rtk` commands or upgrade to WSL.
- **Hook not working**: Run `rtk init --show` to verify installation.
- **Name collision**: If `rtk gain` fails, you installed "Rust Type Kit" instead. Use `cargo install --git https://github.com/rtk-ai/rtk` instead.

## Learn More

- **Full guide**: https://rtk-ai.app/guide
- **GitHub**: https://github.com/rtk-ai/rtk
- **Discord**: https://discord.gg/RySmvNF5kF
