package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const installHint = "install it with 'brew install ffmpeg' (macOS), 'sudo apt install ffmpeg' (Debian/Ubuntu), 'sudo dnf install ffmpeg-free' (Fedora) or 'apk add ffmpeg' (Alpine)"

// tool resolves the path of ffmpeg or ffprobe. FFMPEG_PATH / FFPROBE_PATH override the PATH lookup.
func tool(name string) (string, error) {
	envName := strings.ToUpper(name) + "_PATH"
	if custom := os.Getenv(envName); custom != "" {
		if info, err := os.Stat(custom); err != nil || info.IsDir() {
			return "", errorf("%s not found at %s (set by %s); %s", name, custom, envName, installHint)
		}
		return custom, nil
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return "", errorf("%s is not installed or not in PATH; %s", name, installHint)
	}
	return path, nil
}

// mediaPath turns a local file path into an ffmpeg URL that can never be read as an
// option (leading "-") or as another protocol (http:, concat:, subfile:, ...).
func mediaPath(path string) string {
	return "file:" + path
}

// runFFmpeg executes ffmpeg with an argv array (never through a shell).
func runFFmpeg(args ...string) error {
	_, err := execFFmpegIO("error", args, nil, nil)
	return err
}

// runFFmpegReport executes ffmpeg at info log level and returns its stderr, where filters
// like loudnorm and silencedetect write their reports.
func runFFmpegReport(stdin io.Reader, args ...string) (string, error) {
	return execFFmpegIO("info", args, stdin, nil)
}

// execFFmpegIO runs ffmpeg with optional stdin (for pipe:0) and stdout (for pipe:1).
// Anything ffmpeg prints is captured and only surfaced in error messages, so stdout
// carries nothing but media bytes.
func execFFmpegIO(logLevel string, args []string, stdin io.Reader, stdout io.Writer) (string, error) {
	bin, err := tool("ffmpeg")
	if err != nil {
		return "", err
	}

	base := []string{"-hide_banner", "-nostats", "-v", logLevel}
	if stdin == nil {
		base = append(base, "-nostdin")
	}
	cmd := exec.Command(bin, append(base, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = &stderr
	}
	if stdin != nil {
		cmd.Stdin = stdin
	}

	if err := cmd.Run(); err != nil {
		return stderr.String(), fmt.Errorf("ffmpeg failed: %s", lastLines(stderr.String(), 3))
	}
	return stderr.String(), nil
}

func lastLines(text string, n int) string {
	lines := []string{}
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return "unknown error"
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "; ")
}

type probeStream struct {
	CodecType     string `json:"codec_type"`
	CodecName     string `json:"codec_name"`
	SampleRate    string `json:"sample_rate"`
	Channels      int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
	BitRate       string `json:"bit_rate"`
	Duration      string `json:"duration"`
}

type probeFormat struct {
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
}

type probeResult struct {
	Streams []probeStream `json:"streams"`
	Format  probeFormat   `json:"format"`
}

// AudioInfo is the normalized description of an audio file.
type AudioInfo struct {
	File          string  `json:"file"`
	Format        string  `json:"format"`
	Codec         string  `json:"codec"`
	Duration      float64 `json:"duration"`
	SampleRate    int     `json:"sampleRate"`
	Channels      int     `json:"channels"`
	ChannelLayout string  `json:"channelLayout"`
	Bitrate       int64   `json:"bitrate"`
	Size          int64   `json:"size"`
}

func probe(path string) (*AudioInfo, error) {
	bin, err := tool("ffprobe")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin, "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "a:0", mediaPath(path))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errorf("cannot read %s: %s", path, lastLines(stderr.String(), 2))
	}

	var result probeResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, errorf("cannot parse ffprobe output for %s: %s", path, err.Error())
	}

	if len(result.Streams) == 0 {
		return nil, errorf("no audio stream found in %s", path)
	}

	stream := result.Streams[0]
	info := &AudioInfo{
		File:          path,
		Format:        result.Format.FormatName,
		Codec:         stream.CodecName,
		SampleRate:    atoi(stream.SampleRate),
		Channels:      stream.Channels,
		ChannelLayout: stream.ChannelLayout,
		Size:          atoi64(result.Format.Size),
	}

	if info.ChannelLayout == "" {
		switch info.Channels {
		case 1:
			info.ChannelLayout = "mono"
		case 2:
			info.ChannelLayout = "stereo"
		}
	}

	info.Bitrate = atoi64(stream.BitRate)
	if info.Bitrate == 0 {
		info.Bitrate = atoi64(result.Format.BitRate)
	}

	duration := atof(result.Format.Duration)
	if duration == 0 {
		duration = atof(stream.Duration)
	}
	if duration == 0 {
		// Streamed recordings (e.g. browser MediaRecorder webm) carry no duration header:
		// decode the audio once to measure it.
		duration, err = decodeDuration(path)
		if err != nil {
			return nil, err
		}
	}
	info.Duration = math.Round(duration*1000) / 1000

	if info.Bitrate == 0 && info.Duration > 0 && info.Size > 0 {
		info.Bitrate = int64(math.Round(float64(info.Size*8) / info.Duration))
	}

	return info, nil
}

// decodeDuration measures duration by decoding the first audio stream to the null muxer.
func decodeDuration(path string) (float64, error) {
	bin, err := tool("ffmpeg")
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(bin, "-hide_banner", "-nostdin", "-v", "error", "-i", mediaPath(path), "-map", "0:a:0", "-f", "null", "-progress", "pipe:1", "-nostats", "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, errorf("cannot measure duration of %s: %s", path, lastLines(stderr.String(), 2))
	}

	var micros int64
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_us=") {
			if value, err := strconv.ParseInt(strings.TrimPrefix(line, "out_time_us="), 10, 64); err == nil {
				micros = value
			}
		}
	}
	return float64(micros) / 1e6, nil
}

func atoi(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}

func atoi64(s string) int64 {
	value, _ := strconv.ParseInt(s, 10, 64)
	return value
}

func atof(s string) float64 {
	value, _ := strconv.ParseFloat(s, 64)
	return value
}

// codecArgs returns the encoder arguments for an output extension.
func codecArgs(ext string) []string {
	switch ext {
	case "mp3":
		return []string{"-c:a", "libmp3lame"}
	case "wav":
		return []string{"-c:a", "pcm_s16le"}
	case "flac":
		return []string{"-c:a", "flac"}
	case "ogg":
		if hasEncoder("libvorbis") {
			return []string{"-c:a", "libvorbis"}
		}
		return []string{"-c:a", "libopus"}
	case "opus", "webm":
		return []string{"-c:a", "libopus"}
	case "m4a", "aac":
		return []string{"-c:a", "aac"}
	}
	return nil
}

func hasEncoder(name string) bool {
	bin, err := tool("ffmpeg")
	if err != nil {
		return false
	}
	out, err := exec.Command(bin, "-hide_banner", "-encoders").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == name {
			return true
		}
	}
	return false
}
