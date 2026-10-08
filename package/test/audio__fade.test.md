# audio fade

```beforeAll
mkdir -p tmp-fade
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=3 tmp-fade/tone.wav
```

```afterAll
rm -rf tmp-fade
```

## with fade-in and fade-out

### should apply both fades

```execute
aux4 audio fade tmp-fade/tone.wav --fadeIn 1 --fadeOut 0.5
```

```expect
Faded tmp-fade/tone.wav (in: 1s, out: 0.5s) -> tmp-fade/tone-fade.wav
```

### should keep the duration

```execute
aux4 audio info tmp-fade/tone-fade.wav
```

```expect:partial
  "duration": 3,
```

## with only a fade-out

### should leave the beginning untouched

```execute
aux4 audio fade tmp-fade/tone.wav --fadeOut 2 --output tmp-fade/out.wav
```

```expect
Faded tmp-fade/tone.wav (in: 0s, out: 2s) -> tmp-fade/out.wav
```

## with invalid options

### should require a fade length

```execute
aux4 audio fade tmp-fade/tone.wav --output tmp-fade/none.wav
```

```error:partial
Error: provide --fadeIn and/or --fadeOut (in seconds)
```

### should reject a fade longer than the audio

```execute
aux4 audio fade tmp-fade/tone.wav --fadeIn 5 --output tmp-fade/long.wav
```

```error:partial
Error: fade is longer than the audio (3s)
```
