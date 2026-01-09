<div align="center">

<img src="./assets/logo.svg" alt="logo" width="300"/>

# `runar`

### Zero-config script tracking

</div>

## Why

Ever ran a script and forgot the exact command or parameters that worked? Runar automatically tracks every execution with full context — so you can always find what you ran, when, and what happened.

## Installation

```bash
go install github.com/dmitryglhf/runar/cmd/runar@latest
```

## Quick start

Run any command with `runar`:

```bash
runar python train.py --epochs 10
```

```
[runar] ▶ run_4d848d52
─────────────────────────────────────────
Training started...
Epoch 1/10: loss=0.89
...
─────────────────────────────────────────
[runar] ✓ Done (exit 0) | 2m 34s
```

List your runs:

```bash
runar ls
```

```
ID            STATUS   COMMAND                       DURATION
run_4d848d52  success  python train.py --epochs 10   2m 34s
run_a1b2c3d4  failure  python train.py --epochs 50   12m 05s
run_f7e8d9c0  success  ./scripts/preprocess.sh       45s
```

Show run details:

```bash
runar show run_4d848d52
```

```
ID:       run_4d848d52
Command:  python train.py --epochs 10
Status:   success
Duration: 2m 34s
Git:      main@a1b2c3d
Workdir:  /home/user/project
Exit:     0
Logs:     .runar/logs/run_4d848d52.log
```

View logs:

```bash
runar logs run_4d848d52
```

```
Training started...
Epoch 1/10: loss=0.89
...
```

## Use cases

- **ML experiments** — track training runs with hyperparameters and git state
- **Build scripts** — log builds to debug what changed when things break
- **Data pipelines** — keep history of ETL jobs and their outputs
- **Any long-running command** — never lose track of what you ran

## What gets tracked

Automatically captured for every run:

- Command and arguments
- Start/end time, duration
- Exit code and status
- stdout/stderr (saved to file)
- Git commit, branch, dirty state
- Working directory

## Commands

```bash
runar <command>                 # run and track
runar run <command>             # explicit run (same as above)
runar --name "experiment" <cmd> # run with custom name

runar ls                        # list all runs
runar ls --limit 10             # list last 10 runs

runar show <id>                 # show run details
runar logs <id>                 # show stdout/stderr

runar rm <id>                   # delete a run
runar clean --keep 10           # keep only last 10 runs
runar clean --older 7d          # delete runs older than 7 days
runar clean --keep 5 --dry-run  # preview what would be deleted
```

## Storage

All data is stored locally in `.runar/` directory:

```
.runar/
├── experiments.db    # SQLite database with run metadata
└── logs/             # stdout/stderr for each run
    ├── run_4d848d52.log
    └── ...
```

Add to `.gitignore`:

```
.runar/
```

## License

MIT
