#### Description

The `speed` command plays audio faster or slower while keeping the original pitch, so voices do not sound higher or lower. A `--factor` of `2` makes the audio twice as fast (half the duration); `0.5` makes it half as fast (twice the duration). Factors from `0.1` to `10` are supported.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). When no input file is given, the audio is read from stdin.

#### Usage

```bash
aux4 audio speed [<input>] --factor <factor> [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>      Audio file to change (default: read from stdin)
--factor     Speed factor from 0.1 to 10 (e.g. 1.25 = 25% faster)
--output     Output file path (default: write to stdout)
--format     Output format when writing to stdout (default: input format, else wav)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio speed lecture.mp3 --factor 1.5 --output lecture-fast.mp3
```

```text
Changed speed of lecture.mp3 (1.5x) -> lecture-fast.mp3
```
