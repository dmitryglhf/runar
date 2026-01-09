<div align="center">

<img src="./assets/logo.svg" alt="logo" width="300"/>

# `runar`

### Zero-config script tracking

</div>

## Overview

Runar is a zero-config CLI tool that tracks your script executions. Wrap any command with `runar` to automatically log runtime, exit status, git state, and stdout with persistent storage.

## Installation

```bash
# from source
go install github.com/dmitryglhf/runar/cmd/runar@latest
```

## Usage

Wrap any command with `runar`:

```bash
runar --name "baseline" python train.py --epochs 10
runar uv run train.py --lr 0.001
runar ./scripts/experiment.sh
```

Simple example:

1. Run the command:

```bash
runar echo "hello world2"
```

```bash
[runar] ▶ run_4d848d52
─────────────────────────────────────────
hello world2
─────────────────────────────────────────
[runar] ✓ Done (exit 0) | 0s
```

2. List runs:

```bash
runar ls
```

```bash
ID            STATUS   COMMAND            DURATION
run_4d848d52  success  echo hello world2  0s
```

3. Show run details:

```bash
runar show run_4d848d52
```

```bash
ID:       run_4d848d52
Command:  echo hello world2
Status:   success
Duration: 0s
Git:      main@05ee2f7
Workdir:  /Users/mac/dev/runar
Exit:     0
Logs:     .runar/logs/run_4d848d52.log
```

4. Show run logs:

```bash
runar logs run_4d848d52
```

```bash
hello world2
```

## What gets tracked (automatically)

- Command and arguments
- Start/end time, duration
- Exit code
- stdout/stderr (saved)
- Git commit, branch, dirty state
- Working directory

## Commands

```bash
runar                           # TUI dashboard (WIP)
runar <command>                 # run and track
runar ls                        # list runs
runar ls --limit 10             # list runs with limit
runar rm <id>                   # delete run
runar show <id>                 # show details
runar logs <id>                 # show stdout/stderr
runar clean --keep 10           # keep the last 10
runar clean --older 7d          # delete older than 7 days
runar clean --keep 5 --dry-run  # see what will be deleted
```
