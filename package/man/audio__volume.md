#### Description

The `volume` command makes audio louder or quieter by a fixed amount. The `--level` can be:

- **a multiplier** — `0.5` halves the amplitude, `2` doubles it
- **a gain in decibels** — `6dB` is louder, `-3dB` is quieter

A negative gain must be written with an equals sign (`--level=-3dB`) so it is not read as another flag. Raising the volume can clip loud passages; use `normalize` to reach a target loudness safely.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). When no input file is given, the audio is read from stdin.

#### Usage

```bash
aux4 audio volume [<input>] --level <level> [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>      Audio file to adjust (default: read from stdin)
--level      Multiplier (e.g. 0.5, 1.5) or gain in decibels (e.g. 6dB, -3dB)
--output     Output file path (default: write to stdout)
--format     Output format when writing to stdout (default: input format, else wav)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio volume background.mp3 --level 0.3 --output background-quiet.mp3
```

```text
Adjusted volume of background.mp3 (0.3) -> background-quiet.mp3
```

```bash
aux4 audio volume voice.wav --level=-6dB --output voice-quiet.wav
```

```text
Adjusted volume of voice.wav (-6dB) -> voice-quiet.wav
```
