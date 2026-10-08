# audio speed

```beforeAll
mkdir -p tmp-speed
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=3 tmp-speed/tone.wav
```

```afterAll
rm -rf tmp-speed
```

## faster

### should speed up the audio

```execute
aux4 audio speed tmp-speed/tone.wav --factor 2
```

```expect
Changed speed of tmp-speed/tone.wav (2x) -> tmp-speed/tone-speed.wav
```

### should halve the duration

```execute
aux4 audio info tmp-speed/tone-speed.wav
```

```expect:regex
"duration": 1\.[45]\d*,
```

## slower

### should slow down beyond the single-stage range

```execute
aux4 audio speed tmp-speed/tone.wav --factor 0.25 --output tmp-speed/slow.wav && aux4 audio info tmp-speed/slow.wav
```

```expect:regex
Changed speed of tmp-speed/tone\.wav \(0\.25x\) -> tmp-speed/slow\.wav[\s\S]*"duration": 1[12]\.\d+,
```

## with an invalid factor

### should reject a factor out of range

```execute
aux4 audio speed tmp-speed/tone.wav --factor 20 --output tmp-speed/bad.wav
```

```error:partial
Error: invalid factor "20": must be between 0.1 and 10
```

### should reject a non-numeric factor

```execute
aux4 audio speed tmp-speed/tone.wav --factor fast --output tmp-speed/bad.wav
```

```error:partial
Error: invalid factor "fast": must be a positive number
```
