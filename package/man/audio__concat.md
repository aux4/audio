#### Description

The `concat` command joins two or more audio files, in the order given, into a single file. The inputs may have different formats, sample rates and channel counts: every input is converted to the sample rate and channel count of the first file before joining. The output is encoded in the format of its extension.

`--output` is required. The output may not be one of the inputs, and an existing file is only replaced with `--overwrite true`.

#### Usage

```bash
aux4 audio concat <input> <input> [<input> ...] --output <file> [--overwrite <true|false>]
```

<input>      Audio files to join, in order (at least two)
--output     Output file path
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio concat intro.wav interview.mp3 outro.wav --output episode.mp3
```

```text
Concatenated 3 files -> episode.mp3
```
