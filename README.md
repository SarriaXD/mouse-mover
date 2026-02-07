# mouse-mover (`mm`)

A tiny CLI to keep your machine active by moving the mouse with office-like, non-robotic patterns.

## Features

- Default behavior: run forever (`mm`)
- Simple time setting: `mm 120` means run 120 minutes
- Human-like movement style:
  - mostly short drifts
  - sometimes medium/long moves
  - smooth easing + tiny jitters + occasional pauses
- Cross-platform implementation included:
  - Windows: supported
  - macOS: supported (requires Accessibility permission)
  - Linux: placeholder (not implemented yet)

## Quick commands (tutorial style)

```bash
# 1) Start now, run forever
mm

# 2) Work for 2 hours
mm 120

# 3) Lunch break cover (90 min)
mm 90

# 4) Fast anti-idle mode (move every 20s)
mm -i 20

# 5) Test without moving cursor
mm 5 --dry-run

# Show tutorial text again
mm --tutorial
```

## Flags

- `-m <minutes>`: run minutes, `0` means forever
- `-i <seconds>`: interval between movement cycles (default `30`)
- `--dry-run`: print movement points, do not move cursor
- `--tutorial`: print quick usage examples
- `--seed <n>`: set random seed for reproducible movement

## Build

### Build current platform

```bash
go build -o mm
```

### Build Windows binary from macOS/Linux

```bash
GOOS=windows GOARCH=amd64 go build -o mm.exe
```

### (Optional) Build macOS binary explicitly

```bash
GOOS=darwin GOARCH=arm64 go build -o mm
```

## Install locally

Put the binary in any directory already in your `PATH`, for example `~/bin`.

Then you can run:

```bash
mm
```

## Notes

- On macOS, first run may fail until terminal app is granted Accessibility permission in:
  - `System Settings > Privacy & Security > Accessibility`
