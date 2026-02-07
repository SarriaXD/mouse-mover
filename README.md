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

## Install from GitHub Release

### macOS (arm64)

```bash
curl -L -o mm https://github.com/SarriaXD/mouse-mover/releases/latest/download/mm_darwin_arm64
chmod +x mm
mv mm /usr/local/bin/mm
```

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri "https://github.com/SarriaXD/mouse-mover/releases/latest/download/mm_windows_amd64.exe" -OutFile "mm.exe"
```

### Basic usage

```bash
# run forever (default)
mm

# run for 2 hours
mm 120

# test mode (no real cursor movement)
mm 5 --dry-run
```

## Notes

- On macOS, first run may fail until terminal app is granted Accessibility permission in:
  - `System Settings > Privacy & Security > Accessibility`
- Linux build is included in release assets, but runtime mouse movement on Linux is currently not implemented.

## Release (GitHub Actions, automatic)

This repo includes an auto-release workflow at:

- `.github/workflows/release.yml`

How it works:

1. You push a tag like `v0.1.0`
2. GitHub Actions builds binaries on 3 runners:
   - macOS arm64: `mm_darwin_arm64` and `mm_v0.1.0_darwin_arm64`
   - Windows amd64: `mm_windows_amd64.exe` and `mm_v0.1.0_windows_amd64.exe`
   - Linux amd64: `mm_linux_amd64` and `mm_v0.1.0_linux_amd64`
3. Action creates a GitHub Release and uploads those files automatically

### First-time release checklist

1. Ensure your default branch is up to date (usually `main` or `master`):

```bash
git checkout <default-branch>
git pull origin <default-branch>
```

2. Commit your changes:

```bash
git add .
git commit -m "chore: prepare release"
git push origin <default-branch>
```

3. Create and push a version tag:

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

4. Check Actions:

- Open: `https://github.com/SarriaXD/mouse-mover/actions`
- Wait for workflow `release` to finish (green)

5. Check Release page:

- Open: `https://github.com/SarriaXD/mouse-mover/releases`
- Confirm assets are attached

### Next release

Use a new tag each time, for example:

```bash
git tag -a v0.1.1 -m "v0.1.1"
git push origin v0.1.1
```

Never reuse an existing tag name.
