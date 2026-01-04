<div align="center">

<img src="./assets/logo.svg" alt="logo" width="300"/>

# `runar`

### Leave runes, not configs

</div>

## Overview

Runar is a lightweight tool for managing and executing data science experiments. It provides a simple and intuitive command-line interface for defining and running python scripts.

## Installation

To install Runar, run the following command:

```bash
pip install runar
```

## Usage

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
