package main

import (
	"fmt"
	"os"
)

// usageError is returned for invalid user input. It is printed as-is.
type usageError struct {
	msg string
}

func (e *usageError) Error() string {
	return e.msg
}

func errorf(format string, args ...interface{}) error {
	return &usageError{msg: fmt.Sprintf(format, args...)}
}

var commands = map[string]func(args []string) error{
	"info":         runInfo,
	"convert":      runConvert,
	"trim":         runTrim,
	"extract":      runExtract,
	"concat":       runConcat,
	"split":        runSplit,
	"normalize":    runNormalize,
	"volume":       runVolume,
	"speed":        runSpeed,
	"fade":         runFade,
	"trim-silence": runTrimSilence,
	"speech-prep":  runSpeechPrep,
	"waveform":     runWaveform,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: aux4-audio <command> [args...]")
		os.Exit(1)
	}

	run, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n", os.Args[1])
		os.Exit(1)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

// arg returns the positional argument at index i, or "" if absent.
func arg(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return ""
}
