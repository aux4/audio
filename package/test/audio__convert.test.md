# audio convert

```beforeAll
mkdir -p tmp-convert
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 -ac 2 tmp-convert/tone.wav
```

```afterAll
rm -rf tmp-convert
```

## with an output file

### should write the converted file

```execute
aux4 audio convert tmp-convert/tone.wav --format mp3 --output tmp-convert/tone.mp3
```

```expect
Converted tmp-convert/tone.wav -> tmp-convert/tone.mp3
```

### should produce an mp3 file

```execute
aux4 audio info tmp-convert/tone.mp3
```

```expect:partial
  "format": "mp3",
  "codec": "mp3",
```

## with bitrate

### should encode at the requested bitrate

```execute
aux4 audio convert tmp-convert/tone.wav --format mp3 --bitrate 64k --output tmp-convert/tone-64k.mp3 && aux4 audio info tmp-convert/tone-64k.mp3
```

```expect:partial
Converted tmp-convert/tone.wav -> tmp-convert/tone-64k.mp3
**
  "bitrate": 64000,
```

## with sample rate and channels

### should resample to mono flac

```execute
aux4 audio convert tmp-convert/tone.wav --format flac --sampleRate 16000 --channels 1 --output tmp-convert/tone.flac && aux4 audio info tmp-convert/tone.flac
```

```expect:partial
Converted tmp-convert/tone.wav -> tmp-convert/tone.flac
**
  "format": "flac",
  "codec": "flac",
  "duration": 2,
  "sampleRate": 16000,
  "channels": 1,
```

## with other formats

### should convert to opus

```execute
aux4 audio convert tmp-convert/tone.wav --format opus --output tmp-convert/tone.opus && aux4 audio info tmp-convert/tone.opus
```

```expect:partial
  "codec": "opus",
```

### should convert to m4a

```execute
aux4 audio convert tmp-convert/tone.wav --format m4a --output tmp-convert/tone.m4a && aux4 audio info tmp-convert/tone.m4a
```

```expect:partial
  "codec": "aac",
```

### should convert to webm

```execute
aux4 audio convert tmp-convert/tone.wav --format webm --output tmp-convert/tone.webm && aux4 audio info tmp-convert/tone.webm
```

```expect:partial
  "format": "matroska,webm",
  "codec": "opus",
```

### should infer the format from the output extension

```execute
aux4 audio convert tmp-convert/tone.wav --output tmp-convert/from-ext.ogg && aux4 audio info tmp-convert/from-ext.ogg
```

```expect:partial
Converted tmp-convert/tone.wav -> tmp-convert/from-ext.ogg
**
  "format": "ogg",
```

## with an existing output file

### should refuse to overwrite it

```execute
aux4 audio convert tmp-convert/tone.wav --output tmp-convert/tone.mp3
```

```error:partial
Error: output file already exists: tmp-convert/tone.mp3 (use --overwrite true to replace it)
```

### should replace it with --overwrite true

```execute
aux4 audio convert tmp-convert/tone.wav --output tmp-convert/tone.mp3 --overwrite true
```

```expect
Converted tmp-convert/tone.wav -> tmp-convert/tone.mp3
```

## with invalid options

### should reject an unsupported format

```execute
aux4 audio convert tmp-convert/tone.wav --format xyz
```

```error:partial
Error: unsupported format "xyz": use one of mp3, wav, flac, ogg, opus, m4a, aac, webm
```

### should require a format when writing to stdout

```execute
aux4 audio convert tmp-convert/tone.wav
```

```error:partial
Error: --format is required when writing to stdout (no --output given): use one of mp3, wav, flac, ogg, opus, m4a, aac, webm
```

### should reject an output extension that does not match the format

```execute
aux4 audio convert tmp-convert/tone.wav --format mp3 --output tmp-convert/mismatch.ogg
```

```error:partial
Error: output file tmp-convert/mismatch.ogg does not match --format mp3: use a .mp3 extension
```

### should reject an invalid bitrate

```execute
aux4 audio convert tmp-convert/tone.wav --format mp3 --bitrate fast --output tmp-convert/bad.mp3
```

```error:partial
Error: invalid bitrate "fast": use a value like 128k or 192000
```

### should refuse to write over the input file

```execute
aux4 audio convert tmp-convert/tone.wav --output tmp-convert/tone.wav
```

```error:partial
Error: output file is the same as the input: tmp-convert/tone.wav (use --output to choose another file)
```

### should fail when the input does not exist

```execute
aux4 audio convert tmp-convert/missing.wav --format mp3
```

```error:partial
Error: input file not found: tmp-convert/missing.wav
```

## when ffmpeg cannot be found

### should explain how to install it

```execute
FFMPEG_PATH=/nonexistent/ffmpeg aux4 audio convert tmp-convert/tone.wav --format flac --output tmp-convert/no-ffmpeg.flac
```

```error:partial
Error: ffmpeg not found at /nonexistent/ffmpeg (set by FFMPEG_PATH); install it with 'brew install ffmpeg' (macOS)*
```
