# audio compress

```beforeAll
mkdir -p tmp-compress
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=30 -ac 2 tmp-compress/speech.wav
```

```afterAll
rm -rf tmp-compress
```

## with default options

### should compress to a speech-friendly mp3

```execute
aux4 audio compress tmp-compress/speech.wav --output tmp-compress/speech.mp3
```

```expect:regex
^Compressed tmp-compress/speech\.wav -> tmp-compress/speech\.mp3 \(\d+ bytes, 32k\)$
```

### should produce mono 16 kHz audio at 32 kbps

```execute
aux4 audio info tmp-compress/speech.mp3
```

```expect:partial
  "codec": "mp3",
  "duration": 30,
  "sampleRate": 16000,
  "channels": 1,
  "channelLayout": "mono",
  "bitrate": 32000,
```

### should be much smaller than the input

```execute
test "$(wc -c < tmp-compress/speech.mp3)" -lt "$(wc -c < tmp-compress/speech.wav)" && echo smaller
```

```expect
smaller
```

### should name the output after the input when --output is omitted

```execute
aux4 audio compress tmp-compress/speech.wav
```

```expect:regex
^Compressed tmp-compress/speech\.wav -> tmp-compress/speech-compressed\.mp3 \(\d+ bytes, 32k\)$
```

## with other formats

### should compress to opus

```execute
aux4 audio compress tmp-compress/speech.wav --format opus --output tmp-compress/speech.opus && aux4 audio info tmp-compress/speech.opus
```

```expect:partial
  "codec": "opus",
**
  "channels": 1,
```

### should compress to m4a

```execute
aux4 audio compress tmp-compress/speech.wav --format m4a --bitrate 24k --output tmp-compress/speech.m4a && aux4 audio info tmp-compress/speech.m4a
```

```expect:partial
  "codec": "aac",
**
  "sampleRate": 16000,
  "channels": 1,
```

## with a maximum size

### should lower the bitrate so the file fits

```execute
aux4 audio compress tmp-compress/speech.wav --maxSize 50KB --output tmp-compress/fit.mp3
```

```expect:regex
^Compressed tmp-compress/speech\.wav -> tmp-compress/fit\.mp3 \(\d+ bytes, 8k\)$
```

### should stay under the maximum size

```execute
test "$(wc -c < tmp-compress/fit.mp3)" -le 50000 && echo fits
```

```expect
fits
```

### should fit opus under the maximum size too

```execute
aux4 audio compress tmp-compress/speech.wav --format opus --maxSize 50000 --output tmp-compress/fit.opus && test "$(wc -c < tmp-compress/fit.opus)" -le 50000 && echo fits
```

```expect:partial
fits
```

### should keep the requested bitrate when it already fits

```execute
aux4 audio compress tmp-compress/speech.wav --maxSize 24MB --output tmp-compress/roomy.mp3
```

```expect:regex
\(\d+ bytes, 32k\)$
```

### should fail and ask to split when even the minimum bitrate is too large

```execute
aux4 audio compress tmp-compress/speech.wav --maxSize 10KB --output tmp-compress/tiny.mp3
```

```error:partial
Error: tmp-compress/speech.wav (30s) cannot fit in 10KB even at the minimum bitrate of 8k: split it first (e.g. aux4 audio split --segment <seconds>)
```

## with invalid options

### should reject an unsupported format

```execute
aux4 audio compress tmp-compress/speech.wav --format wav --output tmp-compress/bad.wav
```

```error:partial
Error: unsupported format "wav": use one of mp3, opus, m4a
```

### should reject an output extension that does not match the format

```execute
aux4 audio compress tmp-compress/speech.wav --format opus --output tmp-compress/bad.mp3
```

```error:partial
Error: output file tmp-compress/bad.mp3 does not match --format opus: use a .opus extension
```

### should reject an invalid maximum size

```execute
aux4 audio compress tmp-compress/speech.wav --maxSize huge --output tmp-compress/bad.mp3
```

```error:partial
Error: invalid maxSize "huge": use bytes (e.g. 25000000) or a size like 500KB, 24MB
```

### should reject a bitrate below the minimum

```execute
aux4 audio compress tmp-compress/speech.wav --bitrate 4k --output tmp-compress/bad.mp3
```

```error:partial
Error: invalid bitrate "4k": must be at least 8k
```

### should refuse to overwrite an existing file

```execute
aux4 audio compress tmp-compress/speech.wav --output tmp-compress/speech.mp3
```

```error:partial
Error: output file already exists: tmp-compress/speech.mp3 (use --overwrite true to replace it)
```
