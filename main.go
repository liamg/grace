package main

import (
	"fmt"
	"os"

	"github.com/liamg/grace/printer"

	"github.com/liamg/grace/tracer"
	"github.com/spf13/cobra"
)

var (
	flagDisableColours      = false
	flagMaxStringLen        = 32
	flagHexDumpLongStrings  = true
	flagMaxHexDumpLen       = 4096
	flagPID                 = 0
	flagSuppressOutput      = false
	flagMaxObjectProperties = 2
	flagVerbose             = false
	flagExtraNewLine        = false
	flagMultiline           = false
)

var rootCmd = &cobra.Command{
	Use:     "grace [flags] [command [args]]",
	Example: `grace -- cat /etc/passwd`,
	Short: `grace is a CLI tool for monitoring and modifying syscalls for a given process.

It's essentially strace, in Go, with colours and pretty output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		if len(args) == 0 && flagPID == 0 {
			return cmd.Help()
		}

		var t *tracer.Tracer
		var err error
		if flagPID > 0 {
			t = tracer.New(flagPID)
		} else {
			t, err = tracer.FromCommand(flagSuppressOutput, args[0], args[1:]...)
			if err != nil {
				return err
			}
		}

		p := printer.New(cmd.OutOrStdout())

		p.SetUseColours(!flagDisableColours)
		p.SetMaxStringLen(flagMaxStringLen)
		p.SetMaxHexDumpLen(flagMaxHexDumpLen)
		p.SetExtraNewLine(flagExtraNewLine)
		p.SetMultiLine(flagMultiline)

		if flagVerbose {
			p.SetHexDumpLongStrings(true)
			p.SetMaxObjectProperties(0)
		} else {
			p.SetHexDumpLongStrings(flagHexDumpLongStrings)
			p.SetMaxObjectProperties(flagMaxObjectProperties)
		}

		t.SetSyscallEnterHandler(p.PrintSyscallEnter)
		t.SetSyscallExitHandler(p.PrintSyscallExit)
		t.SetSignalHandler(p.PrintSignal)
		t.SetProcessExitHandler(p.PrintProcessExit)

		// TODO: set signal handler!

		defer func() { _, _ = fmt.Fprintln(cmd.ErrOrStderr(), "") }()

		return t.Start()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&flagDisableColours, "no-colours", "C", flagDisableColours, "disable colours in output")
	rootCmd.Flags().IntVarP(&flagMaxStringLen, "max-string-len", "s", flagMaxStringLen, "maximum length of strings to print")
	rootCmd.Flags().BoolVarP(&flagHexDumpLongStrings, "hex-dump-long-strings", "x", flagHexDumpLongStrings, "hex dump strings longer than --max-string-len")
	rootCmd.Flags().IntVarP(&flagMaxHexDumpLen, "max-hex-dump-len", "l", flagMaxHexDumpLen, "maximum length of hex dumps")
	rootCmd.Flags().IntVarP(&flagPID, "pid", "p", flagPID, "trace an existing process by PID")
	rootCmd.Flags().BoolVarP(&flagSuppressOutput, "suppress-output", "S", flagSuppressOutput, "suppress output of command")
	rootCmd.Flags().IntVarP(&flagMaxObjectProperties, "max-object-properties", "o", flagMaxObjectProperties, "maximum number of properties to print for objects (recursive) - this also applies to array elements")
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", flagVerbose, "enable verbose output (overrides other verbosity settings)")
	rootCmd.Flags().BoolVarP(&flagExtraNewLine, "extra-newline", "n", flagExtraNewLine, "print an extra newline after each syscall to aid readability")
	rootCmd.Flags().BoolVarP(&flagMultiline, "multiline", "m", flagMultiline, "print each syscall argument on a separate line to aid readability")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		if err.Error() == "no such process" {
			os.Exit(0)
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
