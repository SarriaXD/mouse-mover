# mouse-mover (`mm`)

A tiny CLI to keep your machine active with human-like mouse movement and wheel scrolling.

## Features

- Minimal input: only one optional time parameter
  - `mm` runs forever
  - `mm 120` runs for 120 minutes
- More realistic behavior loop:
  - smooth non-linear movement with micro jitters
  - occasional tiny corrections after move
  - random short/medium/long pauses
  - wheel scroll bursts in both directions
- Cross-platform implementation included:
  - Windows: supported
  - macOS: supported (requires Accessibility permission)
  - Linux: placeholder (not implemented yet)

## Usage

```bash
# run forever
mm

# run for 2 hours
mm 120

# show help
mm --help
```

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

Default install (recommended):

```bash
curl -L -o mm https://github.com/SarriaXD/mouse-mover/releases/latest/download/mm_darwin_arm64
chmod +x mm
sudo mv mm /usr/local/bin/mm
```

Verify:

```bash
which mm
mm --help
```

#### Update to latest version

```bash
curl -L -o /usr/local/bin/mm https://github.com/SarriaXD/mouse-mover/releases/latest/download/mm_darwin_arm64
chmod +x /usr/local/bin/mm
```

#### Uninstall

```bash
sudo rm -f /usr/local/bin/mm
which mm
```

### Windows (PowerShell)

Step 1: download `mm.exe` into your current directory.

```powershell
Invoke-WebRequest -Uri "https://github.com/SarriaXD/mouse-mover/releases/latest/download/mm_windows_amd64.exe" -OutFile "mm.exe"
```

Step 2: move `mm.exe` into a folder already in your `PATH`, or add its folder to `PATH`.

Step 3: open a new PowerShell window and verify:

```powershell
Get-Command mm
mm --help
```

## Notes

- On macOS, first run may fail until terminal app is granted Accessibility permission in:
  - `System Settings > Privacy & Security > Accessibility`
- Linux build is included in release assets, but runtime mouse movement/scrolling on Linux is currently not implemented.

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
