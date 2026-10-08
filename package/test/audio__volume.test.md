# audio volume

```beforeAll
mkdir -p tmp-volume
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 tmp-volume/tone.wav
```

```afterAll
rm -rf tmp-volume
```

## with a multiplier

### should scale the volume

```execute
aux4 audio volume tmp-volume/tone.wav --level 0.5 --output tmp-volume/tone-volume.wav
```

```expect
Adjusted volume of tmp-volume/tone.wav (0.5) -> tmp-volume/tone-volume.wav
```

### should keep the duration

```execute
aux4 audio info tmp-volume/tone-volume.wav
```

```expect:partial
  "duration": 2,
```

## with decibels

### should boost by a positive gain

```execute
aux4 audio volume tmp-volume/tone.wav --level 3dB --output tmp-volume/louder.wav
```

```expect
Adjusted volume of tmp-volume/tone.wav (3dB) -> tmp-volume/louder.wav
```

### should cut by a negative gain using the equals form

```execute
aux4 audio volume tmp-volume/tone.wav --level=-6dB --output tmp-volume/quieter.wav
```

```expect
Adjusted volume of tmp-volume/tone.wav (-6dB) -> tmp-volume/quieter.wav
```

## with an invalid level

### should reject it

```execute
aux4 audio volume tmp-volume/tone.wav --level loud --output tmp-volume/bad.wav
```

```error:partial
Error: invalid level "loud": use a multiplier (e.g. 0.5, 1.5) or decibels (e.g. 6dB, -3dB)
```

### should reject filter expressions

```execute
aux4 audio volume tmp-volume/tone.wav --level "1,atempo=2" --output tmp-volume/bad.wav
```

```error:partial
Error: invalid level "1,atempo=2": use a multiplier (e.g. 0.5, 1.5) or decibels (e.g. 6dB, -3dB)
```
