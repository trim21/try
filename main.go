package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/spf13/pflag"
)

func main() {
	var opt = Option{}

	cli := pflag.NewFlagSet("try", pflag.ContinueOnError)
	cli.UintVar(&opt.Limit, "limit", 5, "max retry, set limit to 0 to disable limit")
	cli.DurationVar(&opt.Delay, "delay", time.Millisecond*100, "retry delay")
	cli.BoolVar(&opt.Quiet, "quiet", false, "hide command stdout/stderr")
	flags, cmd := partitionCommand(os.Args[1:])
	if len(cmd) == 0 {
		// handle help message
		fmt.Println("Usage: try [flags] -- command")
		fmt.Println("\nflags:")
		cli.PrintDefaults()
		os.Exit(1)
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
	Limit uint
	Delay time.Duration
	Quiet bool
}

func (o Option) Retry(cmd string, args []string) error {
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
		retry.DelayType(retry.FixedDelay),
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
	var splitIndex = -1
	for i, arg := range args {
		if arg == "--" {
			splitIndex = i
			break
		}
	}

	if splitIndex == -1 {
		return args, []string{}
	}

	return args[:splitIndex], args[splitIndex+1:]
}
