#### Description

The `trim` command keeps a section of an audio file and discards the rest. The section is defined by:

- **--start** — where the section begins (default: the beginning of the file)
- **--end** — where the section ends (default: the end of the file)
- **--duration** — how long the section is, as an alternative to `--end`

At least one of the three must be given, and `--end` and `--duration` cannot be combined. Times are seconds (`90`, `1.5`) or timecodes (`01:30`, `00:01:30.250`).

The command fails if `--end` is not after `--start`, or if `--start` is beyond the end of the audio. The output is re-encoded, so cuts are sample-accurate.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). When no input file is given, the audio is read from stdin.

#### Usage

```bash
aux4 audio trim [<input>] [--start <time>] [--end <time>] [--duration <time>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>      Audio file to trim (default: read from stdin)
--start      Start time in seconds or [HH:]MM:SS[.ms]
--end        End time in seconds or [HH:]MM:SS[.ms]
--duration   Length to keep in seconds or [HH:]MM:SS[.ms]
--output     Output file path (default: write to stdout)
--format     Output format when writing to stdout (default: input format, else wav)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio trim episode.mp3 --start 00:01:30 --end 00:02:45 --output clip.mp3
```

```text
Trimmed episode.mp3 (90s - 165s) -> clip.mp3
```

```bash
aux4 audio trim voicemail.wav --duration 10 > first-ten-seconds.wav
```
