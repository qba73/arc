package arc

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Values of these vars are passed during the build time.
var version, commit, buildtime string

type config struct {
	args   []string
	output io.Writer
}

type option func(*config)

// WithArgs constructs an option.
func WithArgs(args []string) option {
	return func(c *config) {
		c.args = args
	}
}

// WithOoutput constructs an option.
func WithOutput(out io.Writer) option {
	return func(c *config) {
		c.output = out
	}
}

// CLI runs the application with provided arguments
// and prints the output to supplied io.Writer.
func CLI(opts ...option) error {
	c := config{
		args:   os.Args[1:],
		output: os.Stdout,
	}

	for _, opt := range opts {
		opt(&c)
	}

	flagset := flag.NewFlagSet("arc", flag.ExitOnError)
	printVersion := flagset.Bool("version", false, "shows the arc version: arc -version")

	flagset.Parse(c.args)
	flagset.SetOutput(c.output)

	if *printVersion {
		fmt.Fprintf(c.output, "Version: %s\nGitRef: %s\nBuild Time: %s\n", version, commit, buildtime)
		return nil
	}

	return nil
}
