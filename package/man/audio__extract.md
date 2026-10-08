#### Description

The `extract` command saves the audio track of a video file as an audio file. Only the first audio track is extracted; video and subtitle tracks are dropped.

By default the audio is copied without re-encoding (no quality loss) into a container that fits its codec:

- **aac / alac** → `.m4a`
- **mp3** → `.mp3`
- **opus / vorbis** → `.ogg`
- **flac** → `.flac`
- **pcm** → `.wav`
- anything else → `.mka`

Give `--format` (or an `--output` with a different extension) to re-encode instead; `--bitrate` sets the bitrate when re-encoding. The command fails if the file has no audio track.

#### Usage

```bash
aux4 audio extract <input> [--format <format>] [--output <file>] [--bitrate <rate>] [--overwrite <true|false>]
```

<input>      Video file to extract the audio from
--format     Re-encode to mp3, wav, flac, ogg, opus, m4a, aac or webm (default: copy the original audio)
--output     Output file path (default: input name with the audio extension)
--bitrate    Audio bitrate when re-encoding, e.g. 128k
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio extract meeting.mp4
```

```text
Extracted audio from meeting.mp4 -> meeting.m4a
```

```bash
aux4 audio extract meeting.mp4 --format mp3 --bitrate 96k
```

```text
Extracted audio from meeting.mp4 -> meeting.mp3
```
