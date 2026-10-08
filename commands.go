package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// job describes a single-input ffmpeg transcode.
type job struct {
	input     string
	output    string
	inputArgs []string // options placed before -i (e.g. -ss)
	filter    string   // -af filter graph
	args      []string // extra output options
}

// run encodes the output with the codec that matches its extension.
func (j job) run() error {
	argv := append([]string{}, j.inputArgs...)
	argv = append(argv, "-i", mediaPath(j.input), "-vn")
	if j.filter != "" {
		argv = append(argv, "-af", j.filter)
	}
	argv = append(argv, codecArgs(extOf(j.output))...)
	argv = append(argv, j.args...)
	argv = append(argv, "-y", mediaPath(j.output))
	return runFFmpeg(argv...)
}

func runInfo(args []string) error {
	input := arg(args, 0)
	if err := requireInput(input); err != nil {
		return err
	}

	info, err := probe(input)
	if err != nil {
		return err
	}

	out, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(out))
	return nil
}

func runConvert(args []string) error {
	input, format, output := arg(args, 0), strings.ToLower(arg(args, 1)), arg(args, 2)
	bitrate, sampleRate, channels, overwrite := arg(args, 3), arg(args, 4), arg(args, 5), isTrue(arg(args, 6))

	if err := requireInput(input); err != nil {
		return err
	}
	if format == "" && output != "" {
		format = extOf(output)
	}
	if format == "" {
		return errorf("provide --format (%s) or an --output file with one of these extensions", strings.Join(audioFormats, ", "))
	}
	if !isAudioFormat(format) {
		return errorf("unsupported format %q: use one of %s", format, strings.Join(audioFormats, ", "))
	}

	output = defaultOutput(output, input, "", format)
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	encode, err := encodeArgs(bitrate, sampleRate, channels)
	if err != nil {
		return err
	}

	argv := []string{"-i", mediaPath(input), "-vn"}
	argv = append(argv, codecArgs(format)...)
	argv = append(argv, encode...)
	argv = append(argv, "-f", muxerFor(format), "-y", mediaPath(output))
	if err := runFFmpeg(argv...); err != nil {
		return err
	}

	fmt.Printf("Converted %s -> %s\n", input, output)
	return nil
}

func runTrim(args []string) error {
	input, start, end, duration, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), arg(args, 4), isTrue(arg(args, 5))

	if err := requireInput(input); err != nil {
		return err
	}
	if start == "" && end == "" && duration == "" {
		return errorf("provide at least one of --start, --end or --duration")
	}
	if end != "" && duration != "" {
		return errorf("use either --end or --duration, not both")
	}

	from := 0.0
	var err error
	if start != "" {
		if from, err = parseTime("start", start); err != nil {
			return err
		}
	}

	info, err := probe(input)
	if err != nil {
		return err
	}
	if from >= info.Duration {
		return errorf("start (%ss) is beyond the end of the audio (%ss)", formatSeconds(from), formatSeconds(info.Duration))
	}

	length := info.Duration - from
	if end != "" {
		to, err := parseTime("end", end)
		if err != nil {
			return err
		}
		if to <= from {
			return errorf("end (%ss) must be after start (%ss)", formatSeconds(to), formatSeconds(from))
		}
		length = to - from
	} else if duration != "" {
		if length, err = parseTime("duration", duration); err != nil {
			return err
		}
		if length <= 0 {
			return errorf("duration must be greater than 0")
		}
	}

	output = defaultOutput(output, input, "-trimmed", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	j := job{
		input:     input,
		output:    output,
		inputArgs: []string{"-ss", formatSeconds(from)},
		args:      []string{"-t", formatSeconds(length)},
	}
	if err := j.run(); err != nil {
		return err
	}

	stop := math.Min(from+length, info.Duration)
	fmt.Printf("Trimmed %s (%ss - %ss) -> %s\n", input, formatSeconds(from), formatSeconds(math.Round(stop*1000)/1000), output)
	return nil
}

// copyContainers maps a codec to a file extension that can hold it without re-encoding.
var copyContainers = map[string]string{
	"aac":    "m4a",
	"alac":   "m4a",
	"mp3":    "mp3",
	"opus":   "ogg",
	"vorbis": "ogg",
	"flac":   "flac",
	"ac3":    "ac3",
	"eac3":   "eac3",
}

func runExtract(args []string) error {
	input, format, output, bitrate, overwrite := arg(args, 0), strings.ToLower(arg(args, 1)), arg(args, 2), arg(args, 3), isTrue(arg(args, 4))

	if err := requireInput(input); err != nil {
		return err
	}

	info, err := probe(input)
	if err != nil {
		return err
	}

	copyStream := format == "" && output == ""
	if format == "" && output != "" {
		format = extOf(output)
		copyStream = copyContainers[info.Codec] == format
	}

	argv := []string{"-i", mediaPath(input), "-map", "0:a:0", "-vn", "-sn", "-dn"}
	if copyStream {
		ext, ok := copyContainers[info.Codec]
		if !ok {
			if strings.HasPrefix(info.Codec, "pcm_") {
				ext = "wav"
			} else {
				ext = "mka"
			}
		}
		format = ext
		argv = append(argv, "-c:a", "copy")
	} else {
		if !isAudioFormat(format) {
			return errorf("unsupported format %q: use one of %s", format, strings.Join(audioFormats, ", "))
		}
		encode, err := encodeArgs(bitrate, "", "")
		if err != nil {
			return err
		}
		argv = append(argv, codecArgs(format)...)
		argv = append(argv, encode...)
	}

	output = defaultOutput(output, input, "", format)
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	argv = append(argv, "-y", mediaPath(output))
	if err := runFFmpeg(argv...); err != nil {
		return err
	}

	fmt.Printf("Extracted audio from %s -> %s\n", input, output)
	return nil
}

func runConcat(args []string) error {
	output, overwrite, inputsJSON := arg(args, 0), isTrue(arg(args, 1)), arg(args, 2)

	inputs, err := parseList(inputsJSON)
	if err != nil {
		return err
	}
	if len(inputs) < 2 {
		return errorf("concat needs at least two input files")
	}
	for _, input := range inputs {
		if err := requireInput(input); err != nil {
			return err
		}
	}
	if output == "" {
		return errorf("--output is required")
	}
	if err := checkOutput(output, overwrite, inputs...); err != nil {
		return err
	}

	first, err := probe(inputs[0])
	if err != nil {
		return err
	}
	layout := "stereo"
	if first.Channels == 1 {
		layout = "mono"
	} else if first.Channels > 2 {
		layout = strconv.Itoa(first.Channels) + "c"
	}
	rate := first.SampleRate
	if rate == 0 {
		rate = 44100
	}

	argv := []string{}
	graph := ""
	labels := ""
	for i, input := range inputs {
		argv = append(argv, "-i", mediaPath(input))
		graph += fmt.Sprintf("[%d:a:0]aresample=%d,aformat=sample_rates=%d:channel_layouts=%s[a%d];", i, rate, rate, layout, i)
		labels += fmt.Sprintf("[a%d]", i)
	}
	graph += fmt.Sprintf("%sconcat=n=%d:v=0:a=1[out]", labels, len(inputs))

	argv = append(argv, "-filter_complex", graph, "-map", "[out]")
	argv = append(argv, codecArgs(extOf(output))...)
	argv = append(argv, "-y", mediaPath(output))
	if err := runFFmpeg(argv...); err != nil {
		return err
	}

	fmt.Printf("Concatenated %d files -> %s\n", len(inputs), output)
	return nil
}

var silenceStartPattern = regexp.MustCompile(`silence_start: (-?[\d.]+)`)
var silenceEndPattern = regexp.MustCompile(`silence_end: (-?[\d.]+)`)

func runSplit(args []string) error {
	input, segment, silence := arg(args, 0), arg(args, 1), isTrue(arg(args, 2))
	threshold, minSilence, outputDir := arg(args, 3), arg(args, 4), arg(args, 5)
	prefix, format, overwrite := arg(args, 6), strings.ToLower(arg(args, 7)), isTrue(arg(args, 8))

	if err := requireInput(input); err != nil {
		return err
	}
	if silence == (segment != "") {
		return errorf("use either --segment <seconds> or --silence true")
	}

	if format == "" {
		format = extOf(input)
	}
	if !isAudioFormat(format) {
		return errorf("unsupported format %q: use one of %s", format, strings.Join(audioFormats, ", "))
	}
	if outputDir == "" {
		outputDir = "."
	}
	if prefix == "" {
		prefix = filepath.Base(basename(input))
	}
	if strings.ContainsAny(prefix, `/\`) {
		return errorf("invalid prefix %q: must not contain path separators", prefix)
	}

	info, err := probe(input)
	if err != nil {
		return err
	}

	cuts := []float64{}
	if silence {
		db, err := parseSigned("threshold", threshold)
		if err != nil {
			return err
		}
		if db > 0 {
			return errorf("invalid threshold %q: must be a level in dB at or below 0 (e.g. -35)", threshold)
		}
		gap, err := parsePositive("minSilence", minSilence)
		if err != nil {
			return err
		}

		report, err := runFFmpegReport("-i", mediaPath(input), "-vn", "-af",
			fmt.Sprintf("silencedetect=noise=%sdB:d=%s", formatSeconds(db), formatSeconds(gap)), "-f", "null", "-")
		if err != nil {
			return err
		}

		starts := silenceStartPattern.FindAllStringSubmatch(report, -1)
		ends := silenceEndPattern.FindAllStringSubmatch(report, -1)
		for i := 0; i < len(starts) && i < len(ends); i++ {
			from, to := atof(starts[i][1]), atof(ends[i][1])
			if from <= 0.001 || to >= info.Duration-0.001 {
				continue
			}
			cuts = append(cuts, (from+to)/2)
		}
	} else {
		size, err := parsePositive("segment", segment)
		if err != nil {
			return err
		}
		for at := size; at < info.Duration-0.001; at += size {
			cuts = append(cuts, at)
		}
	}

	bounds := append([]float64{0}, cuts...)
	bounds = append(bounds, info.Duration)

	outputs := []string{}
	for i := 0; i+1 < len(bounds); i++ {
		outputs = append(outputs, filepath.Join(outputDir, fmt.Sprintf("%s-%03d.%s", prefix, i+1, format)))
	}
	for _, output := range outputs {
		if err := checkOutput(output, overwrite, input); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return errorf("cannot create output directory %s: %s", outputDir, err.Error())
	}

	for i, output := range outputs {
		from, to := bounds[i], bounds[i+1]
		j := job{
			input:     input,
			output:    output,
			inputArgs: []string{"-ss", formatSeconds(from)},
			args:      []string{"-t", formatSeconds(math.Round((to-from)*1e6) / 1e6)},
		}
		if err := j.run(); err != nil {
			return err
		}
	}

	fmt.Printf("Split %s into %d segments:\n", input, len(outputs))
	for _, output := range outputs {
		fmt.Println(output)
	}
	return nil
}

type loudnormReport struct {
	InputI       string `json:"input_i"`
	InputTP      string `json:"input_tp"`
	InputLRA     string `json:"input_lra"`
	InputThresh  string `json:"input_thresh"`
	TargetOffset string `json:"target_offset"`
}

func runNormalize(args []string) error {
	input, target, truePeak, loudnessRange, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), arg(args, 4), isTrue(arg(args, 5))

	if err := requireInput(input); err != nil {
		return err
	}

	i, err := parseSigned("target", target)
	if err != nil {
		return err
	}
	if i < -70 || i > -5 {
		return errorf("invalid target %q: must be between -70 and -5 LUFS", target)
	}
	tp, err := parseSigned("truePeak", truePeak)
	if err != nil {
		return err
	}
	if tp < -9 || tp > 0 {
		return errorf("invalid truePeak %q: must be between -9 and 0 dBTP", truePeak)
	}
	lra, err := parsePositive("loudnessRange", loudnessRange)
	if err != nil {
		return err
	}
	if lra < 1 || lra > 50 {
		return errorf("invalid loudnessRange %q: must be between 1 and 50 LU", loudnessRange)
	}

	output = defaultOutput(output, input, "-normalized", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	info, err := probe(input)
	if err != nil {
		return err
	}

	base := fmt.Sprintf("loudnorm=I=%s:TP=%s:LRA=%s", formatSeconds(i), formatSeconds(tp), formatSeconds(lra))
	report, err := runFFmpegReport("-i", mediaPath(input), "-vn", "-af", base+":print_format=json", "-f", "null", "-")
	if err != nil {
		return err
	}

	open, close := strings.LastIndex(report, "{"), strings.LastIndex(report, "}")
	if open < 0 || close < open {
		return errorf("could not measure loudness of %s", input)
	}
	var measured loudnormReport
	if err := json.Unmarshal([]byte(report[open:close+1]), &measured); err != nil {
		return errorf("could not parse loudness measurement of %s: %s", input, err.Error())
	}
	for _, value := range []string{measured.InputI, measured.InputTP, measured.InputLRA, measured.InputThresh, measured.TargetOffset} {
		if !signedPattern.MatchString(value) {
			return errorf("cannot normalize %s: the audio is silent or too short to measure", input)
		}
	}

	filter := fmt.Sprintf("%s:measured_I=%s:measured_TP=%s:measured_LRA=%s:measured_thresh=%s:offset=%s:linear=true",
		base, measured.InputI, measured.InputTP, measured.InputLRA, measured.InputThresh, measured.TargetOffset)

	j := job{input: input, output: output, filter: filter}
	if info.SampleRate > 0 {
		j.args = []string{"-ar", strconv.Itoa(info.SampleRate)}
	}
	if err := j.run(); err != nil {
		return err
	}

	fmt.Printf("Normalized %s (%s LUFS) -> %s\n", input, formatSeconds(i), output)
	return nil
}

func runVolume(args []string) error {
	input, level, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), isTrue(arg(args, 3))

	if err := requireInput(input); err != nil {
		return err
	}
	if !volumePattern.MatchString(level) {
		return errorf("invalid level %q: use a multiplier (e.g. 0.5, 1.5) or decibels (e.g. 6dB, -3dB)", level)
	}

	output = defaultOutput(output, input, "-volume", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	if err := (job{input: input, output: output, filter: "volume=" + level}).run(); err != nil {
		return err
	}

	fmt.Printf("Adjusted volume of %s (%s) -> %s\n", input, level, output)
	return nil
}

func runSpeed(args []string) error {
	input, factor, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), isTrue(arg(args, 3))

	if err := requireInput(input); err != nil {
		return err
	}
	value, err := parsePositive("factor", factor)
	if err != nil {
		return err
	}
	if value < 0.1 || value > 10 {
		return errorf("invalid factor %q: must be between 0.1 and 10", factor)
	}

	output = defaultOutput(output, input, "-speed", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	if err := (job{input: input, output: output, filter: atempoChain(value)}).run(); err != nil {
		return err
	}

	fmt.Printf("Changed speed of %s (%sx) -> %s\n", input, formatSeconds(value), output)
	return nil
}

// atempoChain keeps every atempo stage inside the 0.5-2.0 range it handles best.
func atempoChain(factor float64) string {
	stages := []string{}
	for factor > 2 {
		stages = append(stages, "atempo=2")
		factor /= 2
	}
	for factor < 0.5 {
		stages = append(stages, "atempo=0.5")
		factor /= 0.5
	}
	stages = append(stages, "atempo="+strconv.FormatFloat(factor, 'f', 6, 64))
	return strings.Join(stages, ",")
}

func runFade(args []string) error {
	input, fadeIn, fadeOut, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), isTrue(arg(args, 4))

	if err := requireInput(input); err != nil {
		return err
	}
	if fadeIn == "" {
		fadeIn = "0"
	}
	if fadeOut == "" {
		fadeOut = "0"
	}
	in, err := parseNonNegative("fadeIn", fadeIn)
	if err != nil {
		return err
	}
	out, err := parseNonNegative("fadeOut", fadeOut)
	if err != nil {
		return err
	}
	if in == 0 && out == 0 {
		return errorf("provide --fadeIn and/or --fadeOut (in seconds)")
	}

	info, err := probe(input)
	if err != nil {
		return err
	}
	if in > info.Duration || out > info.Duration {
		return errorf("fade is longer than the audio (%ss)", formatSeconds(info.Duration))
	}

	output = defaultOutput(output, input, "-fade", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	filters := []string{}
	if in > 0 {
		filters = append(filters, fmt.Sprintf("afade=t=in:st=0:d=%s", formatSeconds(in)))
	}
	if out > 0 {
		filters = append(filters, fmt.Sprintf("afade=t=out:st=%s:d=%s", formatSeconds(math.Max(0, info.Duration-out)), formatSeconds(out)))
	}

	if err := (job{input: input, output: output, filter: strings.Join(filters, ",")}).run(); err != nil {
		return err
	}

	fmt.Printf("Faded %s (in: %ss, out: %ss) -> %s\n", input, formatSeconds(in), formatSeconds(out), output)
	return nil
}

func runTrimSilence(args []string) error {
	input, threshold, maxPause, all, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), isTrue(arg(args, 3)), arg(args, 4), isTrue(arg(args, 5))

	if err := requireInput(input); err != nil {
		return err
	}
	db, err := parseSigned("threshold", threshold)
	if err != nil {
		return err
	}
	if db > 0 {
		return errorf("invalid threshold %q: must be a level in dB at or below 0 (e.g. -50)", threshold)
	}
	gap, err := parsePositive("maxPause", maxPause)
	if err != nil {
		return err
	}

	output = defaultOutput(output, input, "-nosilence", "")
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	level := formatSeconds(db) + "dB"
	edge := "silenceremove=start_periods=1:start_duration=0:start_threshold=" + level
	filter := edge + ",areverse," + edge + ",areverse"
	if all {
		filter += fmt.Sprintf(",silenceremove=stop_periods=-1:stop_duration=%s:stop_threshold=%s", formatSeconds(gap), level)
	}

	if err := (job{input: input, output: output, filter: filter}).run(); err != nil {
		return err
	}

	fmt.Printf("Removed silence from %s -> %s\n", input, output)
	return nil
}

func runSpeechPrep(args []string) error {
	input, sampleRate, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), isTrue(arg(args, 3))

	if err := requireInput(input); err != nil {
		return err
	}
	if sampleRate == "" {
		sampleRate = "16000"
	}
	encode, err := encodeArgs("", sampleRate, "1")
	if err != nil {
		return err
	}

	output = defaultOutput(output, input, "-16k", "wav")
	if extOf(output) != "wav" {
		return errorf("speech-prep writes WAV audio: the output file must end in .wav (got %s)", output)
	}
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	argv := []string{"-i", mediaPath(input), "-map", "0:a:0", "-vn", "-c:a", "pcm_s16le"}
	argv = append(argv, encode...)
	argv = append(argv, "-y", mediaPath(output))
	if err := runFFmpeg(argv...); err != nil {
		return err
	}

	fmt.Printf("Prepared %s for speech recognition (%s Hz mono WAV) -> %s\n", input, sampleRate, output)
	return nil
}

func runWaveform(args []string) error {
	input, size, color, output, overwrite := arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), isTrue(arg(args, 4))

	if err := requireInput(input); err != nil {
		return err
	}
	if !sizePattern.MatchString(size) {
		return errorf("invalid size %q: use WIDTHxHEIGHT (e.g. 1200x200)", size)
	}
	if !colorPattern.MatchString(color) {
		return errorf("invalid color %q: use a color name (e.g. blue) or hex value (e.g. #3b82f6)", color)
	}

	output = defaultOutput(output, input, "-waveform", "png")
	ext := extOf(output)
	if ext != "png" && ext != "jpg" && ext != "jpeg" {
		return errorf("waveform output must be a .png or .jpg file (got %s)", output)
	}
	if err := checkOutput(output, overwrite, input); err != nil {
		return err
	}

	filter := fmt.Sprintf("[0:a:0]showwavespic=s=%s:colors=%s[out]", size, color)
	if err := runFFmpeg("-i", mediaPath(input), "-filter_complex", filter, "-map", "[out]", "-frames:v", "1", "-update", "1", "-y", mediaPath(output)); err != nil {
		return err
	}

	fmt.Printf("Rendered waveform of %s (%s) -> %s\n", input, size, output)
	return nil
}

func isAudioFormat(format string) bool {
	for _, known := range audioFormats {
		if known == format {
			return true
		}
	}
	return false
}

// muxerFor maps a format to the ffmpeg muxer name.
func muxerFor(format string) string {
	switch format {
	case "m4a":
		return "ipod"
	case "aac":
		return "adts"
	}
	return format
}

// parseList accepts a JSON array of strings (as produced for repeatable aux4 variables)
// or a single plain value.
func parseList(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if strings.HasPrefix(value, "[") {
		var list []string
		if err := json.Unmarshal([]byte(value), &list); err != nil {
			return nil, errorf("invalid input list: %s", err.Error())
		}
		return list, nil
	}
	return []string{value}, nil
}
