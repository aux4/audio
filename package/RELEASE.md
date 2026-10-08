# Release notes

## New package

First release of `aux4/audio`, an audio toolkit powered by FFmpeg.

- `info` — audio metadata as JSON (format, codec, duration, sample rate, channels, bitrate, size); measures the real duration of streamed browser recordings that have no duration header.
- `convert` — mp3, wav, flac, ogg, opus, m4a, aac and webm, with `--bitrate`, `--sampleRate` and `--channels`.
- `trim` — keep a section by `--start`, `--end` or `--duration` (seconds or timecodes).
- `extract` — audio track from a video, copied without re-encoding by default.
- `concat` — join files of different formats, sample rates and channel counts.
- `split` — fixed-length parts or cuts in the middle of each pause; `--json true` lists each part with its start offset and duration, so per-part timestamps can be shifted back to the original timeline.
- `normalize` — two-pass EBU R128 loudness normalization.
- `volume`, `speed` (pitch preserved), `fade` and `trim-silence`.
- `compress` — speech-friendly mp3, opus or m4a (mono, 16 kHz, 32 kbps by default); `--maxSize` picks the bitrate from the duration so the file fits an upload limit, and fails with a clear "split it first" message when it cannot.
- `speech-prep` — 16 kHz mono 16-bit WAV for Whisper and other speech recognizers, from any recording including Chrome/Firefox WebM and Safari MP4.
- `waveform` — PNG or JPG waveform image.

Streaming and pipes: the single-file commands read audio from stdin when no input file is given, and stream the result to stdout when `--output` is omitted, so `cat recording.webm | aux4 audio speech-prep | aux4 whisper transcribe` works. stdout then carries only audio bytes and all messages go to stderr. MP4/M4A input and commands that need the whole recording buffer stdin to a temporary file that is always removed.

Every command validates its input, never replaces an existing file unless `--overwrite true` is given, and explains how to install FFmpeg when it is missing.
