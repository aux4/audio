# audio speech-prep

```beforeAll
mkdir -p tmp-speech-prep
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:duration=2 -ac 2 -c:a libopus -f webm pipe:1 > tmp-speech-prep/chrome.webm
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 -ac 2 -c:a aac tmp-speech-prep/safari.mp4
```

```afterAll
rm -rf tmp-speech-prep
```

## with a browser webm/opus recording

### should write a 16 kHz mono WAV

```execute
aux4 audio speech-prep tmp-speech-prep/chrome.webm
```

```expect
Prepared tmp-speech-prep/chrome.webm for speech recognition (16000 Hz mono WAV) -> tmp-speech-prep/chrome-16k.wav
```

### should have the format speech recognizers expect

```execute
aux4 audio info tmp-speech-prep/chrome-16k.wav
```

```expect:partial
  "format": "wav",
  "codec": "pcm_s16le",
  "duration": 2,
  "sampleRate": 16000,
  "channels": 1,
```

## with a Safari mp4/aac recording

### should convert it too

```execute
aux4 audio speech-prep tmp-speech-prep/safari.mp4 --output tmp-speech-prep/safari.wav && aux4 audio info tmp-speech-prep/safari.wav
```

```expect:partial
Prepared tmp-speech-prep/safari.mp4 for speech recognition (16000 Hz mono WAV) -> tmp-speech-prep/safari.wav
**
  "codec": "pcm_s16le",
**
  "sampleRate": 16000,
  "channels": 1,
```

## with a custom sample rate

### should use it

```execute
aux4 audio speech-prep tmp-speech-prep/chrome.webm --sampleRate 8000 --output tmp-speech-prep/phone.wav
```

```expect
Prepared tmp-speech-prep/chrome.webm for speech recognition (8000 Hz mono WAV) -> tmp-speech-prep/phone.wav
```

## with a non-wav output

### should reject it

```execute
aux4 audio speech-prep tmp-speech-prep/chrome.webm --output tmp-speech-prep/out.mp3
```

```error:partial
Error: speech-prep writes WAV audio: the output file must end in .wav (got tmp-speech-prep/out.mp3)
```
