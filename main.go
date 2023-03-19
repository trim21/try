package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

const reportTpl = `
retry --limit=4 -- curl ...

--`

func main() {
	flags, cmd := partitionCommand(os.Args[1:])

	var opt = Option{}

	cli := pflag.NewFlagSet("try", pflag.ContinueOnError)
	cli.IntVar(&opt.Limit, "limit", 5, "max retry")
	if err := cli.Parse(flags); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err := opt.Retry(cmd[0], cmd[1:])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

type Option struct {
	Limit int
}

func (o Option) Retry(cmd string, args []string) error {
	if o.Limit > 0 {
		o.Limit = 3
	}

	for i := 0; i < o.Limit; i++ {
		fmt.Printf("--- retry %d run ---\n", i+1)
		c := exec.Command(cmd, args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err := c.Run()
		if err == nil {
			return nil
		}

		fmt.Printf("--- run %d failed, exit code %d ---\n", i+1, c.ProcessState.ExitCode())
	}

	return fmt.Errorf("failed to run command after retries")
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
		fmt.Println(reportTpl)
		os.Exit(1)
	}

	return args[:splitIndex], args[splitIndex+1:]
}
