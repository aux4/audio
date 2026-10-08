#### Description

The `waveform` command renders a picture of the audio's waveform, useful for previews, thumbnails and visually spotting silence or clipping. The image size is given as `WIDTHxHEIGHT` and the color as a name (`blue`, `white`) or a hex value (`#3b82f6`). The background is transparent in PNG output.

With `--output`, the image is saved to that file, which must end in `.png` or `.jpg`. Without it, a PNG is streamed to stdout. When no input file is given, the audio is read from stdin.

#### Usage

```bash
aux4 audio waveform [<input>] [--size <WxH>] [--color <color>] [--output <file>] [--overwrite <true|false>]
```

<input>      Audio file to render (default: read from stdin)
--size       Image size as WIDTHxHEIGHT (default: 1200x200)
--color      Waveform color name or hex value (default: #3b82f6)
--output     Output image path, .png or .jpg (default: PNG to stdout)
--overwrite  Replace the output file if it already exists (default: false)

#### Example

```bash
aux4 audio waveform episode.mp3 --size 800x120 --color "#111827" --output cover-wave.png
```

```text
Rendered waveform of episode.mp3 (800x120) -> cover-wave.png
```
