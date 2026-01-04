<div align="center">

<img src="./assets/logo.svg" alt="logo" width="300"/>

# `runar`

### Zero-config experiment tracking

</div>

## Overview

Runar is a lightweight tool for managing and executing experiments. It provides a simple and intuitive command-line interface for defining and running python scripts. This tool contains a persistent SQLite database to store experiment metadata.

## Installation

```bash
# macOS/Linux
brew install runar

# or from source
go install github.com/<username>/runar@latest
```

## Usage

Wrap any command and `runar` tracks the rest:

```bash
runar python train.py --epochs 10
runar uv run train.py --lr 0.001
runar ./scripts/experiment.sh
```

```
▶ Experiment: run_x7k2m9
├─ Command: python train.py --epochs 10
├─ Git: main@a4f2c1d (dirty)
│
│ Epoch 0: loss=0.6234
│ Epoch 1: loss=0.4521
│ ...
│
✓ Completed in 2h 15m (exit 0)
```

## What gets tracked (automatically)

- Command and arguments
- Start/end time, duration
- Exit code
- stdout/stderr (saved)
- Git commit, branch, dirty state
- Working directory

**No SDK. No config. No dependencies in your project.**

## Commands

```bash
runar <command>        # run and track
runar                  # TUI dashboard
runar ls               # list experiments
runar show <id>        # show details
runar logs <id>        # show stdout/stderr
runar rm <id>          # delete experiment
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
