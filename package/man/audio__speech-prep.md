#### Description

The `speech-prep` command converts any recording into the format speech-recognition engines such as Whisper and whisper.cpp expect: 16 kHz, mono, 16-bit PCM WAV. It accepts audio files and video files, including browser recordings:

- **Chrome / Firefox** — WebM or Ogg with Opus audio, often without a duration header
- **Safari** — MP4 with AAC audio

Only the first audio track is used; video is dropped and stereo is mixed down to mono. Use `--sampleRate` when an engine needs a different rate (for example 8000 Hz for telephony models).

With `--output`, the WAV is saved to that file, which must end in `.wav`. Without it, the WAV is streamed to stdout, ready to pipe into a transcriber. When no input file is given, the recording is read from stdin; MP4 recordings are buffered to a temporary file first because they cannot be read from a pipe. A WAV written to a pipe cannot record its own length in the header, so it marks the length as unknown; readers that understand streamed WAV (FFmpeg and most audio tools) read it to the end.

#### Usage

```bash
aux4 audio speech-prep [<input>] [--sampleRate <hz>] [--output <file.wav>] [--overwrite <true|false>]
```

<input>       Audio or video file to prepare (default: read from stdin)
--sampleRate  Sample rate in Hz (default: 16000)
--output      Output WAV file path (default: write to stdout)
--overwrite   Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio speech-prep recording.webm --output recording-16k.wav
```

```text
Prepared recording.webm for speech recognition (16000 Hz mono WAV) -> recording-16k.wav
```

```bash
cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe
```
