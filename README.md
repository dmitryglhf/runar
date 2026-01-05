<div align="center">

<img src="./assets/logo.svg" alt="logo" width="300"/>

# `runar`

### Zero-config scripts tracking

</div>

## Overview

Runar is a zero-config CLI tool that tracks your script executions. Wrap and command with `runar` to automatically log runtime, exit status, git state, and stdout with persistent storage.

## Installation

```bash
# from source
go install github.com/dmitryglhf/runar@latest
```

## Usage

Wrap any command and `runar` tracks the rest:

```bash
runar python train.py --epochs 10
runar uv run train.py --lr 0.001
runar ./scripts/experiment.sh
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
runar <command>        # run and track
runar                  # TUI dashboard
runar ls               # list experiments
runar rm <id>          # delete experiment
runar show <id>        # get details
runar logs <id>        # show stdout/stderr
```

## TUI Dashboard

```bash
$ runar
```

```
┌─ Experiments ────────────────────────────────── runar v0.1 ─┐
│                                                             │
├─────┬──────────────────┬────────┬──────────────┬────────────┤
│  #  │ Name             │ Status │ Command      │ Duration   │
├─────┼──────────────────┼────────┼──────────────┼────────────┤
│ ► 1 │ run_x7k2m9       │ ✓ done │ python tr... │ 2h 15m     │
│   2 │ run_p3n8v2       │ ✓ done │ uv run tr... │ 1h 23m     │
│   3 │ run_abc123       │ ✗ fail │ ./train.sh   │ 0h 05m     │
├─────┴──────────────────┴────────┴──────────────┴────────────┤
│ ⏎ details   d delete   / search   q quit   ? help           │
└─────────────────────────────────────────────────────────────┘
```

## Philosophy

- **Zero changes** to your code
- **Zero dependencies** in your project
- **Zero config** — just prefix with `runar`
- Works with **any runner**: python, uv, conda, bash

## Roadmap

- Add more commands
- Interactive TUI
- Python-SDK
