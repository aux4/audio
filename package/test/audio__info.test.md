# audio info

```beforeAll
mkdir -p tmp-info
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 -ac 2 -fflags +bitexact tmp-info/tone.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:duration=2 -c:a libopus -f webm pipe:1 > tmp-info/recording.webm
ffmpeg -hide_banner -loglevel error -y -f lavfi -i testsrc=size=64x64:rate=10:duration=2 -f lavfi -i sine=frequency=440:duration=2 -c:v mpeg4 -c:a aac -shortest tmp-info/video.mp4
ffmpeg -hide_banner -loglevel error -y -f lavfi -i testsrc=size=64x64:rate=10:duration=1 -c:v mpeg4 tmp-info/silent-video.mp4
```

```afterAll
rm -rf tmp-info
```

## with a wav file

### should print the audio metadata as JSON

```execute
aux4 audio info tmp-info/tone.wav
```

```expect:partial
{
  "file": "tmp-info/tone.wav",
  "format": "wav",
  "codec": "pcm_s16le",
  "duration": 2,
  "sampleRate": 44100,
  "channels": 2,
  "channelLayout": "stereo",
  "bitrate": 1411200,
  "size": *
}
```

## with a streamed browser recording

### should measure the duration even when the file has no duration header

```execute
aux4 audio info tmp-info/recording.webm
```

```expect:partial
  "format": "matroska,webm",
  "codec": "opus",
  "duration": 2,
  "sampleRate": 48000,
  "channels": 1,
```

## with a video file

### should describe the audio track

```execute
aux4 audio info tmp-info/video.mp4
```

```expect:partial
  "codec": "aac",
  "duration": 2,
  "sampleRate": 44100,
  "channels": 1,
```

## with a video without audio

### should fail with a clear error

```execute
aux4 audio info tmp-info/silent-video.mp4
```

```error:partial
Error: no audio stream found in tmp-info/silent-video.mp4
```

## with a missing file

### should fail with a clear error

```execute
aux4 audio info tmp-info/missing.wav
```

```error:partial
Error: input file not found: tmp-info/missing.wav
```

## when ffprobe cannot be found

### should explain how to install it

```execute
FFPROBE_PATH=/nonexistent/ffprobe aux4 audio info tmp-info/tone.wav
```

```error:partial
Error: ffprobe not found at /nonexistent/ffprobe (set by FFPROBE_PATH); install it with 'brew install ffmpeg' (macOS), 'sudo apt install ffmpeg' (Debian/Ubuntu)*
```
