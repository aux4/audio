package main

import (
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// minCompressBitrate is the lowest bitrate (bits per second) compress will use when fitting
// a file under --maxSize. Below it speech stops being intelligible.
const minCompressBitrate = 8000

var compressFormats = []string{"mp3", "opus", "m4a"}

var sizePatternBytes = regexp.MustCompile(`^(?i)(\d+(?:\.\d+)?)\s*(b|kb|mb|gb|kib|mib|gib|k|m|g)?$`)

// Valid libmp3lame bitrates in kbps: MPEG-1 for 32 kHz and above, MPEG-2/2.5 below.
var mp3BitratesHigh = []int{32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}
var mp3BitratesLow = []int{8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160}

// parseBitrate converts 32k, 0.5M or 32000 into bits per second.
func parseBitrate(value string) (int, error) {
	if !bitratePattern.MatchString(value) {
		return 0, errorf("invalid bitrate %q: use a value like 32k or 32000", value)
	}
	multiplier := 1.0
	number := value
	switch strings.ToLower(value[len(value)-1:]) {
	case "k":
		multiplier, number = 1000, value[:len(value)-1]
	case "m":
		multiplier, number = 1000000, value[:len(value)-1]
	}
	parsed, _ := strconv.ParseFloat(number, 64)
	bits := int(parsed * multiplier)
	if bits < minCompressBitrate {
		return 0, errorf("invalid bitrate %q: must be at least %dk", value, minCompressBitrate/1000)
	}
	return bits, nil
}

// parseByteSize converts 24MB, 500KB, 2MiB or 1048576 into bytes. KB/MB/GB are decimal
// (1 MB = 1,000,000 bytes); KiB/MiB/GiB are binary.
func parseByteSize(value string) (int64, error) {
	match := sizePatternBytes.FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		return 0, errorf("invalid maxSize %q: use bytes (e.g. 25000000) or a size like 500KB, 24MB", value)
	}
	number, _ := strconv.ParseFloat(match[1], 64)
	units := map[string]float64{
		"": 1, "b": 1,
		"k": 1e3, "kb": 1e3, "m": 1e6, "mb": 1e6, "g": 1e9, "gb": 1e9,
		"kib": 1 << 10, "mib": 1 << 20, "gib": 1 << 30,
	}
	bytes := int64(number * units[strings.ToLower(match[2])])
	if bytes < 1024 {
		return 0, errorf("invalid maxSize %q: must be at least 1KB", value)
	}
	return bytes, nil
}

// snapBitrate rounds a bitrate down to a value the encoder accepts. It returns 0 when no
// valid bitrate is at or below the requested one.
func snapBitrate(format string, bits, sampleRate int) int {
	if format != "mp3" {
		return bits
	}
	table := mp3BitratesLow
	if sampleRate >= 32000 {
		table = mp3BitratesHigh
	}
	snapped := 0
	for _, kbps := range table {
		if kbps*1000 <= bits {
			snapped = kbps * 1000
		}
	}
	return snapped
}

// fitBitrate returns the highest bitrate whose output should stay under maxSize, leaving
// room for container overhead.
func fitBitrate(maxSize int64, duration float64) int {
	budget := float64(maxSize)*0.97 - 2048
	if budget <= 0 || duration <= 0 {
		return 0
	}
	return int(math.Floor(budget * 8 / duration))
}

func runCompress(args []string) error {
	input, output, format := arg(args, 0), arg(args, 1), strings.ToLower(arg(args, 2))
	bitrate, sampleRate, channels := arg(args, 3), arg(args, 4), arg(args, 5)
	maxSize, overwrite := arg(args, 6), isTrue(arg(args, 7))

	if format == "" {
		format = "mp3"
	}
	known := false
	for _, f := range compressFormats {
		known = known || f == format
	}
	if !known {
		return errorf("unsupported format %q: use one of %s", format, strings.Join(compressFormats, ", "))
	}
	if bitrate == "" {
		bitrate = "32k"
	}
	if sampleRate == "" {
		sampleRate = "16000"
	}
	if channels == "" {
		channels = "1"
	}

	bits, err := parseBitrate(bitrate)
	if err != nil {
		return err
	}
	// Validates sampleRate and channels with the same rules as convert.
	if _, err := encodeArgs("", sampleRate, channels); err != nil {
		return err
	}
	rate, _ := strconv.Atoi(sampleRate)
	if format == "opus" && !isOpusRate(rate) {
		return errorf("invalid sampleRate %q for opus: use 8000, 12000, 16000, 24000 or 48000", sampleRate)
	}

	var limit int64
	if maxSize != "" {
		if limit, err = parseByteSize(maxSize); err != nil {
			return err
		}
	}

	// Fitting a size limit needs the duration, so stdin is read in full first.
	in, err := openInput(input, limit > 0)
	if err != nil {
		return err
	}
	out, err := resolveOutput(output, format, in, format, overwrite)
	if err != nil {
		return err
	}

	encodeTo := func(target *mediaOutput, bitsPerSecond int) error {
		after := []string{"-map", "0:a:0", "-vn", "-map_metadata", "-1"}
		after = append(after, codecArgs(format)...)
		after = append(after, "-b:a", strconv.Itoa(bitsPerSecond), "-ar", sampleRate, "-ac", channels)
		return transcode(in, target, nil, after)
	}

	if limit == 0 {
		encoded := snapBitrate(format, bits, rate)
		if err := encodeTo(out, encoded); err != nil {
			return err
		}
		out.say("Compressed %s -> %s (%d bytes, %dk)", in.label, out.label(), outputSize(out), encoded/1000)
		return nil
	}

	info, err := probe(in.path)
	if err != nil {
		return err
	}
	if info.Duration <= 0 {
		return errorf("cannot compress %s: the audio has no duration", in.label)
	}

	// Encode to a file first so the size can be checked; when streaming, the file is a
	// temporary one that is copied to stdout once it fits.
	file := out
	if out.streaming() {
		path, err := tempFile(format)
		if err != nil {
			return err
		}
		file = &mediaOutput{path: path, format: format}
	}

	target := bits
	if fit := fitBitrate(limit, info.Duration); fit < target {
		target = fit
	}

	written := false
	for attempt := 0; attempt < 3; attempt++ {
		encoded := snapBitrate(format, target, rate)
		if encoded < minCompressBitrate {
			if written {
				file.discardPartial()
			}
			return errorf("%s (%ss) cannot fit in %s even at the minimum bitrate of %dk: split it first (e.g. aux4 audio split --segment <seconds>)",
				in.label, formatSeconds(info.Duration), maxSize, minCompressBitrate/1000)
		}

		if err := encodeTo(file, encoded); err != nil {
			return err
		}
		written = true

		stat, err := os.Stat(file.path)
		if err != nil {
			return errorf("compressed file was not written: %s", file.path)
		}
		if stat.Size() <= limit {
			if out.streaming() {
				if err := copyToStdout(file.path, out); err != nil {
					return err
				}
			}
			out.say("Compressed %s -> %s (%d bytes, %dk)", in.label, out.label(), stat.Size(), encoded/1000)
			return nil
		}

		// The encoder overshot (container overhead, bitrate floors): scale down and retry.
		target = int(float64(encoded) * float64(limit) / float64(stat.Size()) * 0.95)
	}

	file.discardPartial()
	return errorf("could not compress %s under %s: split it first (e.g. aux4 audio split --segment <seconds>)", in.label, maxSize)
}

// outputSize is the number of bytes written to the output file or stdout.
func outputSize(out *mediaOutput) int64 {
	if out.streaming() {
		if out.stdout == nil {
			return 0
		}
		return out.stdout.n
	}
	if stat, err := os.Stat(out.path); err == nil {
		return stat.Size()
	}
	return 0
}

func copyToStdout(path string, out *mediaOutput) error {
	file, err := os.Open(path)
	if err != nil {
		return errorf("cannot read %s: %s", path, err.Error())
	}
	defer file.Close()
	if _, err := io.Copy(out.writer(), file); err != nil {
		return errorf("cannot write to stdout: %s", err.Error())
	}
	return nil
}

func isOpusRate(rate int) bool {
	switch rate {
	case 8000, 12000, 16000, 24000, 48000:
		return true
	}
	return false
}
