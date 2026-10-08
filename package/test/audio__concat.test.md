# audio concat

```beforeAll
mkdir -p tmp-concat
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=1 -ac 2 tmp-concat/first.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=660:sample_rate=44100:duration=2 -ac 2 tmp-concat/second.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=880:sample_rate=16000:duration=1 -c:a libmp3lame tmp-concat/third.mp3
```

```afterAll
rm -rf tmp-concat
```

## with two files

### should join them in order

```execute
aux4 audio concat tmp-concat/first.wav tmp-concat/second.wav --output tmp-concat/joined.wav
```

```expect
Concatenated 2 files -> tmp-concat/joined.wav
```

### should add up the durations

```execute
aux4 audio info tmp-concat/joined.wav
```

```expect:partial
  "duration": 3,
  "sampleRate": 44100,
  "channels": 2,
```

## with files in different formats

### should convert every file to the format of the first one

```execute
aux4 audio concat tmp-concat/first.wav tmp-concat/third.mp3 tmp-concat/second.wav --output tmp-concat/mixed.mp3 && aux4 audio info tmp-concat/mixed.mp3
```

```expect:partial
Concatenated 3 files -> tmp-concat/mixed.mp3
**
  "codec": "mp3",
**
  "sampleRate": 44100,
  "channels": 2,
```

## with invalid input

### should require at least two files

```execute
aux4 audio concat tmp-concat/first.wav --output tmp-concat/single.wav
```

```error:partial
Error: concat needs at least two input files
```

### should fail when one of the files is missing

```execute
aux4 audio concat tmp-concat/first.wav tmp-concat/missing.wav --output tmp-concat/broken.wav
```

```error:partial
Error: input file not found: tmp-concat/missing.wav
```

### should refuse to write over one of the inputs

```execute
aux4 audio concat tmp-concat/first.wav tmp-concat/second.wav --output tmp-concat/second.wav
```

```error:partial
Error: output file is the same as the input: tmp-concat/second.wav (use --output to choose another file)
```
