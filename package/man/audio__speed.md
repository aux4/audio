#### Description

The `speed` command plays audio faster or slower while keeping the original pitch, so voices do not sound higher or lower. A `--factor` of `2` makes the audio twice as fast (half the duration); `0.5` makes it half as fast (twice the duration). Factors from `0.1` to `10` are supported.

#### Usage

```bash
aux4 audio speed <input> --factor <factor> [--output <file>] [--overwrite <true|false>]
```

<input>      Audio file to change
--factor     Speed factor from 0.1 to 10 (e.g. 1.25 = 25% faster)
--output     Output file path (default: <input>-speed.<ext>)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio speed lecture.mp3 --factor 1.5
```

```text
Changed speed of lecture.mp3 (1.5x) -> lecture-speed.mp3
```
