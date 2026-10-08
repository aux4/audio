#### Description

The `trim-silence` command removes the silence at the beginning and end of a recording, such as the dead air before someone starts speaking and after they stop. Audio quieter than `--threshold` dB counts as silence.

With `--all true`, pauses inside the recording are shortened as well: any pause longer than `--maxPause` seconds is cut down to that length, which tightens up long hesitations without running words together.

Use a lower threshold (e.g. `--threshold=-60`) for quiet recordings and a higher one (e.g. `--threshold=-40`) for noisy ones. Negative values must be written with an equals sign.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). When no input file is given, the audio is read from stdin.

#### Usage

```bash
aux4 audio trim-silence [<input>] [--threshold <dB>] [--all <true|false>] [--maxPause <seconds>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>      Audio file to clean up (default: read from stdin)
--threshold  Level in dB below which audio counts as silence (default: -50)
--all        Also shorten pauses inside the audio (default: false)
--maxPause   Longest pause in seconds kept when --all is true (default: 0.5)
--output     Output file path (default: write to stdout)
--format     Output format when writing to stdout (default: input format, else wav)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio trim-silence voice-note.m4a --output voice-note-clean.m4a
```

```text
Removed silence from voice-note.m4a -> voice-note-clean.m4a
```

```bash
aux4 audio trim-silence interview.wav --all true --maxPause 0.3 --output interview-tight.wav
```

```text
Removed silence from interview.wav -> interview-tight.wav
```
