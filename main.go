package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var opt = Option{}

	cli := pflag.NewFlagSet("try", pflag.ContinueOnError)
	cli.UintVar(&opt.Limit, "limit", 5, "max retry, set limit to 0 to disable limit")
	cli.DurationVar(&opt.Delay, "delay", time.Millisecond*100, "retry delay")
	cli.DurationVar(&opt.MaxDelay, "max-delay", time.Second, "max retry delay when using non-fixed delay type")
	cli.BoolVar(&opt.Quiet, "quiet", false, "hide command stdout/stderr")
	cli.StringVar(&opt.DelayType, "delay-type", "fixed", "delay type, can 'fixed' / 'backoff' / 'off'")
	flags, cmd := partitionCommand(os.Args[1:])
	if len(cmd) == 0 {
		if slices.Contains(flags, "--version") {
			fmt.Printf("version: %s\n", version)
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("build at: %s\n", date)
			os.Exit(0)
		} else {
			// handle help message
			fmt.Println("Usage: try [flags] -- command")
			fmt.Println("\nflags:")
			cli.PrintDefaults()
			os.Exit(1)
		}
	}
	if err := cli.Parse(flags); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err := opt.Retry(cmd[0], cmd[1:])
	if opt.Quiet {
		fmt.Println()
	}
	if err != nil {
		if errors.As(err, &retry.Error{}) {
			fmt.Println("All attempts fail")
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(2)
	}
}

type Option struct {
	Limit     uint
	Delay     time.Duration
	MaxDelay  time.Duration
	DelayType string
	Quiet     bool
}

func (o Option) Retry(cmd string, args []string) error {
	var delayType retry.DelayTypeFunc
	switch o.DelayType {
	case "fixed":
		delayType = retry.FixedDelay
	case "backoff":
		delayType = retry.BackOffDelay
	case "off":
		delayType = nil
	default:
		return fmt.Errorf("unknown delay type: %s", o.DelayType)
	}

	return retry.Do(
		func() error {
			c := exec.Command(cmd, args...)
			if !o.Quiet {
				c.Stderr = os.Stderr
				c.Stdout = os.Stdout
			}
			return c.Run()
		},
		retry.Attempts(o.Limit),
		retry.Delay(o.Delay),
		retry.MaxDelay(o.MaxDelay),
		retry.DelayType(delayType),
		retry.OnRetry(func(n uint, err error) {
			if o.Quiet {
				fmt.Print(".")
			} else {
				fmt.Printf("--- failed %d time(s), err: %s ---\n", n+1, err)
			}
		}),
	)
}

func partitionCommand(args []string) ([]string, []string) {
	var splitIndex = slices.Index(args, "--")

	if splitIndex == -1 {
		return args, []string{}
	}

	return args[:splitIndex], args[splitIndex+1:]
}
