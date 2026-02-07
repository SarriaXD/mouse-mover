package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SarriaXD/mouse-mover/internal/mouse"
)

type config struct {
	minutes int
}

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		printUsage(os.Stderr)
		os.Exit(2)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if runtime.GOOS == "darwin" {
			fmt.Fprintln(os.Stderr, "hint: grant Accessibility permission to your terminal in System Settings > Privacy & Security > Accessibility")
		}
		os.Exit(1)
	}
}

func parseArgs(args []string) (config, error) {
	cfg := config{}
	if len(args) == 0 {
		return cfg, nil
	}
	if len(args) > 1 {
		return cfg, fmt.Errorf("only one time parameter is supported")
	}

	a := strings.TrimSpace(args[0])
	if a == "-h" || a == "--help" {
		printUsage(os.Stdout)
		os.Exit(0)
	}

	minutes, err := strconv.Atoi(a)
	if err != nil {
		return cfg, fmt.Errorf("invalid minutes value %q", a)
	}
	if minutes < 0 {
		return cfg, fmt.Errorf("minutes must be >= 0")
	}
	cfg.minutes = minutes
	return cfg, nil
}

func printUsage(w *os.File) {
	fmt.Fprintln(w, "Usage: mm [minutes]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  mm       # run forever")
	fmt.Fprintln(w, "  mm 120   # run for 120 minutes")
}

type mover struct {
	rng *rand.Rand
}

func run(cfg config) error {
	m := mover{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var end time.Time
	if cfg.minutes > 0 {
		end = time.Now().Add(time.Duration(cfg.minutes) * time.Minute)
	}

	fmt.Printf("mm started: minutes=%d\n", cfg.minutes)
	if !end.IsZero() {
		fmt.Printf("mm will stop at: %s\n", end.Format(time.RFC3339))
	}
	fmt.Println("press Ctrl+C to stop")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	stopCh := make(chan struct{})
	go func() {
		<-sigCh
		close(stopCh)
	}()

	for {
		if !end.IsZero() && time.Now().After(end) {
			fmt.Println("mm finished: configured duration reached")
			return nil
		}

		select {
		case <-stopCh:
			fmt.Println("mm stopped by signal")
			return nil
		default:
		}

		if err := m.humanCycle(stopCh); err != nil {
			if err == errStopped {
				fmt.Println("mm stopped by signal")
				return nil
			}
			return err
		}
	}
}

var errStopped = fmt.Errorf("stopped")

func (m mover) humanCycle(stopCh <-chan struct{}) error {
	if err := m.humanMove(stopCh); err != nil {
		return err
	}

	if m.rng.Float64() < 0.58 {
		if err := m.humanScrollBurst(stopCh); err != nil {
			return err
		}
	}

	pause := m.pauseDuration()
	if err := sleepInterruptible(stopCh, pause); err != nil {
		return err
	}
	return nil
}

func (m mover) humanMove(stopCh <-chan struct{}) error {
	pos, err := mouse.Position()
	if err != nil {
		return err
	}
	size, err := mouse.ScreenSize()
	if err != nil {
		return err
	}

	targetX, targetY := m.pickTarget(pos, size)
	steps := 18 + m.rng.Intn(48)
	wobbleAmpX := float64(1 + m.rng.Intn(4))
	wobbleAmpY := float64(1 + m.rng.Intn(4))
	wobbleFreq := 1.2 + m.rng.Float64()*3.4
	phase := m.rng.Float64() * 2 * math.Pi

	for i := 1; i <= steps; i++ {
		select {
		case <-stopCh:
			return errStopped
		default:
		}

		t := float64(i) / float64(steps)
		eased := easeInOut(t)

		x := float64(pos.X) + float64(targetX-pos.X)*eased
		y := float64(pos.Y) + float64(targetY-pos.Y)*eased

		wobbleBase := math.Sin(phase+t*wobbleFreq*2*math.Pi) * (1.0 - math.Abs(0.5-t))
		x += wobbleBase*wobbleAmpX + float64(m.rng.Intn(3)-1)
		y += math.Cos(phase+t*wobbleFreq*2*math.Pi)*wobbleAmpY + float64(m.rng.Intn(3)-1)

		mx := clamp(int(x), 0, max(0, size.X-1))
		my := clamp(int(y), 0, max(0, size.Y-1))

		if err := mouse.MoveTo(mx, my); err != nil {
			return err
		}

		sleepMs := 7 + m.rng.Intn(19)
		if m.rng.Float64() < 0.12 {
			sleepMs += 10 + m.rng.Intn(30)
		}
		if err := sleepInterruptible(stopCh, time.Duration(sleepMs)*time.Millisecond); err != nil {
			return err
		}
	}

	microFixes := m.rng.Intn(3)
	for i := 0; i < microFixes; i++ {
		select {
		case <-stopCh:
			return errStopped
		default:
		}

		fx := clamp(targetX+(m.rng.Intn(7)-3), 0, max(0, size.X-1))
		fy := clamp(targetY+(m.rng.Intn(7)-3), 0, max(0, size.Y-1))
		if err := mouse.MoveTo(fx, fy); err != nil {
			return err
		}
		if err := sleepInterruptible(stopCh, time.Duration(25+m.rng.Intn(70))*time.Millisecond); err != nil {
			return err
		}
	}

	return nil
}

func (m mover) humanScrollBurst(stopCh <-chan struct{}) error {
	bursts := 1 + m.rng.Intn(3)
	for b := 0; b < bursts; b++ {
		select {
		case <-stopCh:
			return errStopped
		default:
		}

		lines := 1 + m.rng.Intn(4)
		dir := 1
		if m.rng.Float64() < 0.48 {
			dir = -1
		}
		if b > 0 && m.rng.Float64() < 0.22 {
			dir *= -1
		}

		for i := 0; i < lines; i++ {
			select {
			case <-stopCh:
				return errStopped
			default:
			}

			if err := mouse.ScrollVertical(dir); err != nil {
				return err
			}
			if err := sleepInterruptible(stopCh, time.Duration(45+m.rng.Intn(190))*time.Millisecond); err != nil {
				return err
			}
		}

		if b < bursts-1 {
			if err := sleepInterruptible(stopCh, time.Duration(250+m.rng.Intn(1300))*time.Millisecond); err != nil {
				return err
			}
		}
	}
	return nil
}

func sleepInterruptible(stopCh <-chan struct{}, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-stopCh:
		return errStopped
	case <-t.C:
		return nil
	}
}

func (m mover) pickTarget(pos, size mouse.Point) (int, int) {
	r := m.rng.Float64()
	var dxMin, dxMax, dyMin, dyMax int

	switch {
	case r < 0.68:
		dxMin, dxMax = 20, min(180, max(30, size.X/6))
		dyMin, dyMax = 15, min(120, max(20, size.Y/6))
	case r < 0.93:
		dxMin, dxMax = 120, min(480, max(180, size.X/3))
		dyMin, dyMax = 80, min(280, max(120, size.Y/3))
	default:
		dxMin, dxMax = 220, min(900, max(260, size.X/2))
		dyMin, dyMax = 120, min(520, max(160, size.Y/2))
	}

	targetX := clamp(pos.X+randomSigned(m.rng, dxMin, dxMax), 8, max(8, size.X-8))
	targetY := clamp(pos.Y+randomSigned(m.rng, dyMin, dyMax), 8, max(8, size.Y-8))
	return targetX, targetY
}

func (m mover) pauseDuration() time.Duration {
	r := m.rng.Float64()
	switch {
	case r < 0.62:
		return time.Duration(900+m.rng.Intn(3800)) * time.Millisecond
	case r < 0.88:
		return time.Duration(4+m.rng.Intn(9)) * time.Second
	default:
		return time.Duration(12+m.rng.Intn(22)) * time.Second
	}
}

func randomSigned(rng *rand.Rand, minAbs, maxAbs int) int {
	if maxAbs <= minAbs {
		maxAbs = minAbs + 1
	}
	v := minAbs + rng.Intn(maxAbs-minAbs)
	if rng.Intn(2) == 0 {
		return -v
	}
	return v
}

func easeInOut(t float64) float64 {
	return 0.5 - 0.5*math.Cos(math.Pi*t)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
