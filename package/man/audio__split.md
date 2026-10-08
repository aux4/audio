#### Description

The `split` command cuts an audio file into numbered parts. Choose one mode:

- **--segment <seconds>** — fixed-length parts; the last part holds the remainder
- **--silence true** — cut in the middle of every pause, which keeps words and sentences intact. A pause is audio quieter than `--threshold` dB for at least `--minSilence` seconds. Silence at the very beginning or end of the file does not create a cut.

Parts are named `<prefix>-001.<ext>`, `<prefix>-002.<ext>`, ... inside `--outputDir` (created if needed). The prefix defaults to the input file name and the format defaults to the input format. If any part already exists, nothing is written unless `--overwrite true` is given.

The command prints the number of parts followed by one path per line, which makes the output easy to feed into other commands.

#### Usage

```bash
aux4 audio split <input> (--segment <seconds> | --silence true) [--threshold <dB>] [--minSilence <seconds>] [--outputDir <dir>] [--prefix <name>] [--format <format>] [--overwrite <true|false>]
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
