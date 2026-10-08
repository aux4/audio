package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	timePattern    = regexp.MustCompile(`^(\d+:){0,2}\d+(\.\d+)?$`)
	numberPattern  = regexp.MustCompile(`^\d+(\.\d+)?$`)
	signedPattern  = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	bitratePattern = regexp.MustCompile(`^\d+(\.\d+)?[kKmM]?$`)
	volumePattern  = regexp.MustCompile(`^(\d+(\.\d+)?|-?\d+(\.\d+)?dB)$`)
	sizePattern    = regexp.MustCompile(`^\d{1,5}x\d{1,5}$`)
	colorPattern   = regexp.MustCompile(`^#?[A-Za-z0-9]{1,32}$`)
)

var audioFormats = []string{"mp3", "wav", "flac", "ogg", "opus", "m4a", "aac", "webm"}

func isTrue(value string) bool {
	return strings.EqualFold(value, "true")
}

// requireInput checks that the input file exists and is a regular file.
func requireInput(path string) error {
	if path == "" {
		return errorf("input file is required")
	}
	info, err := os.Stat(path)
	if err != nil {
		return errorf("input file not found: %s", path)
	}
	if info.IsDir() {
		return errorf("input is a directory, not a file: %s", path)
	}
	return nil
}

// checkOutput refuses to overwrite an existing file unless overwrite is set, and never
// allows writing over one of the inputs.
func checkOutput(output string, overwrite bool, inputs ...string) error {
	if output == "" {
		return errorf("output file is required")
	}

	outAbs, _ := filepath.Abs(output)
	for _, input := range inputs {
		inAbs, _ := filepath.Abs(input)
		if outAbs == inAbs {
			return errorf("output file is the same as the input: %s (use --output to choose another file)", output)
		}
	}

	if info, err := os.Stat(output); err == nil {
		if info.IsDir() {
			return errorf("output is a directory: %s", output)
		}
		if !overwrite {
			return errorf("output file already exists: %s (use --overwrite true to replace it)", output)
		}
	}
	return nil
}

func extOf(path string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
}

func basename(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

// parseTime accepts seconds (90, 1.5) or [HH:]MM:SS[.ms] and returns seconds.
func parseTime(name, value string) (float64, error) {
	if !timePattern.MatchString(value) {
		return 0, errorf("invalid %s %q: use seconds (e.g. 1.5) or [HH:]MM:SS[.ms]", name, value)
	}
	total := 0.0
	for _, part := range strings.Split(value, ":") {
		number, _ := strconv.ParseFloat(part, 64)
		total = total*60 + number
	}
	return total, nil
}

func parsePositive(name, value string) (float64, error) {
	if !numberPattern.MatchString(value) {
		return 0, errorf("invalid %s %q: must be a positive number", name, value)
	}
	number, _ := strconv.ParseFloat(value, 64)
	if number <= 0 {
		return 0, errorf("invalid %s %q: must be greater than 0", name, value)
	}
	return number, nil
}

func parseNonNegative(name, value string) (float64, error) {
	if !numberPattern.MatchString(value) {
		return 0, errorf("invalid %s %q: must be a number >= 0", name, value)
	}
	number, _ := strconv.ParseFloat(value, 64)
	return number, nil
}

func parseSigned(name, value string) (float64, error) {
	if !signedPattern.MatchString(value) {
		return 0, errorf("invalid %s %q: must be a number", name, value)
	}
	number, _ := strconv.ParseFloat(value, 64)
	return number, nil
}

func formatSeconds(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// encodeArgs validates and returns the shared --bitrate / --sampleRate / --channels arguments.
func encodeArgs(bitrate, sampleRate, channels string) ([]string, error) {
	args := []string{}
	if bitrate != "" {
		if !bitratePattern.MatchString(bitrate) {
			return nil, errorf("invalid bitrate %q: use a value like 128k or 192000", bitrate)
		}
		args = append(args, "-b:a", bitrate)
	}
	if sampleRate != "" {
		rate, err := strconv.Atoi(sampleRate)
		if err != nil || rate < 1000 || rate > 768000 {
			return nil, errorf("invalid sampleRate %q: use a rate in Hz such as 16000, 44100 or 48000", sampleRate)
		}
		args = append(args, "-ar", sampleRate)
	}
	if channels != "" {
		count, err := strconv.Atoi(channels)
		if err != nil || count < 1 || count > 8 {
			return nil, errorf("invalid channels %q: use a number from 1 to 8", channels)
		}
		args = append(args, "-ac", channels)
	}
	return args, nil
}
