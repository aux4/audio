# audio trim

```beforeAll
mkdir -p tmp-trim
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=3 tmp-trim/tone.wav
```

```afterAll
rm -rf tmp-trim
```

## with start and end

### should keep the section between start and end

```execute
aux4 audio trim tmp-trim/tone.wav --start 1 --end 2.5 --output tmp-trim/tone-trimmed.wav
```

```expect
Trimmed tmp-trim/tone.wav (1s - 2.5s) -> tmp-trim/tone-trimmed.wav
```

### should produce a 1.5 second file

```execute
aux4 audio info tmp-trim/tone-trimmed.wav
```

```expect:partial
  "duration": 1.5,
```

## with duration

### should keep the given length from the beginning

```execute
aux4 audio trim tmp-trim/tone.wav --duration 0.75 --output tmp-trim/first.wav && aux4 audio info tmp-trim/first.wav
```

```expect:partial
Trimmed tmp-trim/tone.wav (0s - 0.75s) -> tmp-trim/first.wav
**
  "duration": 0.75,
```

## with a timecode

### should accept MM:SS.ms values

```execute
aux4 audio trim tmp-trim/tone.wav --start 00:02 --output tmp-trim/tail.wav && aux4 audio info tmp-trim/tail.wav
```

```expect:partial
Trimmed tmp-trim/tone.wav (2s - 3s) -> tmp-trim/tail.wav
**
  "duration": 1,
```

## with invalid options

### should require at least one boundary

```execute
aux4 audio trim tmp-trim/tone.wav --output tmp-trim/none.wav
```

```error:partial
Error: provide at least one of --start, --end or --duration
```

### should reject an end before the start

```execute
aux4 audio trim tmp-trim/tone.wav --start 2 --end 1 --output tmp-trim/bad.wav
```

```error:partial
Error: end (1s) must be after start (2s)
```

### should reject a start beyond the end of the audio

```execute
aux4 audio trim tmp-trim/tone.wav --start 10 --output tmp-trim/bad.wav
```

```error:partial
Error: start (10s) is beyond the end of the audio (3s)
```

### should reject end and duration together

```execute
aux4 audio trim tmp-trim/tone.wav --end 2 --duration 1 --output tmp-trim/bad.wav
```

```error:partial
Error: use either --end or --duration, not both
```

### should reject an invalid time

```execute
aux4 audio trim tmp-trim/tone.wav --start 1s --output tmp-trim/bad.wav
```

```error:partial
Error: invalid start "1s": use seconds (e.g. 1.5) or [HH:]MM:SS[.ms]
```
