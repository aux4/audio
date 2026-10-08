#### Description

The `fade` command gradually raises the volume at the start of the audio (fade-in), lowers it at the end (fade-out), or both. Lengths are in seconds; at least one of `--fadeIn` and `--fadeOut` must be greater than zero, and neither may be longer than the audio. The duration of the audio does not change.

With `--output`, the result is saved in the format of the file extension. Without it, the result is streamed to stdout in `--format`, or in the input format (WAV when it cannot be recognized). Piped input is read in full first, because the fade-out position depends on the total duration.

#### Usage

```bash
aux4 audio fade [<input>] [--fadeIn <seconds>] [--fadeOut <seconds>] [--output <file>] [--format <format>] [--overwrite <true|false>]
```

<input>      Audio file to fade (default: read from stdin)
--fadeIn     Fade-in length in seconds (default: 0)
--fadeOut    Fade-out length in seconds (default: 0)
--output     Output file path (default: write to stdout)
--format     Output format when writing to stdout (default: input format, else wav)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio fade jingle.wav --fadeIn 0.5 --fadeOut 2 --output jingle-fade.wav
```

```text
Faded jingle.wav (in: 0.5s, out: 2s) -> jingle-fade.wav
```
