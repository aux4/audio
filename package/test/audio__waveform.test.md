# audio waveform

```beforeAll
mkdir -p tmp-waveform
ffmpeg -hide_banner -loglevel error -y -f lavfi -i sine=frequency=440:sample_rate=44100:duration=2 tmp-waveform/tone.wav
```

```afterAll
rm -rf tmp-waveform
```

## with an output file

### should render a PNG file

```execute
aux4 audio waveform tmp-waveform/tone.wav --output tmp-waveform/tone-waveform.png && test -s tmp-waveform/tone-waveform.png && echo "image written"
```

```expect
Rendered waveform of tmp-waveform/tone.wav (1200x200) -> tmp-waveform/tone-waveform.png
image written
```

## with size and color

### should render with the given options

```execute
aux4 audio waveform tmp-waveform/tone.wav --size 600x100 --color red --output tmp-waveform/small.png
```

```expect
Rendered waveform of tmp-waveform/tone.wav (600x100) -> tmp-waveform/small.png
```

## with invalid options

### should reject an invalid size

```execute
aux4 audio waveform tmp-waveform/tone.wav --size big --output tmp-waveform/bad.png
```

```error:partial
Error: invalid size "big": use WIDTHxHEIGHT (e.g. 1200x200)
```

### should reject an invalid color

```execute
aux4 audio waveform tmp-waveform/tone.wav --color "red:split_channels=1" --output tmp-waveform/bad.png
```

```error:partial
Error: invalid color "red:split_channels=1": use a color name (e.g. blue) or hex value (e.g. #3b82f6)
```

### should reject a non-image output

```execute
aux4 audio waveform tmp-waveform/tone.wav --output tmp-waveform/wave.wav
```

```error:partial
Error: waveform output must be a .png or .jpg file (got tmp-waveform/wave.wav)
```
