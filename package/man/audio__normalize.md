#### Description

The `normalize` command adjusts the perceived loudness of a recording to a target level using the EBU R128 standard. It runs in two passes: the first measures the loudness of the whole file, the second applies a single linear gain so the result hits the target without changing the dynamics. The original sample rate is preserved.

Common targets:

- **-16 LUFS** — podcasts and spoken word (default)
- **-14 LUFS** — music streaming services
- **-23 LUFS** — broadcast (EBU R128)

`--truePeak` limits the maximum peak level and `--loudnessRange` sets the allowed loudness range. Negative values must be written with an equals sign (`--target=-14`). Files that are completely silent cannot be normalized and produce an error.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). Piped input is read in full before the first pass.

#### Usage

```bash
aux4 audio normalize [<input>] [--target <LUFS>] [--truePeak <dBTP>] [--loudnessRange <LU>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>          Audio file to normalize (default: read from stdin)
--target         Integrated loudness target in LUFS, -70 to -5 (default: -16)
--truePeak       Maximum true peak in dBTP, -9 to 0 (default: -1.5)
--loudnessRange  Loudness range target in LU, 1 to 50 (default: 11)
--output         Output file path (default: write to stdout)
--format         Output format when writing to stdout (default: input format, else wav)
--overwrite      Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio normalize episode.wav --output episode-normalized.wav
```

```text
Normalized episode.wav (-16 LUFS) -> episode-normalized.wav
```

```bash
aux4 audio normalize song.flac --target=-14 --truePeak=-1 --output song-master.flac
```

```text
Normalized song.flac (-14 LUFS) -> song-master.flac
```
