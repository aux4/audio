# audio trim-silence

```beforeAll
mkdir -p tmp-trim-silence
ffmpeg -hide_banner -loglevel error -y -f lavfi -i "aevalsrc=if(between(t\,1\,2)\,0.5*sin(2*PI*440*t)\,0):s=44100:d=3" tmp-trim-silence/padded.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i "aevalsrc=if(between(t\,1\,2)\,0\,0.5*sin(2*PI*440*t)):s=44100:d=3" tmp-trim-silence/pause.wav
```

```afterAll
rm -rf tmp-trim-silence
```

## leading and trailing silence

### should remove silence at both ends

```execute
aux4 audio trim-silence tmp-trim-silence/padded.wav
```

```expect
Removed silence from tmp-trim-silence/padded.wav -> tmp-trim-silence/padded-nosilence.wav
```

### should keep only the sound

```execute
aux4 audio info tmp-trim-silence/padded-nosilence.wav
```

```expect:regex
"duration": (1|0\.9\d*|1\.0\d*),
```

## pauses inside the audio

### should leave inner pauses alone by default

```execute
aux4 audio trim-silence tmp-trim-silence/pause.wav --output tmp-trim-silence/kept.wav && aux4 audio info tmp-trim-silence/kept.wav
```

```expect:regex
"duration": (3|2\.9\d*),
```

### should shorten long pauses with --all true

```execute
aux4 audio trim-silence tmp-trim-silence/pause.wav --all true --maxPause 0.2 --output tmp-trim-silence/tight.wav && aux4 audio info tmp-trim-silence/tight.wav
```

```expect:regex
"duration": 2\.[1-3]\d*,
```

## with an invalid threshold

### should reject a positive level

```execute
aux4 audio trim-silence tmp-trim-silence/padded.wav --threshold 10 --output tmp-trim-silence/bad.wav
```

```error:partial
Error: invalid threshold "10": must be a level in dB at or below 0 (e.g. -50)
```
