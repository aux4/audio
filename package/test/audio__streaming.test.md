# audio streaming and pipes

```beforeAll
mkdir -p tmp-streaming
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=3 -ac 2 tmp-streaming/tone.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:duration=2 -ac 2 -c:a libopus -f webm pipe:1 > tmp-streaming/chrome.webm
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 -ac 2 -c:a aac tmp-streaming/safari.mp4
ffmpeg -hide_banner -loglevel error -y -f lavfi -i "aevalsrc=if(between(t\,1\,2)\,0.5*sin(2*PI*440*t)\,0):s=44100:d=3" tmp-streaming/padded.wav
```

```afterAll
rm -rf tmp-streaming
```

## pipe in and pipe out

### should prepare a piped browser recording for speech recognition

```execute
cat tmp-streaming/chrome.webm | aux4 audio speech-prep | aux4 audio info
```

```expect:partial
{
  "file": "stdin",
  "format": "wav",
  "codec": "pcm_s16le",
  "duration": 2,
  "sampleRate": 16000,
  "channels": 1,
```

```error
Prepared stdin for speech recognition (16000 Hz mono WAV) -> stdout
```

### should chain several commands

```execute
cat tmp-streaming/tone.wav | aux4 audio trim --start 1 | aux4 audio volume --level 0.5 --format mp3 | aux4 audio info
```

```expect:partial
  "codec": "mp3",
```

## file in and pipe out

### should write only audio bytes to stdout

```execute
aux4 audio convert tmp-streaming/tone.wav --format wav > tmp-streaming/piped.wav && head -c 4 tmp-streaming/piped.wav
```

```expect
RIFF
```

```error
Converted tmp-streaming/tone.wav -> stdout
```

### should stream as many bytes as a file conversion writes

```execute
aux4 audio convert tmp-streaming/tone.wav --output tmp-streaming/direct.wav && test "$(aux4 audio convert tmp-streaming/tone.wav --format wav | wc -c)" -eq "$(wc -c < tmp-streaming/direct.wav)" && echo "same size"
```

```expect:partial
same size
```

### should pass the bytes through unchanged

```execute
test "$(aux4 audio convert tmp-streaming/tone.wav --format wav | shasum)" = "$(ffmpeg -hide_banner -nostdin -v error -i file:tmp-streaming/tone.wav -map 0:a:0 -vn -c:a pcm_s16le -f wav pipe:1 | shasum)" && echo identical
```

```expect
identical
```

### should carry the same audio as a file conversion

```execute
test "$(aux4 audio convert tmp-streaming/tone.wav --format wav | ffmpeg -v error -i pipe:0 -f md5 -)" = "$(ffmpeg -v error -i tmp-streaming/direct.wav -f md5 -)" && echo "same audio"
```

```expect
same audio
```

### should stream m4a in a pipe-safe layout

```execute
aux4 audio convert tmp-streaming/tone.wav --format m4a | aux4 audio info
```

```expect:partial
  "codec": "aac",
```

### should keep the input format when no format is given

```execute
aux4 audio volume tmp-streaming/tone.wav --level 2 | aux4 audio info
```

```expect:partial
  "format": "wav",
```

### should stream a waveform PNG

```execute
aux4 audio waveform tmp-streaming/tone.wav > tmp-streaming/wave.png && head -c 4 tmp-streaming/wave.png | tail -c 3
```

```expect
PNG
```

## stdin in and file out

### should read the input from stdin

```execute
cat tmp-streaming/tone.wav | aux4 audio convert --output tmp-streaming/from-stdin.flac && aux4 audio info tmp-streaming/from-stdin.flac
```

```expect:partial
Converted stdin -> tmp-streaming/from-stdin.flac
**
  "codec": "flac",
  "duration": 3,
```

### should accept /dev/stdin as the input

```execute
cat tmp-streaming/tone.wav | aux4 audio speed /dev/stdin --factor 2 --output tmp-streaming/fast.wav
```

```expect
Changed speed of stdin (2x) -> tmp-streaming/fast.wav
```

## containers that cannot be read from a pipe

### should spool a piped mp4 recording before converting it

```execute
cat tmp-streaming/safari.mp4 | aux4 audio speech-prep --output tmp-streaming/safari.wav && aux4 audio info tmp-streaming/safari.wav
```

```expect:partial
Prepared stdin for speech recognition (16000 Hz mono WAV) -> tmp-streaming/safari.wav
**
  "sampleRate": 16000,
  "channels": 1,
```

### should extract the audio of a piped video

```execute
cat tmp-streaming/safari.mp4 | aux4 audio extract | aux4 audio info
```

```expect:partial
  "codec": "aac",
```

## commands that need the whole input

### should normalize piped audio

```execute
cat tmp-streaming/tone.wav | aux4 audio normalize | aux4 audio info
```

```expect:partial
  "format": "wav",
  "codec": "pcm_s16le",
  "duration": 3,
```

### should remove silence from piped audio

```execute
cat tmp-streaming/padded.wav | aux4 audio trim-silence --all true | aux4 audio info
```

```expect:regex
"duration": (1|0\.9\d*|1\.0\d*),
```

### should describe piped audio

```execute
cat tmp-streaming/safari.mp4 | aux4 audio info
```

```expect:partial
  "file": "stdin",
  "format": "mov,mp4,m4a,3gp,3g2,mj2",
  "codec": "aac",
```

## errors

### should refuse to stream without a format it can infer

```execute
cat tmp-streaming/tone.wav | aux4 audio convert
```

```error
Error: --format is required when writing to stdout (no --output given): use one of mp3, wav, flac, ogg, opus, m4a, aac, webm
```

### should fail when there is no input file and nothing on stdin

```execute
aux4 audio volume --level 0.5 < /dev/null
```

```error
Error: no input: pass an input file or pipe audio into the command
```

### should report invalid piped data without printing audio

```execute
printf 'not audio at all' | aux4 audio convert --format wav | wc -c
```

```expect:partial
0
```

```error:partial
Error: ffmpeg failed:
```
