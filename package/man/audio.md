#### Description

The `audio` command provides a suite of audio tools powered by FFmpeg. It converts between formats, trims and splits recordings, joins files, normalizes loudness, changes volume and speed, applies fades, removes silence, extracts audio from video, renders waveform images and prepares recordings for speech recognition.

Every command validates its input file, derives a sensible output file name when `--output` is omitted, and refuses to replace an existing file unless `--overwrite true` is given. Requires `ffmpeg` and `ffprobe` to be installed.

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
- **speech-prep** — Prepare audio for speech recognition (16 kHz mono WAV)
- **waveform** — Render a waveform image of the audio

#### Example

```bash
aux4 audio convert interview.wav --format mp3 --bitrate 128k
```

```text
Converted interview.wav -> interview.mp3
```
