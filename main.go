package main

import (
	"errors"
	"flag"
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

	"mouse-mover/internal/mouse"
)

const (
	defaultInterval = 30 * time.Second
)

type config struct {
	minutes      int
	interval     time.Duration
	dryRun       bool
	seed         int64
	tutorialOnly bool
}

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if cfg.tutorialOnly {
		printTutorial()
		return
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
	fs := flag.NewFlagSet("mm", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg config
	var intervalSeconds int

	fs.IntVar(&cfg.minutes, "m", 0, "run minutes; 0 means run forever")
	fs.IntVar(&intervalSeconds, "i", int(defaultInterval.Seconds()), "seconds between movement cycles")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "print actions without moving cursor")
	fs.Int64Var(&cfg.seed, "seed", 0, "random seed (0 uses current time)")
	fs.BoolVar(&cfg.tutorialOnly, "tutorial", false, "print quick usage examples")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mm [minutes] [flags]\n\n")
		fmt.Fprintln(fs.Output(), "Examples:")
		fmt.Fprintln(fs.Output(), "  mm            # run forever")
		fmt.Fprintln(fs.Output(), "  mm 120        # run for 120 minutes")
		fmt.Fprintln(fs.Output(), "  mm -m 45      # run for 45 minutes")
		fmt.Fprintln(fs.Output(), "  mm -i 20      # move every 20 seconds")
		fmt.Fprintln(fs.Output(), "  mm --tutorial # show quick tutorial commands")
		fmt.Fprintln(fs.Output(), "\nFlags:")
		fs.PrintDefaults()
	}

	flagArgs, rest, err := reorderArgs(args)
	if err != nil {
		return cfg, err
	}

	if err := fs.Parse(flagArgs); err != nil {
		return cfg, err
	}

	if intervalSeconds <= 0 {
		return cfg, errors.New("-i must be > 0")
	}
	cfg.interval = time.Duration(intervalSeconds) * time.Second

	rest = append(rest, fs.Args()...)
	if len(rest) > 1 {
		return cfg, fmt.Errorf("too many positional args: %s", strings.Join(rest, " "))
	}
	if len(rest) == 1 {
		minutes, err := strconv.Atoi(rest[0])
		if err != nil {
			return cfg, fmt.Errorf("invalid minutes value %q", rest[0])
		}
		if minutes < 0 {
			return cfg, errors.New("minutes must be >= 0")
		}
		cfg.minutes = minutes
	}
	if cfg.minutes < 0 {
		return cfg, errors.New("-m must be >= 0")
	}

	if cfg.seed == 0 {
		cfg.seed = time.Now().UnixNano()
	}

	return cfg, nil
}

func reorderArgs(args []string) ([]string, []string, error) {
	flagArgs := make([]string, 0, len(args))
	positional := make([]string, 0, 1)

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}

		if !strings.HasPrefix(a, "-") || a == "-" {
			positional = append(positional, a)
			continue
		}

		flagArgs = append(flagArgs, a)
		if takesValue(a) {
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("flag %s requires a value", a)
			}
			i++
			flagArgs = append(flagArgs, args[i])
		}
	}

	return flagArgs, positional, nil
}

func takesValue(flagArg string) bool {
	switch flagArg {
	case "-m", "-i", "--seed":
		return true
	default:
		return false
	}
}

type mover struct {
	rng *rand.Rand
}

func run(cfg config) error {
	m := mover{rng: rand.New(rand.NewSource(cfg.seed))}

	var end time.Time
	if cfg.minutes > 0 {
		end = time.Now().Add(time.Duration(cfg.minutes) * time.Minute)
	}

	fmt.Printf("mm started: minutes=%d interval=%s dry_run=%v seed=%d\n", cfg.minutes, cfg.interval, cfg.dryRun, cfg.seed)
	if !end.IsZero() {
		fmt.Printf("mm will stop at: %s\n", end.Format(time.RFC3339))
	}
	fmt.Println("press Ctrl+C to stop")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if err := m.cycle(cfg.dryRun); err != nil {
		return err
	}

	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()

	for {
		if !end.IsZero() && time.Now().After(end) {
			fmt.Println("mm finished: configured duration reached")
			return nil
		}

		select {
		case <-sigCh:
			fmt.Println("mm stopped by signal")
			return nil
		case <-ticker.C:
			if err := m.cycle(cfg.dryRun); err != nil {
				return err
			}
		}
	}
}

func (m mover) cycle(dryRun bool) error {
	pos, err := mouse.Position()
	if err != nil {
		return err
	}
	size, err := mouse.ScreenSize()
	if err != nil {
		return err
	}

	dx, dy := m.pickDelta(size)
	targetX := clamp(pos.X+dx, 10, max(10, size.X-10))
	targetY := clamp(pos.Y+dy, 10, max(10, size.Y-10))

	steps := 10 + m.rng.Intn(22)
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		eased := easeInOut(t)

		x := int(float64(pos.X) + float64(targetX-pos.X)*eased)
		y := int(float64(pos.Y) + float64(targetY-pos.Y)*eased)

		if i < steps {
			x += m.rng.Intn(5) - 2
			y += m.rng.Intn(5) - 2
		}
		x = clamp(x, 0, max(0, size.X-1))
		y = clamp(y, 0, max(0, size.Y-1))

		if dryRun {
			fmt.Printf("dry-run move: (%d,%d)\n", x, y)
		} else if err := mouse.MoveTo(x, y); err != nil {
			return err
		}
		time.Sleep(time.Duration(8+m.rng.Intn(18)) * time.Millisecond)
	}

	if m.rng.Float64() < 0.35 {
		pause := time.Duration(200+m.rng.Intn(1800)) * time.Millisecond
		time.Sleep(pause)
	}

	return nil
}

func (m mover) pickDelta(size mouse.Point) (int, int) {
	r := m.rng.Float64()
	if r < 0.65 {
		return randomSigned(m.rng, 10, 120), randomSigned(m.rng, 8, 80)
	}
	if r < 0.9 {
		return randomSigned(m.rng, 120, min(360, max(130, size.X/3))), randomSigned(m.rng, 80, min(240, max(90, size.Y/3)))
	}
	return randomSigned(m.rng, 240, min(800, max(260, size.X/2))), randomSigned(m.rng, 120, min(500, max(130, size.Y/2)))
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

func printTutorial() {
	fmt.Println("Tutorial commands (easy to memorize):")
	fmt.Println("  1) Start now, run forever:")
	fmt.Println("     mm")
	fmt.Println("  2) Work for 2 hours:")
	fmt.Println("     mm 120")
	fmt.Println("  3) Lunch break cover (90 min):")
	fmt.Println("     mm 90")
	fmt.Println("  4) Fast anti-idle mode (move every 20s):")
	fmt.Println("     mm -i 20")
	fmt.Println("  5) Test without moving cursor:")
	fmt.Println("     mm 5 --dry-run")
}
