#### Description

The `speech-prep` command converts any recording into the format speech-recognition engines such as Whisper and whisper.cpp expect: 16 kHz, mono, 16-bit PCM WAV. It accepts audio files and video files, including browser recordings:

- **Chrome / Firefox** — WebM or Ogg with Opus audio, often without a duration header
- **Safari** — MP4 with AAC audio

Only the first audio track is used; video is dropped and stereo is mixed down to mono. Use `--sampleRate` when an engine needs a different rate (for example 8000 Hz for telephony models). The output must be a `.wav` file.

#### Usage

```bash
aux4 audio speech-prep <input> [--sampleRate <hz>] [--output <file.wav>] [--overwrite <true|false>]
```

<input>       Audio or video file to prepare
--sampleRate  Sample rate in Hz (default: 16000)
--output      Output WAV file path (default: <input>-16k.wav)
--overwrite   Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio speech-prep recording.webm
```

```text
Prepared recording.webm for speech recognition (16000 Hz mono WAV) -> recording-16k.wav
```
