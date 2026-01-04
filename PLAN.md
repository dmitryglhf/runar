# Runar: ML Experiment Runner & Tracker

> Название: **Runar** — отсылка к `run` и рунам как артефактам проекта.

## Концепция

Локальный ML experiment runner с TUI на Bubble Tea. Минималистичный инструмент в стиле `uv` — одна команда, zero configuration, just works.

**Философия:**
- Одна команда `runar` для всего
- Никаких `init`, `serve`, конфигов
- Работает как обёртка над Python-скриптами
- Автоматический трекинг без изменения кода

**Ключевые решения:**
- Storage: SQLite (прямая запись, без сервера)
- CLI: Experiment runner + TUI dashboard
- TUI: Bubble Tea + Lip Gloss
- Python SDK: Минимальный API для логирования метрик

---

## Пайплайн взаимодействия

### Основной flow: Experiment Runner

```bash
$ runar train.py --epochs 10 -n "baseline"

  ▶ Experiment: baseline (run_x7k2m9)
  ├─ Git: main@a4f2c1d
  ├─ Python: 3.11.4
  ├─ Args: --epochs 10
  │
  │ Epoch 0: loss=0.6234, acc=0.7821
  │ Epoch 1: loss=0.4521, acc=0.8234
  │ ...
  │
  ✓ Completed in 2h 15m
  ✓ Artifacts: model.pt (420MB)
```

**Что происходит автоматически:**
- Создаётся `.runar/` если не существует
- Генерируется `run_id`
- Захватывается: git state, Python version, CLI args, env
- Стримится stdout/stderr
- Мониторится процесс (время, exit code)
- Всё пишется в SQLite

**Структура .runar/:**
```
my-ml-project/
├── .runar/
│   ├── experiments.db       # SQLite: метаданные, метрики
│   └── artifacts/           # бинарные файлы
│       └── run_x7k2m9/
│           ├── model.pt
│           └── stdout.log
├── train.py
└── ...
```

---

## CLI Commands

```bash
# Запуск экспериментов
runar train.py                      # запуск с авто-трекингом
runar train.py -n "experiment-1"    # с именем
runar train.py --epochs 10          # аргументы передаются скрипту

# Dashboard
runar                               # TUI (основная команда)
runar ls                            # список экспериментов (stdout)
runar show <id>                     # детали эксперимента
runar logs <id>                     # stdout/stderr эксперимента
runar diff <id1> <id2>              # сравнить два эксперимента

# Управление
runar rm <id>                       # удалить эксперимент
runar clean                         # удалить всё (с подтверждением)

# Для скриптов/CI
runar ls --json                     # JSON output
runar rm <id> --yes                 # без подтверждения
```

**Минимальный набор:** `runar`, `runar train.py`, `runar rm`

---

## Python SDK

### Уровень 0: Zero changes (просто обёртка)

```python
# train.py — обычный скрипт, никаких изменений
import torch

model = train(data)
print(f"Accuracy: {accuracy}")  # runar захватит stdout
torch.save(model, "model.pt")
```

```bash
$ runar train.py   # трекает время, git, args, stdout
```

---

### Уровень 1: Логирование метрик

```python
# train.py
from runar import log, save, config

config({"lr": 3e-5, "epochs": 10})

for epoch in range(10):
    loss = train_epoch(model)
    log({"epoch": epoch, "loss": loss})

save("model.pt", model.state_dict())
```

**Работает и с `runar train.py`, и с `python train.py`.**

---

### Уровень 2: Явный эксперимент

```python
# train.py
from runar import Run

with Run("bert-classifier") as run:
    run.config({"model": "bert-base", "lr": 3e-5})

    for epoch in range(10):
        loss = train_epoch(model)
        run.log({"epoch": epoch, "loss": loss})

    run.save("model.pt", model.state_dict())
    run.tag("baseline")
```

---

### Уровень 3: Декоратор

```python
# train.py
from runar import experiment, log, save

@experiment
def train(lr: float = 3e-5, epochs: int = 10):
    """Аргументы функции → config эксперимента"""
    for epoch in range(epochs):
        loss = train_step()
        log({"epoch": epoch, "loss": loss})

    save("model.pt", model)
    return {"final_accuracy": 0.95}  # return → финальные метрики

if __name__ == "__main__":
    train()
```

```bash
$ runar train.py --lr 1e-4 --epochs 20
# runar парсит аргументы и передаёт в функцию
```

---

### Python SDK API

```python
from runar import (
    log,        # log({"metric": value}) — логировать метрики
    save,       # save("name", object) — сохранить артефакт
    config,     # config({"param": value}) — записать конфиг
    tag,        # tag("baseline") — добавить тег
    Run,        # with Run("name"): ... — явный эксперимент
    experiment, # @experiment — декоратор
)
```

---

## Как это работает

```
runar train.py --lr 0.001
       │
       ▼
┌──────────────────────────────────────────┐
│  RUNAR CLI (Go)                          │
│                                          │
│  1. Создаёт run_id                       │
│  2. Пишет в SQLite: run started          │
│  3. Устанавливает RUNAR_RUN_ID env       │
│  4. Запускает: python train.py --lr 0.001│
│  5. Захватывает stdout/stderr            │
│  6. Финализирует run                     │
└──────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────┐
│  Python SDK (внутри train.py)            │
│                                          │
│  RUNAR_RUN_ID = os.getenv("RUNAR_RUN_ID")│
│                                          │
│  if RUNAR_RUN_ID:                        │
│      # Пишем в существующий run          │
│  else:                                   │
│      # Создаём новый run (standalone)    │
│                                          │
│  Запись: напрямую в SQLite               │
└──────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────┐
│  .runar/experiments.db                   │
│  .runar/artifacts/<run_id>/              │
└──────────────────────────────────────────┘
       │
       │ fsnotify (file system events)
       ▼
┌──────────────────────────────────────────┐
│  RUNAR TUI (Bubble Tea)                  │
│                                          │
│  • Watches experiments.db via fsnotify   │
│  • Real-time updates без сервера         │
│  • Перечитывает данные при изменениях    │
└──────────────────────────────────────────┘
```

**Real-time без сервера:**
- TUI подписывается на изменения `experiments.db` через fsnotify
- При записи метрик из Python — TUI мгновенно обновляется
- Никакого polling, никакого HTTP — просто file system events

---

## TUI Dashboard

```bash
$ runar
```

```
┌─ Experiments ─────────────────────────────────────── runar v0.1 ─┐
│                                                                   │
│  Recent Runs                                         ↑↓ navigate  │
├─────┬──────────────────┬────────┬─────────┬─────────┬────────────┤
│  #  │ Name             │ Status │ Acc     │ Loss    │ Duration   │
├─────┼──────────────────┼────────┼─────────┼─────────┼────────────┤
│ ► 1 │ baseline         │ ✓ done │ 0.891   │ 0.187   │ 2h 15m     │
│   2 │ lower-lr         │ ✓ done │ 0.847   │ 0.234   │ 2h 34m     │
│   3 │ bert-large       │ ● run  │ 0.823...│ 0.256...│ 1h 02m     │
│   4 │ failed-exp       │ ✗ fail │ -       │ -       │ 0h 05m     │
├─────┴──────────────────┴────────┴─────────┴─────────┴────────────┤
│ ⏎ details  c compare  d delete  / search  ? help                 │
└──────────────────────────────────────────────────────────────────┘
```

**TUI Keybindings:**
- `Enter` — детали эксперимента
- `c` — сравнить выбранные
- `d` — удалить
- `D` — удалить все
- `/` — поиск
- `?` — помощь
- `q` — выход

---

## Архитектура проекта

```
runar/
├── cmd/
│   └── runar/
│       └── main.go                 # Entry point
├── internal/
│   ├── cli/                        # CLI commands
│   │   ├── root.go                 # runar (TUI)
│   │   ├── run.go                  # runar train.py
│   │   ├── list.go                 # runar ls
│   │   ├── show.go                 # runar show
│   │   ├── remove.go               # runar rm
│   │   └── diff.go                 # runar diff
│   ├── runner/                     # Experiment runner
│   │   ├── runner.go               # Запуск subprocess
│   │   ├── capture.go              # Захват stdout/stderr
│   │   └── monitor.go              # Мониторинг процесса
│   ├── tui/                        # Bubble Tea UI
│   │   ├── app.go                  # Main model
│   │   ├── watcher.go              # fsnotify for real-time updates
│   │   ├── styles.go               # Lip Gloss
│   │   ├── keys.go                 # Key bindings
│   │   └── views/
│   │       ├── list.go             # Experiments list
│   │       ├── details.go          # Run details
│   │       ├── compare.go          # Comparison
│   │       └── chart.go            # ASCII charts
│   ├── tracker/                    # Core domain
│   │   ├── run.go                  # Run struct
│   │   ├── metrics.go              # Metrics
│   │   └── artifact.go             # Artifacts
│   └── storage/                    # Persistence
│       ├── storage.go              # Interface
│       ├── sqlite.go               # SQLite implementation
│       └── schema.sql              # DB schema
├── python/                         # Python SDK
│   ├── pyproject.toml
│   └── runar/
│       ├── __init__.py             # Public API
│       ├── tracker.py              # Core tracking logic
│       ├── storage.py              # SQLite writer
│       └── decorators.py           # @experiment
├── go.mod
├── go.sum
└── README.md
```

---

## Схема данных (SQLite)

```sql
CREATE TABLE runs (
    id TEXT PRIMARY KEY,              -- run_x7k2m9
    name TEXT,                        -- "baseline"
    status TEXT DEFAULT 'running',    -- running, completed, failed
    command TEXT,                     -- "train.py --epochs 10"
    config JSON,                      -- {"lr": 3e-5, "epochs": 10}
    git_commit TEXT,
    git_branch TEXT,
    git_dirty BOOLEAN,
    python_version TEXT,
    exit_code INTEGER,
    created_at TIMESTAMP,
    finished_at TIMESTAMP,
    tags JSON                         -- ["baseline", "best"]
);

CREATE TABLE metrics (
    id INTEGER PRIMARY KEY,
    run_id TEXT REFERENCES runs(id) ON DELETE CASCADE,
    step INTEGER,
    data JSON,                        -- {"loss": 0.5, "accuracy": 0.8}
    logged_at TIMESTAMP
);

CREATE TABLE artifacts (
    id INTEGER PRIMARY KEY,
    run_id TEXT REFERENCES runs(id) ON DELETE CASCADE,
    name TEXT,                        -- "model.pt"
    path TEXT,                        -- "artifacts/run_x7k2m9/model.pt"
    size_bytes INTEGER,
    checksum TEXT,
    created_at TIMESTAMP
);

CREATE INDEX idx_metrics_run ON metrics(run_id);
CREATE INDEX idx_artifacts_run ON artifacts(run_id);
CREATE INDEX idx_runs_status ON runs(status);
CREATE INDEX idx_runs_created ON runs(created_at DESC);
```

---

## Технологии

**Go:**
- `github.com/spf13/cobra` — CLI
- `github.com/charmbracelet/bubbletea` — TUI
- `github.com/charmbracelet/lipgloss` — стили
- `github.com/charmbracelet/bubbles` — компоненты (table, viewport)
- `github.com/fsnotify/fsnotify` — real-time file watching
- `modernc.org/sqlite` — SQLite без CGO

**Python SDK:**
- Zero dependencies (только stdlib)
- `sqlite3` — запись в БД
- `json` — сериализация
- `functools` — декораторы

---

## Roadmap

### Phase 1: Core CLI + Runner
- [ ] Структура проекта с go modules
- [ ] `runar train.py` — запуск с трекингом
- [ ] Захват stdout/stderr
- [ ] SQLite storage
- [ ] `runar ls`, `runar rm`

### Phase 2: Python SDK
- [ ] `log()`, `save()`, `config()`
- [ ] Детекция RUNAR_RUN_ID
- [ ] Прямая запись в SQLite
- [ ] `@experiment` декоратор

### Phase 3: TUI Dashboard
- [ ] Bubble Tea app
- [ ] List view
- [ ] Details view с метриками
- [ ] ASCII charts
- [ ] Keybindings
- [ ] fsnotify watcher для real-time updates

### Phase 4: Advanced Features
- [ ] Compare view
- [ ] `runar diff`
- [ ] Фильтрация и поиск
- [ ] GPU monitoring (nvidia-smi)

### Phase 5: Research Extensions
- [ ] Reproducibility score
- [ ] Export для LaTeX tables
- [ ] Automatic anomaly detection

---

## Пример полного flow

```bash
# Запустить эксперимент
$ runar train.py --epochs 50 -n "baseline"
  ▶ Experiment: baseline (run_x7k2m9)
  │ Epoch 0: loss=0.6234
  │ Epoch 1: loss=0.4521
  │ ...
  ✓ Done in 1h 23m

# Ещё один эксперимент
$ runar train.py --epochs 50 --lr 0.001 -n "lower-lr"
  ▶ Experiment: lower-lr (run_p3n8v2)
  │ ...
  ✓ Done in 1h 25m

# Посмотреть результаты
$ runar
  # Открывается TUI с обоими экспериментами

# Сравнить
$ runar diff baseline lower-lr

  Config Diff:
    lr: 0.0003 → 0.001

  Metrics:
    accuracy: 0.891 → 0.847 (−0.044)
    loss:     0.187 → 0.234 (+0.047)

# Удалить неудачный
$ runar rm failed-exp
  ✓ Deleted: failed-exp (run_abc123)
```
