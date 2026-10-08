#### Description

The `fade` command gradually raises the volume at the start of the audio (fade-in), lowers it at the end (fade-out), or both. Lengths are in seconds; at least one of `--fadeIn` and `--fadeOut` must be greater than zero, and neither may be longer than the audio. The duration of the audio does not change.

#### Usage

```bash
aux4 audio fade <input> [--fadeIn <seconds>] [--fadeOut <seconds>] [--output <file>] [--overwrite <true|false>]
```

<input>      Audio file to fade
--fadeIn     Fade-in length in seconds (default: 0)
--fadeOut    Fade-out length in seconds (default: 0)
--output     Output file path (default: <input>-fade.<ext>)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio fade jingle.wav --fadeIn 0.5 --fadeOut 2
```

```text
Faded jingle.wav (in: 0.5s, out: 2s) -> jingle-fade.wav
```
