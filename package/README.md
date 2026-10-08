# aux4/audio

Audio toolkit powered by FFmpeg. Convert between formats, trim and split recordings, join files, compress recordings to fit upload limits, normalize loudness, change volume and speed, add fades, remove silence, extract audio from video, render waveform images and prepare recordings for speech recognition — all with short, predictable commands.

## Installation

```bash
aux4 aux4 pkger install aux4/audio
```

Requires `ffmpeg` and `ffprobe` (both ship with FFmpeg). The package tries to install FFmpeg automatically through your system package manager if it is not found. To install it yourself:

```bash
brew install ffmpeg            # macOS
sudo apt install ffmpeg        # Debian / Ubuntu
sudo dnf install ffmpeg-free   # Fedora
apk add ffmpeg                 # Alpine
```

## Quick Start

```bash
# Inspect a file
aux4 audio info podcast.mp3

# Convert to MP3 at 128 kbps
aux4 audio convert interview.wav --format mp3 --bitrate 128k --output interview.mp3

# Keep 1:30 to 2:45
aux4 audio trim episode.mp3 --start 01:30 --end 02:45 --output clip.mp3

# Pull the audio out of a video
aux4 audio extract meeting.mp4 --output meeting.m4a

# Shrink a recording to fit a 24 MB upload limit
aux4 audio compress meeting.wav --output meeting.mp3 --maxSize 24MB

# Turn a browser recording into the 16 kHz mono WAV that Whisper expects, and transcribe it
cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe
```

## Common Behavior

- **Input** — every command checks that the input file exists and contains an audio track before doing any work.
- **Output** — with `--output <file>`, the result is saved to that file and its format follows the file extension. Without `--output`, the result is streamed to stdout (see [Streaming and pipes](#streaming-and-pipes)).
- **No accidental overwrites** — an existing output file is never replaced unless you pass `--overwrite true`, and the output can never be the input file itself.
- **Negative values** — flags that take negative numbers must use an equals sign so the value is not read as another flag: `--target=-14`, `--threshold=-40`, `--level=-6dB`.
- **Errors** — problems are reported on stderr as `Error: ...` with a non-zero exit code.

## Streaming and pipes

The single-file commands — `info`, `convert`, `trim`, `extract`, `normalize`, `volume`, `speed`, `fade`, `trim-silence`, `compress`, `speech-prep` and `waveform` — read from stdin and write to stdout, so they chain with each other and with other tools:

```bash
# Browser recording straight into a transcriber
cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe

# File in, audio out
aux4 audio convert recording.webm --format wav | aux4 whisper transcribe

# Several steps without temporary files
cat raw.wav | aux4 audio trim-silence | aux4 audio normalize | aux4 audio compress --output episode.mp3
```

- **Reading stdin** — leave the input out, or pass `/dev/stdin`. MP4/M4A input (for example Safari recordings) is buffered to a temporary file first, because these containers cannot be read from a pipe. Commands that need the whole recording (`info`, `trim`, `extract`, `fade`, `normalize`, `trim-silence --all`, `compress --maxSize`) buffer piped input the same way. Temporary files are always removed.
- **Writing stdout** — happens whenever `--output` is not given. stdout then carries only the audio (or PNG for `waveform`); every message goes to stderr, and failures exit with a non-zero code.
- **Output format** — `--format` chooses it. `convert` requires it when streaming, `speech-prep` always writes WAV, `compress` defaults to MP3, `waveform` writes PNG, and the other commands keep the input format (WAV when it cannot be recognized).
- **Length headers** — WAV, MP3 and M4A written to a pipe cannot store their final length up front. WAV marks it as unknown, MP3 players estimate it, and M4A is written in fragments. FFmpeg-based tools read all of them correctly; use `--output` when you need a file with exact length metadata.

`concat` and `split` produce files only.

## Commands

### info

Show the properties of the first audio track as JSON: container format, codec, duration (seconds), sample rate, channels, channel layout, bitrate and file size. Works on video files and piped input too, and measures the real duration of streamed browser recordings that have no duration header.

```bash
aux4 audio info [<input>]
```

```bash
aux4 audio info podcast.mp3
```

```json
{
  "file": "podcast.mp3",
  "format": "mp3",
  "codec": "mp3",
  "duration": 1834.512,
  "sampleRate": 44100,
  "channels": 2,
  "channelLayout": "stereo",
  "bitrate": 128000,
  "size": 29352448
}
```

### convert

Convert audio to `mp3`, `wav`, `flac`, `ogg`, `opus`, `m4a`, `aac` or `webm`. The format comes from `--format`, or from the `--output` extension; `--format` is required when streaming to stdout.

```bash
aux4 audio convert [<input>] [--format <format>] [--output <file>] [--bitrate <rate>] [--sampleRate <hz>] [--channels <n>] [--overwrite <true|false>]
```

```bash
aux4 audio convert interview.wav --format mp3 --bitrate 128k --output interview.mp3
```

```text
Converted interview.wav -> interview.mp3
```

`--sampleRate` and `--channels` resample and remix, e.g. `--sampleRate 16000 --channels 1` for mono 16 kHz.

### trim

Keep a section of the audio, defined by `--start`, `--end` and/or `--duration` (`--end` and `--duration` are mutually exclusive). Times are seconds (`1.5`) or timecodes (`01:30`, `00:01:30.250`).

```bash
aux4 audio trim [<input>] [--start <time>] [--end <time>] [--duration <time>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio trim episode.mp3 --start 00:01:30 --end 00:02:45 --output clip.mp3
```

```text
Trimmed episode.mp3 (90s - 165s) -> clip.mp3
```

### extract

Save the audio track of a video file. By default the audio is copied without re-encoding into a matching container (AAC → `.m4a`, MP3 → `.mp3`, Opus/Vorbis → `.ogg`, ...). Pass `--format` to re-encode.

```bash
aux4 audio extract [<input>] [--format <format>] [--output <file>] [--bitrate <rate>] [--overwrite <true|false>]
```

```bash
aux4 audio extract meeting.mp4 --output meeting.m4a
```

```text
Extracted audio from meeting.mp4 -> meeting.m4a
```

### concat

Join two or more files in order. Inputs may differ in format, sample rate and channels; they are converted to match the first file.

```bash
aux4 audio concat <input> <input> [<input> ...] --output <file> [--overwrite <true|false>]
```

```bash
aux4 audio concat intro.wav interview.mp3 outro.wav --output episode.mp3
```

```text
Concatenated 3 files -> episode.mp3
```

### split

Cut a file into numbered parts, either every `--segment` seconds or in the middle of each pause with `--silence true` (pauses are quieter than `--threshold` dB, default `-35`, for at least `--minSilence` seconds, default `0.5`). Parts are written as `<prefix>-001.<ext>`, `<prefix>-002.<ext>`, ... in `--outputDir`, and their paths are printed one per line, in order.

```bash
aux4 audio split <input> (--segment <seconds> | --silence true) [--threshold <dB>] [--minSilence <seconds>] [--outputDir <dir>] [--prefix <name>] [--format <format>] [--overwrite <true|false>] [--json <true|false>]
```

```bash
aux4 audio split lecture.mp3 --segment 600 --outputDir chunks
```

```text
Split lecture.mp3 into 3 segments:
chunks/lecture-001.mp3
chunks/lecture-002.mp3
chunks/lecture-003.mp3
```

With `--json true`, the output is a JSON array with each part's `path`, its `start` offset and its `duration` in seconds, so timestamps computed on a part (such as a transcript) can be shifted back to the original timeline:

```bash
aux4 audio split lecture.mp3 --segment 600 --outputDir chunks --json true
```

```json
[
  {
    "path": "chunks/lecture-001.mp3",
    "start": 0,
    "duration": 600
  },
  {
    "path": "chunks/lecture-002.mp3",
    "start": 600,
    "duration": 600
  },
  {
    "path": "chunks/lecture-003.mp3",
    "start": 1200,
    "duration": 287.412
  }
]
```

### normalize

Bring a recording to a target loudness (EBU R128, two passes, linear gain). Defaults suit podcasts and spoken word: `--target -16` LUFS, `--truePeak -1.5` dBTP, `--loudnessRange 11` LU. Use `--target=-14` for music streaming or `--target=-23` for broadcast.

```bash
aux4 audio normalize [<input>] [--target <LUFS>] [--truePeak <dBTP>] [--loudnessRange <LU>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio normalize episode.wav --output episode-normalized.wav
```

```text
Normalized episode.wav (-16 LUFS) -> episode-normalized.wav
```

### volume

Change the volume by a multiplier (`0.5`, `2`) or a gain in decibels (`6dB`, `--level=-3dB`).

```bash
aux4 audio volume [<input>] --level <level> [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio volume background.mp3 --level 0.3 --output background-quiet.mp3
```

```text
Adjusted volume of background.mp3 (0.3) -> background-quiet.mp3
```

### speed

Play faster or slower without changing the pitch. `--factor` ranges from `0.1` to `10`.

```bash
aux4 audio speed [<input>] --factor <factor> [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio speed lecture.mp3 --factor 1.5 --output lecture-fast.mp3
```

```text
Changed speed of lecture.mp3 (1.5x) -> lecture-fast.mp3
```

### fade

Add a fade-in and/or fade-out, lengths in seconds.

```bash
aux4 audio fade [<input>] [--fadeIn <seconds>] [--fadeOut <seconds>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio fade jingle.wav --fadeIn 0.5 --fadeOut 2 --output jingle-fade.wav
```

```text
Faded jingle.wav (in: 0.5s, out: 2s) -> jingle-fade.wav
```

### trim-silence

Remove silence at the beginning and end (audio quieter than `--threshold` dB, default `-50`). With `--all true`, pauses inside the recording longer than `--maxPause` seconds (default `0.5`) are shortened to that length.

```bash
aux4 audio trim-silence [<input>] [--threshold <dB>] [--all <true|false>] [--maxPause <seconds>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

```bash
aux4 audio trim-silence voice-note.m4a --output voice-note-clean.m4a
```

```text
Removed silence from voice-note.m4a -> voice-note-clean.m4a
```

### compress

Make a recording as small as possible while keeping speech clear: mono, 16 kHz, 32 kbps MP3 by default, or `--format opus` / `--format m4a`. With `--maxSize` (bytes, or `KB`/`MB`/`GB` decimal and `KiB`/`MiB`/`GiB` binary), the bitrate is lowered from the input duration so the result fits, down to 8 kbps. If the audio cannot fit even at 8 kbps, nothing is written and the command fails asking you to `split` the file first. The output path, size in bytes and bitrate used are printed.

```bash
aux4 audio compress [<input>] [--output <file>] [--format <mp3|opus|m4a>] [--bitrate <rate>] [--sampleRate <hz>] [--channels <n>] [--maxSize <size>] [--overwrite <true|false>]
```

```bash
aux4 audio compress all-hands.m4a --output upload.mp3 --maxSize 24MB
```

```text
Compressed all-hands.m4a -> upload.mp3 (21844992 bytes, 24k)
```

### speech-prep

Convert any recording — including browser recordings (WebM/Opus from Chrome and Firefox, MP4/AAC from Safari) — into 16 kHz mono 16-bit WAV, the input format of Whisper, whisper.cpp and most speech-recognition engines. `--sampleRate` changes the rate when an engine needs something else.

```bash
aux4 audio speech-prep [<input>] [--sampleRate <hz>] [--output <file.wav>] [--overwrite <true|false>]
```

```bash
aux4 audio speech-prep recording.webm --output recording-16k.wav
```

```text
Prepared recording.webm for speech recognition (16000 Hz mono WAV) -> recording-16k.wav
```

```bash
cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe
```

### waveform

Render a waveform image (`.png` with a transparent background, or `.jpg`). Without `--output`, a PNG is streamed to stdout.

```bash
aux4 audio waveform [<input>] [--size <WxH>] [--color <color>] [--output <file>] [--overwrite <true|false>]
```

```bash
aux4 audio waveform episode.mp3 --size 800x120 --color "#111827" --output cover-wave.png
```

```text
Rendered waveform of episode.mp3 (800x120) -> cover-wave.png
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `FFMPEG_PATH` | Path to the `ffmpeg` binary to use instead of the one on `PATH` |
| `FFPROBE_PATH` | Path to the `ffprobe` binary to use instead of the one on `PATH` |
