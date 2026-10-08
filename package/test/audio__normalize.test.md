# audio normalize

```beforeAll
mkdir -p tmp-normalize
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=3 -af volume=0.1 tmp-normalize/quiet.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i anullsrc=r=44100:cl=mono -t 1 tmp-normalize/silence.wav
```

```afterAll
rm -rf tmp-normalize
```

## with default target

### should normalize to -16 LUFS

```execute
aux4 audio normalize tmp-normalize/quiet.wav
```

```expect
Normalized tmp-normalize/quiet.wav (-16 LUFS) -> tmp-normalize/quiet-normalized.wav
```

### should keep the original sample rate and duration

```execute
aux4 audio info tmp-normalize/quiet-normalized.wav
```

```expect:partial
  "duration": 3,
  "sampleRate": 44100,
  "channels": 1,
```

## with a custom target

### should accept negative values with the equals form

```execute
aux4 audio normalize tmp-normalize/quiet.wav --target=-23 --truePeak=-2 --output tmp-normalize/broadcast.wav
```

```expect
Normalized tmp-normalize/quiet.wav (-23 LUFS) -> tmp-normalize/broadcast.wav
```

## with invalid options

### should reject a target out of range

```execute
aux4 audio normalize tmp-normalize/quiet.wav --target=-100 --output tmp-normalize/bad.wav
```

```error:partial
Error: invalid target "-100": must be between -70 and -5 LUFS
```

### should refuse to normalize silence

```execute
aux4 audio normalize tmp-normalize/silence.wav
```

```error:partial
Error: cannot normalize tmp-normalize/silence.wav: the audio is silent or too short to measure
```
