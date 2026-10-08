#### Description

The `split` command cuts an audio file into numbered parts. Choose one mode:

- **--segment <seconds>** — fixed-length parts; the last part holds the remainder
- **--silence true** — cut in the middle of every pause, which keeps words and sentences intact. A pause is audio quieter than `--threshold` dB for at least `--minSilence` seconds. Silence at the very beginning or end of the file does not create a cut.

Parts are named `<prefix>-001.<ext>`, `<prefix>-002.<ext>`, ... inside `--outputDir` (created if needed). The prefix defaults to the input file name and the format defaults to the input format. If any part already exists, nothing is written unless `--overwrite true` is given.

The command prints the number of parts followed by one path per line, in order. With `--json true` it prints only a JSON array instead, with the `path` of each part plus its `start` offset and `duration` in seconds within the original audio. Use the start offsets to shift timestamps computed on a part (for example a transcript) back to the original timeline.

#### Usage

```bash
aux4 audio split <input> (--segment <seconds> | --silence true) [--threshold <dB>] [--minSilence <seconds>] [--outputDir <dir>] [--prefix <name>] [--format <format>] [--overwrite <true|false>] [--json <true|false>]
```

<input>       Audio file to split
--segment     Length of each part in seconds
--silence     Split at pauses instead of fixed lengths (default: false)
--threshold   Level in dB below which audio counts as silence (default: -35); write negative values as --threshold=-40
--minSilence  Minimum pause length in seconds (default: 0.5)
--outputDir   Directory for the parts (default: .)
--prefix      File name prefix for the parts (default: input file name)
--format      Format of the parts: mp3, wav, flac, ogg, opus, m4a, aac, webm (default: input format)
--overwrite   Replace existing parts (default: false)
--json        Print a JSON array of {path, start, duration} instead of text (default: false)

#### Example

```bash
aux4 audio split lecture.mp3 --segment 600 --outputDir chunks
```

```text
Split lecture.mp3 into 3 segments:
chunks/lecture-001.mp3
chunks/lecture-002.mp3
chunks/lecture-003.mp3
```

```bash
aux4 audio split dictation.wav --silence true --minSilence 0.8 --prefix sentence
```

```text
Split dictation.wav into 4 segments:
sentence-001.wav
sentence-002.wav
sentence-003.wav
sentence-004.wav
```

```bash
aux4 audio split lecture.mp3 --segment 600 --outputDir chunks --json true
```

```json
[
  {
    "path": "chunks/lecture-001.mp3",
    "start": 0,
    "duration": 600
  },
  {
    "path": "chunks/lecture-002.mp3",
    "start": 600,
    "duration": 600
  },
  {
    "path": "chunks/lecture-003.mp3",
    "start": 1200,
    "duration": 287.412
  }
]
```
