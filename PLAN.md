# Runar: Development Plan

> **runar** — zero-config experiment tracking for ML

## Философия

- **CLI wrapper** — основная фича, оборачивает любую команду
- **Zero changes** — не требует изменения кода
- **Zero config** — работает из коробки
- **Universal** — python, uv, conda, bash — что угодно

## Что трекается автоматически

```
✓ Команда и аргументы
✓ Время старта/окончания, duration
✓ Exit code (success/fail)
✓ stdout/stderr (сохраняется в файл)
✓ Git commit, branch, dirty state
✓ Working directory
```

---

## CLI Commands

```bash
runar <command>           # запуск с трекингом
runar                     # TUI dashboard
runar ls                  # список экспериментов
runar show <id>           # детали
runar logs <id>           # stdout/stderr
runar rm <id>             # удалить
```

**Примеры запуска:**
```bash
runar python train.py --epochs 10
runar uv run train.py --lr 0.001
runar ./scripts/experiment.sh
runar make train
```

---

## Архитектура

```
runar/
├── cmd/
│   └── runar/
│       └── main.go              # Entry point
├── internal/
│   ├── cli/                     # CLI commands (cobra)
│   │   ├── root.go              # runar (без аргументов → TUI)
│   │   ├── run.go               # runar <command>
│   │   ├── list.go              # runar ls
│   │   ├── show.go              # runar show <id>
│   │   ├── logs.go              # runar logs <id>
│   │   └── remove.go            # runar rm <id>
│   ├── runner/                  # Subprocess management
│   │   └── runner.go            # запуск, захват stdout/stderr
│   ├── storage/                 # SQLite
│   │   └── sqlite.go            # схема + CRUD
│   └── tui/                     # Bubble Tea (Phase 2)
│       ├── app.go
│       ├── styles.go
│       └── views/
├── go.mod
└── go.sum
```

---

## Схема данных (SQLite)

```sql
CREATE TABLE runs (
    id TEXT PRIMARY KEY,              -- run_x7k2m9
    name TEXT,                        -- опциональное имя
    command TEXT NOT NULL,            -- "python train.py --epochs 10"
    status TEXT DEFAULT 'running',    -- running, completed, failed
    exit_code INTEGER,
    git_commit TEXT,
    git_branch TEXT,
    git_dirty BOOLEAN,
    workdir TEXT,
    stdout_path TEXT,                 -- путь к файлу с stdout
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME
);

CREATE INDEX idx_runs_created ON runs(created_at DESC);
CREATE INDEX idx_runs_status ON runs(status);
```

**Структура .runar/:**
```
.runar/
├── experiments.db
└── logs/
    ├── run_x7k2m9.log
    └── run_p3n8v2.log
```

---

## Data Flow

```
runar python train.py --epochs 10
       │
       ▼
┌──────────────────────────────────────────┐
│  RUNAR CLI (Go)                          │
│                                          │
│  1. Генерирует run_id                    │
│  2. Собирает git info                    │
│  3. INSERT INTO runs (status=running)    │
│  4. Запускает subprocess                 │
│  5. Стримит stdout → terminal + file     │
│  6. UPDATE runs (status, exit_code)      │
└──────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────┐
│  .runar/experiments.db                   │
│  .runar/logs/run_x7k2m9.log              │
└──────────────────────────────────────────┘
       │
       │ fsnotify
       ▼
┌──────────────────────────────────────────┐
│  RUNAR TUI (Bubble Tea)                  │
│  Real-time updates                       │
└──────────────────────────────────────────┘
```

---

## Технологии

**Go:**
- `github.com/spf13/cobra` — CLI
- `github.com/charmbracelet/bubbletea` — TUI
- `github.com/charmbracelet/lipgloss` — стили
- `github.com/fsnotify/fsnotify` — real-time updates
- `modernc.org/sqlite` — SQLite без CGO

---

## Roadmap

### Phase 1: Core CLI ✅ in progress
- [x] Структура проекта, go mod
- [x] Cobra CLI skeleton
- [ ] SQLite storage (schema, CRUD)
- [ ] `runar ls` — список из БД
- [ ] `runar rm <id>` — удаление

### Phase 2: Runner
- [ ] `runar <command>` — запуск subprocess
- [ ] Захват stdout/stderr (tee: terminal + file)
- [ ] Сбор git info
- [ ] Генерация run_id
- [ ] Запись в SQLite

### Phase 3: TUI Dashboard
- [ ] Bubble Tea app
- [ ] List view
- [ ] Details view
- [ ] fsnotify watcher
- [ ] Keybindings

### Phase 4: Polish
- [ ] `runar show <id>`
- [ ] `runar logs <id>`
- [ ] Цветной output
- [ ] --name флаг для именования
- [ ] --json для скриптов

### Future: Python SDK (опционально)
- [ ] `from runar import log`
- [ ] Логирование метрик внутри скрипта
- [ ] Прямая запись в SQLite

---

## Примеры использования

```bash
# Базовый запуск
$ runar python train.py
▶ run_x7k2m9
│ Training started...
│ Epoch 1/10: loss=0.5234
│ ...
✓ Done in 45m 23s (exit 0)

# С именем
$ runar --name "baseline" python train.py
▶ baseline (run_x7k2m9)
│ ...

# Список экспериментов
$ runar ls
ID          NAME       STATUS  DURATION  COMMAND
run_x7k2m9  baseline   ✓ done  45m 23s   python train.py
run_p3n8v2  -          ✓ done  1h 02m    uv run train.py
run_abc123  -          ✗ fail  2m 15s    ./train.sh

# Детали
$ runar show run_x7k2m9
ID:       run_x7k2m9
Name:     baseline
Command:  python train.py
Status:   completed
Duration: 45m 23s
Git:      main@a4f2c1d (dirty)
Exit:     0

# Логи
$ runar logs run_x7k2m9
[stdout/stderr of the run]

# TUI
$ runar
[opens interactive dashboard]
```
