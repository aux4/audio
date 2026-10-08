# audio split

```beforeAll
mkdir -p tmp-split
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2.5 tmp-split/tone.wav
ffmpeg -hide_banner -loglevel error -y -f lavfi -i "aevalsrc=if(between(t\,1\,2)\,0\,0.5*sin(2*PI*440*t)):s=44100:d=3" tmp-split/speech.wav
```

```afterAll
rm -rf tmp-split
```

## by segment length

### should split into fixed-length parts

```execute
aux4 audio split tmp-split/tone.wav --segment 1 --outputDir tmp-split/parts
```

```expect
Split tmp-split/tone.wav into 3 segments:
tmp-split/parts/tone-001.wav
tmp-split/parts/tone-002.wav
tmp-split/parts/tone-003.wav
```

### should keep the remainder in the last part

```execute
aux4 audio info tmp-split/parts/tone-003.wav
```

```expect:partial
  "duration": 0.5,
```

## at silences

### should cut in the middle of each pause

```execute
aux4 audio split tmp-split/speech.wav --silence true --outputDir tmp-split/phrases --prefix phrase
```

```expect
Split tmp-split/speech.wav into 2 segments:
tmp-split/phrases/phrase-001.wav
tmp-split/phrases/phrase-002.wav
```

### should produce parts that end inside the pause

```execute
aux4 audio info tmp-split/phrases/phrase-001.wav
```

```expect:partial
  "duration": 1.5,
```

## as JSON

### should list each part with its start offset and duration

```execute
aux4 audio split tmp-split/tone.wav --segment 1 --outputDir tmp-split/json --json true
```

```expect:json
[
  {
    "path": "tmp-split/json/tone-001.wav",
    "start": 0,
    "duration": 1
  },
  {
    "path": "tmp-split/json/tone-002.wav",
    "start": 1,
    "duration": 1
  },
  {
    "path": "tmp-split/json/tone-003.wav",
    "start": 2,
    "duration": 0.5
  }
]
```

### should report the cut points found at silences

```execute
aux4 audio split tmp-split/speech.wav --silence true --outputDir tmp-split/json-silence --json true
```

```expect:json
[
  {
    "path": "tmp-split/json-silence/speech-001.wav",
    "start": 0,
    "duration": 1.5
  },
  {
    "path": "tmp-split/json-silence/speech-002.wav",
    "start": 1.5,
    "duration": 1.5
  }
]
```

## with a format

### should encode the parts in that format

```execute
aux4 audio split tmp-split/tone.wav --segment 2 --format mp3 --outputDir tmp-split/mp3
```

```expect
Split tmp-split/tone.wav into 2 segments:
tmp-split/mp3/tone-001.mp3
tmp-split/mp3/tone-002.mp3
```

## with existing parts

### should refuse to overwrite them

```execute
aux4 audio split tmp-split/tone.wav --segment 2 --format mp3 --outputDir tmp-split/mp3
```

```error:partial
Error: output file already exists: tmp-split/mp3/tone-001.mp3 (use --overwrite true to replace it)
```

## with invalid options

### should require a segment length or silence mode

```execute
aux4 audio split tmp-split/tone.wav
```

```error:partial
Error: use either --segment <seconds> or --silence true
```

### should reject a non-numeric segment length

```execute
aux4 audio split tmp-split/tone.wav --segment abc
```

```error:partial
Error: invalid segment "abc": must be a positive number
```
