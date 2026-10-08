#### Description

The `convert` command re-encodes audio into another format. Supported formats:

- **mp3** — MP3 (LAME)
- **wav** — uncompressed 16-bit PCM
- **flac** — lossless FLAC
- **ogg** — Ogg Vorbis (Ogg Opus when the Vorbis encoder is not available)
- **opus** — Ogg Opus
- **m4a** — AAC in an MP4 container
- **aac** — raw AAC (ADTS)
- **webm** — Opus in a WebM container

With `--output`, the converted audio is saved to that file and the format is taken from `--format` or, if omitted, from the file extension (both must agree when given). Without `--output`, the converted audio is streamed to stdout and `--format` is required. When no input file is given, the audio is read from stdin.

Use `--bitrate`, `--sampleRate` and `--channels` to control the encoding; for example `--sampleRate 16000 --channels 1` produces mono 16 kHz audio. Only the first audio track is converted; video is dropped.

An existing output file is never replaced unless `--overwrite true` is given, and the output can never be the input file itself.

#### Usage

```bash
aux4 audio convert [<input>] [--format <format>] [--output <file>] [--bitrate <rate>] [--sampleRate <hz>] [--channels <n>] [--overwrite <true|false>]
```

<input>       Audio or video file to convert (default: read from stdin)
--format      Target format: mp3, wav, flac, ogg, opus, m4a, aac, webm (required when writing to stdout)
--output      Output file path (default: write to stdout)
--bitrate     Audio bitrate, e.g. 96k, 128k, 192k (default: encoder default)
--sampleRate  Sample rate in Hz, e.g. 16000, 44100, 48000 (default: same as input when supported)
--channels    Number of channels, 1 to 8 (default: same as input)
--overwrite   Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio convert interview.wav --format mp3 --bitrate 128k --output interview.mp3
```

```text
Converted interview.wav -> interview.mp3
```

```bash
aux4 audio convert recording.webm --format wav | aux4 whisper transcribe
```
