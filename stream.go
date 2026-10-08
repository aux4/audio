package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// Temporary files (stdin spools, size-limited encodes) are registered here so they are
// removed on every exit path, including Ctrl-C.
var (
	cleanupMu    sync.Mutex
	cleanupFiles []string
)

func registerTemp(path string) {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	cleanupFiles = append(cleanupFiles, path)
}

func removeTemps() {
	cleanupMu.Lock()
	defer cleanupMu.Unlock()
	for _, path := range cleanupFiles {
		_ = os.Remove(path)
	}
	cleanupFiles = nil
}

func handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGPIPE)
	go func() {
		<-signals
		removeTemps()
		os.Exit(130)
	}()
}

func tempFile(ext string) (string, error) {
	suffix := ""
	if ext != "" {
		suffix = "." + ext
	}
	file, err := os.CreateTemp("", "aux4-audio-*"+suffix)
	if err != nil {
		return "", errorf("cannot create a temporary file: %s", err.Error())
	}
	registerTemp(file.Name())
	_ = file.Close()
	return file.Name(), nil
}

// mediaInput is where a command reads audio from: a file, or stdin. Stdin is either
// streamed straight into ffmpeg (pipe:0) or spooled to a temporary file when the whole
// input is needed (probing, two-pass filters) or the container cannot be read from a
// pipe (MP4/M4A keep their index at the end of the file).
type mediaInput struct {
	path   string    // readable file path (empty when streaming stdin)
	label  string    // name shown in messages
	stream io.Reader // set when ffmpeg reads from pipe:0
	sniff  string    // format guessed from the file extension or the first bytes
	stdin  bool
}

func (m *mediaInput) arg() string {
	if m.stream != nil {
		return "pipe:0"
	}
	return mediaPath(m.path)
}

// files lists the real input files, for the "output must not be the input" check.
func (m *mediaInput) files() []string {
	if m.stdin {
		return nil
	}
	return []string{m.path}
}

func isStdinName(input string) bool {
	return input == "" || input == "-" || input == "/dev/stdin"
}

// openInput resolves the input. needFile forces stdin to be spooled to a temporary file.
func openInput(input string, needFile bool) (*mediaInput, error) {
	if !isStdinName(input) {
		if err := requireInput(input); err != nil {
			return nil, err
		}
		sniff := extOf(input)
		if !isAudioFormat(sniff) {
			sniff = ""
		}
		return &mediaInput{path: input, label: input, sniff: sniff}, nil
	}

	info, err := os.Stdin.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice != 0 {
		return nil, errorf("no input: pass an input file or pipe audio into the command")
	}

	reader := bufio.NewReaderSize(os.Stdin, 64*1024)
	head, _ := reader.Peek(16)
	if len(head) == 0 {
		return nil, errorf("no input: stdin is empty")
	}

	in := &mediaInput{label: "stdin", sniff: sniffFormat(head), stdin: true}
	if !needFile && in.sniff != "m4a" {
		in.stream = reader
		return in, nil
	}

	path, err := tempFile(in.sniff)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, errorf("cannot write temporary file: %s", err.Error())
	}
	_, copyErr := io.Copy(file, reader)
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil {
		return nil, errorf("cannot read stdin into a temporary file")
	}
	in.path = path
	return in, nil
}

// sniffFormat recognizes common audio containers from their first bytes.
func sniffFormat(head []byte) string {
	has := func(offset int, magic string) bool {
		return len(head) >= offset+len(magic) && string(head[offset:offset+len(magic)]) == magic
	}
	switch {
	case has(0, "RIFF") && has(8, "WAVE"):
		return "wav"
	case has(0, "fLaC"):
		return "flac"
	case has(0, "OggS"):
		return "ogg"
	case has(0, "\x1a\x45\xdf\xa3"):
		return "webm"
	case has(4, "ftyp"):
		return "m4a"
	case has(0, "ID3"):
		return "mp3"
	case len(head) >= 2 && head[0] == 0xFF && head[1]&0xF6 == 0xF0:
		return "aac"
	case len(head) >= 2 && head[0] == 0xFF && head[1]&0xE0 == 0xE0:
		return "mp3"
	}
	return ""
}

// countingWriter counts the bytes streamed to stdout.
type countingWriter struct {
	w io.Writer
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += int64(n)
	return n, err
}

// mediaOutput is where a command writes audio: a file, or stdout when no --output is given.
type mediaOutput struct {
	path   string // empty => stdout
	format string // output format (extension); drives codec and, for stdout, the muxer
	stdout *countingWriter
}

func (o *mediaOutput) streaming() bool {
	return o.path == ""
}

func (o *mediaOutput) label() string {
	if o.streaming() {
		return "stdout"
	}
	return o.path
}

// args returns the ffmpeg output options and target.
func (o *mediaOutput) args() []string {
	if !o.streaming() {
		return []string{"-y", mediaPath(o.path)}
	}
	args := []string{"-f", muxerFor(o.format)}
	if o.format == "m4a" || o.format == "mp4" {
		// MP4 normally writes its index after the data; fragmenting makes it pipe-safe.
		args = append(args, "-movflags", "frag_keyframe+empty_moov")
	}
	return append(args, "pipe:1")
}

func (o *mediaOutput) writer() io.Writer {
	if !o.streaming() {
		return nil
	}
	if o.stdout == nil {
		o.stdout = &countingWriter{w: os.Stdout}
	}
	return o.stdout
}

// say prints a status message: on stdout when writing a file, on stderr when stdout
// carries the audio.
func (o *mediaOutput) say(format string, args ...interface{}) {
	target := os.Stdout
	if o.streaming() {
		target = os.Stderr
	}
	fmt.Fprintf(target, format+"\n", args...)
}

// discardPartial removes a file output that ffmpeg failed to finish.
func (o *mediaOutput) discardPartial() {
	if !o.streaming() {
		_ = os.Remove(o.path)
	}
}

// resolveOutput decides the output target and format.
//
//   - explicit --format wins, and must match the --output extension when both are given
//   - with --output, the format is the output extension
//   - when streaming, the format falls back to the input format, then to fallback;
//     an empty fallback means the format must be given explicitly
func resolveOutput(output, format string, in *mediaInput, fallback string, overwrite bool) (*mediaOutput, error) {
	format = strings.ToLower(format)
	if format != "" && !isAudioFormat(format) {
		return nil, errorf("unsupported format %q: use one of %s", format, strings.Join(audioFormats, ", "))
	}

	if output != "" {
		if format != "" && extOf(output) != format {
			return nil, errorf("output file %s does not match --format %s: use a .%s extension", output, format, format)
		}
		if err := checkOutput(output, overwrite, in.files()...); err != nil {
			return nil, err
		}
		return &mediaOutput{path: output, format: extOf(output)}, nil
	}

	if format == "" && fallback != "" {
		format = in.sniff
		if format == "" {
			format = fallback
		}
	}
	if format == "" {
		return nil, errorf("--format is required when writing to stdout (no --output given): use one of %s", strings.Join(audioFormats, ", "))
	}
	return &mediaOutput{format: format}, nil
}

// transcode runs one ffmpeg pass from in to out.
func transcode(in *mediaInput, out *mediaOutput, before, after []string) error {
	argv := append([]string{}, before...)
	argv = append(argv, "-i", in.arg())
	argv = append(argv, after...)
	argv = append(argv, out.args()...)
	if _, err := execFFmpegIO("error", argv, in.stream, out.writer()); err != nil {
		out.discardPartial()
		return err
	}
	return nil
}
