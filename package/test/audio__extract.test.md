# audio extract

```beforeAll
mkdir -p tmp-extract
ffmpeg -hide_banner -loglevel error -y -f lavfi -i testsrc=size=64x64:rate=10:duration=2 -f lavfi -i sine=frequency=440:duration=2 -c:v mpeg4 -c:a aac -shortest tmp-extract/video.mp4
ffmpeg -hide_banner -loglevel error -y -f lavfi -i testsrc=size=64x64:rate=10:duration=1 -c:v mpeg4 tmp-extract/silent.mp4
```

```afterAll
rm -rf tmp-extract
```

## with an output file in the original codec's container

### should copy the audio track without re-encoding

```execute
aux4 audio extract tmp-extract/video.mp4 --output tmp-extract/video.m4a
```

```expect
Extracted audio from tmp-extract/video.mp4 -> tmp-extract/video.m4a
```

### should keep the original codec

```execute
aux4 audio info tmp-extract/video.m4a
```

```expect:partial
  "codec": "aac",
```

## with a format

### should re-encode to the requested format

```execute
aux4 audio extract tmp-extract/video.mp4 --format mp3 --bitrate 96k --output tmp-extract/video.mp3 && aux4 audio info tmp-extract/video.mp3
```

```expect:partial
Extracted audio from tmp-extract/video.mp4 -> tmp-extract/video.mp3
**
  "codec": "mp3",
**
  "bitrate": 96000,
```

## with an output file

### should use the output extension as the format

```execute
aux4 audio extract tmp-extract/video.mp4 --output tmp-extract/track.wav && aux4 audio info tmp-extract/track.wav
```

```expect:partial
Extracted audio from tmp-extract/video.mp4 -> tmp-extract/track.wav
**
  "codec": "pcm_s16le",
```

## with a video without audio

### should fail with a clear error

```execute
aux4 audio extract tmp-extract/silent.mp4
```

```error:partial
Error: no audio stream found in tmp-extract/silent.mp4
```
