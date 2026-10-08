#### Description

The `extract` command saves the audio track of a video file as audio. Only the first audio track is extracted; video and subtitle tracks are dropped.

By default the audio is copied without re-encoding (no quality loss) into a container that fits its codec:

- **aac / alac** → `m4a`
- **mp3** → `mp3`
- **opus / vorbis** → `ogg`
- **flac** → `flac`
- **pcm** → `wav`
- anything else → `mka`

Give `--format` (or an `--output` with a different extension) to re-encode instead; `--bitrate` sets the bitrate when re-encoding. Without `--output`, the audio is streamed to stdout in that container. When no input file is given, the video is read from stdin. The command fails if the file has no audio track.

#### Usage

```bash
aux4 audio extract [<input>] [--format <format>] [--output <file>] [--bitrate <rate>] [--overwrite <true|false>]
```

<input>      Video file to extract the audio from (default: read from stdin)
--format     Re-encode to mp3, wav, flac, ogg, opus, m4a, aac or webm (default: copy the original audio)
--output     Output file path (default: write to stdout)
--bitrate    Audio bitrate when re-encoding, e.g. 128k
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio extract meeting.mp4 --output meeting.m4a
```

```text
Extracted audio from meeting.mp4 -> meeting.m4a
```

```bash
aux4 audio extract meeting.mp4 --format mp3 --bitrate 96k --output meeting.mp3
```

```text
Extracted audio from meeting.mp4 -> meeting.mp3
```
