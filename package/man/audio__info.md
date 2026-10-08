#### Description

The `info` command prints the properties of the first audio track in a file as JSON. It works on audio files and on video files that contain audio.

Fields:

- **file** — the path that was inspected
- **format** — container format (e.g. `wav`, `mp3`, `matroska,webm`, `mov,mp4,m4a,3gp,3g2,mj2`)
- **codec** — audio codec (e.g. `pcm_s16le`, `mp3`, `aac`, `opus`)
- **duration** — length in seconds, rounded to milliseconds
- **sampleRate** — samples per second in Hz
- **channels** — number of channels
- **channelLayout** — channel layout name (e.g. `mono`, `stereo`)
- **bitrate** — audio bitrate in bits per second
- **size** — file size in bytes

Recordings streamed by browsers (for example WebM/Opus files from `MediaRecorder`) often carry no duration in their header. In that case the audio is decoded once to measure the real duration.

The command fails with an error if the file does not exist or contains no audio track.

#### Usage

```bash
aux4 audio info <input>
```

<input>  Audio or video file to inspect

#### Example

```bash
aux4 audio info podcast.mp3
```

```json
{
  "file": "podcast.mp3",
  "format": "mp3",
  "codec": "mp3",
  "duration": 1834.512,
  "sampleRate": 44100,
  "channels": 2,
  "channelLayout": "stereo",
  "bitrate": 128000,
  "size": 29352448
}
```
