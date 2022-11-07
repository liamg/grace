package main

import (
	"fmt"
	"os"

	"github.com/liamg/grace/filter"

	"github.com/liamg/grace/printer"

	"github.com/liamg/grace/tracer"
	"github.com/spf13/cobra"
)

var (
	flagDisableColours      = false
	flagMaxStringLen        = 16
	flagHexDumpLongStrings  = false
	flagMaxHexDumpLen       = 4096
	flagPID                 = 0
	flagForwardIO           = false
	flagMaxObjectProperties = 2
	flagVerbose             = false
	flagExtraNewLine        = false
	flagMultiline           = false
	flagFilter              = ""
	flagAbsoluteTimestamps  = false
	flagRelativeTimestamps  = false
	flagSummarise           = false
	flagSortKey             = ""
	flagShowSyscallNumber   = false
	flagFilterPassing       = false
	flagFilterFailing       = false
	flagOutputFile          = ""
	flagRawOutput           = false
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
			t, err = tracer.FromCommand(!flagForwardIO, args[0], args[1:]...)
			if err != nil {
				return err
			}
		}

		output := cmd.OutOrStdout()
		if flagOutputFile != "" {
			output, err = os.Create(flagOutputFile)
			if err != nil {
				return err
			}
			defer func() {
				_ = output.(*os.File).Close()
			}()
		}

		p := printer.New(output)

		p.SetUseColours(!flagDisableColours && flagOutputFile == "")
		p.SetMaxStringLen(flagMaxStringLen)
		p.SetMaxHexDumpLen(flagMaxHexDumpLen)
		p.SetExtraNewLine(flagExtraNewLine)
		p.SetMultiLine(flagMultiline)
		p.SetHexDumpLongStrings(flagHexDumpLongStrings)
		p.SetShowAbsoluteTimestamps(flagAbsoluteTimestamps)
		p.SetShowRelativeTimestamps(flagRelativeTimestamps)
		p.SetShowSyscallNumber(flagShowSyscallNumber)
		p.SetRawOutput(flagRawOutput)

		if flagVerbose {
			p.SetMaxObjectProperties(0)
		} else {
			p.SetMaxObjectProperties(flagMaxObjectProperties)
		}

		fltr, err := filter.Parse(flagFilter)
		if err != nil {
			return fmt.Errorf("failed to parse filter: %s", err)
		}
		fltr.SetFailingOnly(flagFilterFailing)
		fltr.SetPassingOnly(flagFilterPassing)
		p.SetFilter(fltr)

		if flagSummarise {
			configureSummary(t, output, flagSortKey)
		} else {
			t.SetSyscallEnterHandler(p.PrintSyscallEnter)
			t.SetSyscallExitHandler(p.PrintSyscallExit)
			t.SetSignalHandler(p.PrintSignal)
			t.SetProcessExitHandler(p.PrintProcessExit)
			t.SetAttachHandler(p.PrintAttach)
			t.SetDetachHandler(p.PrintDetach)
		}

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
	rootCmd.Flags().BoolVarP(&flagForwardIO, "forward-io", "F", flagForwardIO, "forward stdin/stdout/stderr for the given command")
	rootCmd.Flags().IntVarP(&flagMaxObjectProperties, "max-object-properties", "O", flagMaxObjectProperties, "maximum number of properties to print for objects (recursive) - this also applies to array elements")
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", flagVerbose, "enable verbose output (overrides other verbosity settings)")
	rootCmd.Flags().BoolVarP(&flagExtraNewLine, "extra-newline", "n", flagExtraNewLine, "print an extra newline after each syscall to aid readability")
	rootCmd.Flags().BoolVarP(&flagMultiline, "multiline", "m", flagMultiline, "print each syscall argument on a separate line to aid readability")
	rootCmd.Flags().StringVarP(&flagFilter, "filter", "f", flagFilter, "Filter string to apply to output. The string should be formatted as a query string e.g. 'syscall=write&arg0=stdout'. The syscall parameter filters syscalls by name. The path parameter filters syscalls that reference a particular path. The ret parameter filters by return value (values for ret are assumed to be decimal unless prefixed with 0x). Each parameter can be specified multiple times with an OR match being appied to parameters of that type, and an AND match applied to parameters of a different type.")
	rootCmd.Flags().BoolVarP(&flagAbsoluteTimestamps, "absolute-timestamps", "a", flagAbsoluteTimestamps, "print absolute timestamps for each event")
	rootCmd.Flags().BoolVarP(&flagRelativeTimestamps, "relative-timestamps", "r", flagRelativeTimestamps, "print relative timestamps for each event")
	rootCmd.Flags().BoolVarP(&flagSummarise, "summary", "S", flagSummarise, "summarise counts of all syscalls")
	rootCmd.Flags().StringVarP(&flagSortKey, "sort-column", "c", flagSortKey, "sort key for summary output (time, seconds, count, errors) (default is sort by syscall name)")
	rootCmd.Flags().BoolVarP(&flagShowSyscallNumber, "number", "N", flagShowSyscallNumber, "show syscall numbers in output")
	rootCmd.Flags().BoolVarP(&flagFilterFailing, "only-failing", "Z", flagFilterFailing, "show only failing syscalls")
	rootCmd.Flags().BoolVarP(&flagFilterPassing, "only-passing", "z", flagFilterPassing, "show only passing syscalls")
	rootCmd.Flags().StringVarP(&flagOutputFile, "output-file", "o", flagOutputFile, "output file (default is stdout)")
	rootCmd.Flags().BoolVarP(&flagRawOutput, "raw", "R", flagRawOutput, "Raw output format for arguments and return values (format everything as raw hex values)")
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
