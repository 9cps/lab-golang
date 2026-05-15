# Skill: RTK (Rust Token Killer) Integration

RTK คือ CLI tool ที่ลด token consumption ของ AI assistants ได้ 60-90% โดย filter/compress output ของ commands

## Installation

### Recommended: WSL (Windows) / Linux / macOS
```bash
# WSL/Linux/macOS
curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh
brew install rtk    # macOS via Homebrew

rtk init -g         # Setup for Claude Code (creates hook + RTK.md)
rtk --version       # Verify
```

### Windows (Native)
```bash
# Download rtk-x86_64-pc-windows-msvc.zip from:
# https://github.com/rtk-ai/rtk/releases
# Extract and add to PATH, then:
rtk init -g         # Falls back to CLAUDE.md injection mode
```

---

## Token Savings for Backend Go Projects

| Category | Command | Baseline | With RTK | Savings |
|---|---|---|---|---|
| **Testing** | `go test ./...` (on fail) | 6,000 | 600 | -90% |
| **Linting** | `golangci-lint run` | 3,000 | 450 | -85% |
| **Git** | `git status` | 2,000 | 400 | -80% |
| **Git** | `git diff` | 5,000 | 1,250 | -75% |
| **Files** | `ls -la` | 800 | 160 | -80% |
| **Search** | `grep / rg` | 4,000 | 800 | -80% |
| **30-min session** | Combined | ~30,000 | ~5,000 | **-83%** |

---

## Core Commands for Go Backend

```bash
# Testing
rtk go test ./...               # Shows failures only
rtk go test ./test/...          # Unit test results
rtk test <cmd>                  # Generic test wrapper

# Linting
rtk golangci-lint run           # Grouped by file
rtk err golangci-lint run       # Errors only

# Git
rtk git status                  # Compact status
rtk git diff                    # Condensed diff
rtk git log -n 10               # One-line commits
rtk git add/commit/push         # Minimal confirmation

# File Operations
rtk ls -la                      # Tree view (-80%)
rtk read main.go                # Smart file reading
rtk read main.go -l aggressive  # Signatures only
rtk grep "pattern" .            # Grouped results
rtk find "*.go" .               # Compact find

# Analytics
rtk gain                        # Token savings summary
rtk gain --graph                # 30-day graph
rtk discover                    # Missed opportunities
```

---

## How It Works

### Without RTK
```
Claude --git status--> shell --> git --> ~2,000 tokens (raw output)
```

### With RTK
```
Claude --git status--> RTK --> git --> ~200 tokens (filtered)
```

**Strategies**:
1. Smart Filtering — removes noise (comments, whitespace)
2. Grouping — aggregates similar items (files by dir, errors by type)
3. Truncation — keeps relevant context, cuts redundancy
4. Deduplication — collapses repeated lines with counts

---

## Auto-Rewrite Hook

After `rtk init -g`, commands are transparently rewritten (Bash only):

```bash
go test ./...       # Automatically becomes: rtk go test ./...
git status          # Automatically becomes: rtk git status
ls -la              # Automatically becomes: rtk ls -la
```

**Note**: Claude Code built-in tools (Read, Grep, Glob) bypass the hook.  
Use shell commands for RTK filtering:
```bash
cat file.go         # → Use `rtk read file.go` instead
grep "pattern" .    # → Use `rtk grep "pattern" .` instead
```

---

## Configuration

### Global Config
`~/.config/rtk/config.toml` (Linux/macOS)  
`~/Library/Application Support/rtk/config.toml` (macOS)

```toml
[hooks]
exclude_commands = ["curl", "psql"]  # Skip these commands

[tee]
enabled = true          # Save full output on failure
mode = "failures"       # "failures", "always", "never"
```

### Per-Project Config (optional)
```toml
# .rtk/config.toml (project-scoped)
[filters]
go = { aggressive = true }      # More aggressive filtering
git = { truncate_lines = 50 }   # Custom truncation
```

---

## Best Practices

| Practice | Benefit |
|---|---|
| `rtk init -g` once | Automatic hook setup; zero manual overhead |
| `rtk go test` always | -90% on test failures; crucial for iteration |
| `rtk git` for commits | -80-92% on git operations |
| `rtk gain` weekly | Track savings; identify missed opportunities |
| Explicit `rtk` calls | When built-in tools (Read/Grep) are used |
| WSL on Windows | Full hook support; best experience |

---

## Supported AI Tools

RTK integrates with:
- **Claude Code** — `rtk init -g` (PreToolUse hook)
- **GitHub Copilot** — `rtk init -g --copilot`
- **Cursor** — `rtk init -g --agent cursor`
- **Windsurf** — `rtk init --agent windsurf`
- **Cline / Roo Code** — `rtk init --agent cline`
- **Gemini CLI** — `rtk init -g --gemini`
- **Hermes** — `rtk init --agent hermes`

And 6+ more (see https://rtk-ai.app/guide/getting-started/supported-agents)

---

## Troubleshooting

| Issue | Solution |
|---|---|
| Hook not activating (Windows native) | Use WSL, or call `rtk` explicitly |
| `rtk gain` fails | May have installed "Rust Type Kit" instead; use `cargo install --git https://github.com/rtk-ai/rtk` |
| Need to verify installation | Run `rtk init --show` |
| Disable telemetry | `rtk telemetry disable` or export `RTK_TELEMETRY_DISABLED=1` |

---

## Resources

- **Official Guide**: https://rtk-ai.app/guide
- **GitHub**: https://github.com/rtk-ai/rtk
- **Discord**: https://discord.gg/RySmvNF5kF
- **Architecture**: https://github.com/rtk-ai/rtk/blob/develop/docs/contributing/ARCHITECTURE.md
