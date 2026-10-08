#### Description

The `compress` command makes a recording as small as possible while keeping speech clear, for uploads and speech-recognition APIs that limit file size. The defaults suit spoken word: mono, 16 kHz, 32 kbps MP3. Use `--format opus` for the smallest files at a given quality, or `--format m4a` for the widest playback support.

With `--maxSize`, the bitrate is calculated from the duration of the input so the result fits under the limit:

- the requested `--bitrate` is kept if the file already fits
- otherwise the bitrate is lowered, down to a minimum of 8 kbps
- if the encoder overshoots, the file is re-encoded at a lower bitrate
- if the audio cannot fit even at 8 kbps, the command fails without writing anything and tells you to split the file first (see `aux4 audio split`)

`--maxSize` accepts bytes (`25000000`) or a size with a unit. `KB`, `MB` and `GB` are decimal (1 MB = 1,000,000 bytes); `KiB`, `MiB` and `GiB` are binary (1 MiB = 1,048,576 bytes).

The command prints the output path, its size in bytes and the bitrate used. Video tracks and metadata are dropped. Opus files always report a 48 kHz playback rate in `aux4 audio info`, whatever rate they were encoded at.

#### Usage

```bash
aux4 audio compress <input> --output <file> [--format <mp3|opus|m4a>] [--bitrate <rate>] [--sampleRate <hz>] [--channels <n>] [--maxSize <size>] [--overwrite <true|false>]
```

<input>       Audio or video file to compress
--output      Output file path; the extension must match --format (default: <input>-compressed.<format>)
--format      Output format: mp3, opus or m4a (default: mp3)
--bitrate     Audio bitrate, at least 8k (default: 32k)
--sampleRate  Sample rate in Hz (default: 16000); opus accepts 8000, 12000, 16000, 24000 or 48000
--channels    Number of channels (default: 1)
--maxSize     Maximum output size, e.g. 24MB, 500KB, 25000000 (default: no limit)
--overwrite   Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio compress meeting.wav --output meeting.mp3
```

```text
Compressed meeting.wav -> meeting.mp3 (14402347 bytes, 32k)
```

```bash
aux4 audio compress all-hands.m4a --output upload.mp3 --maxSize 24MB
```

```text
Compressed all-hands.m4a -> upload.mp3 (21844992 bytes, 24k)
```

```bash
aux4 audio compress conference.wav --output upload.mp3 --maxSize 1MB
```

```text
Error: conference.wav (7200s) cannot fit in 1MB even at the minimum bitrate of 8k: split it first (e.g. aux4 audio split --segment <seconds>)
```
