# Release notes

## New package

First release of `aux4/audio`, an audio toolkit powered by FFmpeg.

- `info` — audio metadata as JSON (format, codec, duration, sample rate, channels, bitrate, size); measures the real duration of streamed browser recordings that have no duration header.
- `convert` — mp3, wav, flac, ogg, opus, m4a, aac and webm, with `--bitrate`, `--sampleRate` and `--channels`.
- `trim` — keep a section by `--start`, `--end` or `--duration` (seconds or timecodes).
- `extract` — audio track from a video, copied without re-encoding by default.
- `concat` — join files of different formats, sample rates and channel counts.
- `split` — fixed-length parts or cuts in the middle of each pause.
- `normalize` — two-pass EBU R128 loudness normalization.
- `volume`, `speed` (pitch preserved), `fade` and `trim-silence`.
- `speech-prep` — 16 kHz mono 16-bit WAV for Whisper and other speech recognizers, from any recording including Chrome/Firefox WebM and Safari MP4.
- `waveform` — PNG or JPG waveform image.

Every command validates its input, never replaces an existing file unless `--overwrite true` is given, and explains how to install FFmpeg when it is missing.
