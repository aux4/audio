#### Description

The `audio` command provides a suite of audio tools powered by FFmpeg. It converts between formats, trims and splits recordings, joins files, compresses recordings to fit size limits, normalizes loudness, changes volume and speed, applies fades, removes silence, extracts audio from video, renders waveform images and prepares recordings for speech recognition.

Requires `ffmpeg` and `ffprobe` to be installed.

**Streaming and pipes.** The single-file commands (`info`, `convert`, `trim`, `extract`, `normalize`, `volume`, `speed`, `fade`, `trim-silence`, `compress`, `speech-prep`, `waveform`) work in pipelines:

- **Input** — give a file, or leave the input out (or pass `/dev/stdin`) to read audio piped into the command.
- **Output** — with `--output <file>` the result is saved to that file; without `--output` the result is streamed to stdout. In that case stdout carries only the audio bytes and every message goes to stderr.
- **Format** — when streaming, the output format comes from `--format`; `convert` requires it, `speech-prep` always writes WAV, `compress` defaults to MP3, `waveform` writes PNG, and the other commands keep the input format (WAV when it cannot be recognized).

`concat` and `split` work with files only.

Every command validates its input, refuses to replace an existing file unless `--overwrite true` is given, and exits with a non-zero code and an `Error: ...` message on stderr when something fails.

Flags that take negative numbers (such as `--target`, `--threshold` or a `-6dB` volume level) must be written with an equals sign, for example `--target=-14`, so the value is not mistaken for another flag.

#### Usage

```bash
aux4 audio <command>
```

Available commands:

- **info** — Show audio metadata as JSON
- **convert** — Convert audio to another format
- **trim** — Cut a section of audio by start, end or duration
- **extract** — Extract the audio track from a video file
- **concat** — Join multiple audio files into one
- **split** — Split audio into parts by fixed duration or at silences
- **normalize** — Normalize loudness to a target level
- **volume** — Change the volume by a multiplier or in decibels
- **speed** — Change playback speed without changing pitch
- **fade** — Apply a fade-in and/or fade-out
- **trim-silence** — Remove leading and trailing silence
- **compress** — Shrink audio for upload or speech recognition, optionally under a maximum size
- **speech-prep** — Prepare audio for speech recognition (16 kHz mono WAV)
- **waveform** — Render a waveform image of the audio

#### Example

```bash
aux4 audio convert interview.wav --format mp3 --bitrate 128k --output interview.mp3
```

```text
Converted interview.wav -> interview.mp3
```

```bash
cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe
```
