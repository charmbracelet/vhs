# VHS

<p>
  <img src="https://user-images.githubusercontent.com/42545625/197049037-b38fea25-a885-4945-825e-d29842c5e44b.png#gh-dark-mode-only" width="500" />
  <img src="https://user-images.githubusercontent.com/42545625/197049039-83498ce6-d01d-4a08-8794-64770606ca8e.png#gh-light-mode-only" width="500" />
  <br>
  <a href="https://github.com/charmbracelet/vhs/releases"><img src="https://img.shields.io/github/release/charmbracelet/vhs.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/vhs?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/charmbracelet/vhs/actions"><img src="https://github.com/charmbracelet/vhs/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Write terminal GIFs as code for integration testing and demoing your CLI tools.

<img alt="Welcome to VHS" src="https://stuff.charm.sh/vhs/examples/neofetch_2.gif" width="600" />

The above example was generated with VHS ([view source](./examples/neofetch/neofetch.tape)).

## Tutorial

To get started, [install VHS](#installation) and create a new `.tape` file.

```sh
vhs new demo.tape
```

Open the `.tape` file with your favorite `$EDITOR`.

```sh
vim demo.tape
```

Tape files consist of a series of [commands](#commands). The commands are
instructions for VHS to perform on its virtual terminal.  For a list of all
possible commands see [the command reference](#vhs-command-reference).

```elixir
# Where should we write the GIF?
Output demo.gif

# Set up a 1200x600 terminal with 46px font.
Set FontSize 46
Set Width 1200
Set Height 600

# Type a command in the terminal.
Type "echo 'Welcome to VHS!'"

# Pause for dramatic effect...
Sleep 500ms

# Run the command by pressing enter.
Enter

# Admire the output for a bit.
Sleep 5s
```

Once you've finished, save the file and feed it into VHS.

```sh
vhs < demo.tape
```

All done! You should see a new file called `demo.gif` (or whatever you named
the `Output`) in the directory.

<img alt="A GIF produced by the VHS code above" src="https://stuff.charm.sh/vhs/examples/demo.gif" width="600" />

For more examples see the [`examples/`](https://github.com/charmbracelet/vhs/tree/main/examples) directory.

## Installation

> **Note**
> VHS requires [`ttyd`](https://github.com/tsl0922/ttyd) and [`ffmpeg`](https://ffmpeg.org) to be installed and available on your `PATH`.

Use a package manager:

```sh
# macOS or Linux
brew install charmbracelet/tap/vhs ffmpeg
brew install ttyd --HEAD

# Arch Linux (btw)
yay -S vhs ttyd ffmpeg

# Nix
nix-env -iA nixpkgs.vhs nixpkgs.ttyd nixpkgs.ffmpeg
```

Or, use Docker to run VHS directly, dependencies included:

```sh
docker run ghcr.io/charmbracelet/vhs <cassette>.tape
```

Or, download it:

* [Packages][releases] are available in Debian and RPM formats
* [Binaries][releases] are available for Linux, macOS, and Windows

Or, just install it with `go`:

```sh
go install github.com/charmbracelet/vhs@latest
```

[releases]: https://github.com/charmbracelet/vhs/releases

## The VHS Server

VHS has an SSH server built in! When you self host VHS you can access it as 
though it were installed locally. VHS will have access to commands and
applications on the host so you don't need to install them on your machine.

To start the server run:

```sh
vhs serve
```

Then, simply access VHS from any machine via `ssh`:

```sh
ssh vhs.example.com < demo.tape > demo.gif
```

## VHS Command Reference

> **Note**
> You can view all VHS documentation on the command line with `vhs manual`.

There are a few basic types of VHS commands:

* [`Output <path>`](#output): specify file output
* [`Set <Setting> Value`](#settings): set recording settings
* [`Type "<characters>"`](#type): emulate typing
* [`Backspace`](#backspace) [`Ctrl`](#ctrl) [`Enter`](#enter) [`Left`](#arrow-keys) [`Right`](#arrow-keys) [`Up`](#arrow-keys) [`Down`](#arrow-keys) [`Tab`](#tab) [`Space`](#tab): special keys
* [`Sleep <time>`](#sleep): wait for a certain amount of time
* [`Hide`](#hide): hide commands from output
* [`Show`](#show): stop hiding commands from output

### Output

The `Output` command allows you to specify the location and file format
of the render. You can specify more than one output in a tape file which
will render them to the respective locations.

```elixir
Output out.gif
Output out.mp4
Output out.webm
Output frames/ # a directory of frames as a PNG sequence
```

### Settings

The `Set` command allows you to change global aspects of the terminal, such as
the font settings, window dimensions, and GIF output location.

Setting must be administered at the top of the tape file. Any setting (except
`TypingSpeed`) applied after a non-setting or non-output command will be
ignored.

#### Set Font Size

Set the font size with the `Set FontSize <number>` command.

```elixir
Set FontSize 10
Set FontSize 20
Set FontSize 40
```

<img alt="Example of setting the font size to 10 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-10.gif" width="600" />

<img alt="Example of setting the font size to 20 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-20.gif" width="600" />

<img alt="Example of setting the font size to 40 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-40.gif" width="600" />

#### Set Font Family

Set the font family with the `Set FontFamily "<font>"` command

```elixir
Set FontFamily "Monoflow"
```

<img alt="Example of changing the font family to Monoflow" src="https://stuff.charm.sh/vhs/examples/font-family.gif" width="600" />

#### Set Height

Set the height of the terminal with the `Set Height` command.

```elixir
Set Height 600
Set Height 1000
```

#### Set Width

Set the width of the terminal with the `Set Width` command.

```elixir
Set Width 1200
Set Width 2000
```

#### Set Letter Spacing

Set the spacing between letters (tracking) with the `Set LetterSpacing`
Command.

```elixir
Set LetterSpacing 20
```

<img alt="Example of changing the letter spacing to 20 pixels between characters" src="https://stuff.charm.sh/vhs/examples/letter-spacing.gif" width="600" />

#### Set Line Height

Set the spacing between lines with the `Set LineHeight` Command.

```elixir
Set LineHeight 1.8
```

<img alt="Example of changing the line height to 1.8" src="https://stuff.charm.sh/vhs/examples/line-height.gif" width="600" />

#### Set Typing Speed

```elixir
Set TypingSpeed 500ms # 500ms
Set TypingSpeed 1s    # 1s
```

Set the typing speed of seconds per key press. For example, a typing speed of
`0.1` would result in a `0.1s` (`100ms`) delay between each character being typed.

This setting can also be overridden per command with the `@<time>` syntax.

```elixir
Set TypingSpeed 0.1
Type "100ms delay per character"
Type@500ms "500ms delay per character"
```

<img alt="Example of changing the typing speed to type different words" src="https://stuff.charm.sh/vhs/examples/typing-speed.gif" width="600" />

#### Set Theme

Set the theme of the terminal with the `Set Theme` command. The theme value
should be a JSON string with the base 16 colors and foreground + background.

```elixir
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "purple": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightPurple": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selectionBackground": "#3d3c58", "cursorColor": "#b3b0d6" }
```

<img alt="Example of changing the theme to Whimsy" src="https://stuff.charm.sh/vhs/examples/theme.gif" width="600" />

#### Set Padding

Set the padding (in pixels) of the terminal frame with the `Set Padding`
command.

```elixir
Set Padding 0
```

<img alt="Example of setting padding to 0" src="https://stuff.charm.sh/vhs/examples/padding.gif" width="600" />

#### Set Framerate

Set the rate at which VHS captures frames with the `Set Framerate` command.

```elixir
Set Framerate 60
```

#### Set Playback Speed

Set the playback speed of the final render.

```elixir
Set PlaybackSpeed 0.5 # Make output 2 times slower
Set PlaybackSpeed 1.0 # Keep output at normal speed (default)
Set PlaybackSpeed 2.0 # Make output 2 times faster
```

### Type

Use `Type` to emulate key presses. That is, you can use `Type` to script typing
in a terminal. Type is handy for both entering commands and interacting with
prompts and TUIs in the terminal. The command takes a string argument of the
characters to type.

You can set the standard typing speed with [`Set TypingSpeed`](#set-typing-speed)
and override it in places with a `@time` argument.

```elixir
# Type something
Type "Whatever you want"

# Type something really slowly!
Type@500ms "Slow down there, partner."
```

<img alt="Example of using the Type command in VHS" src="https://stuff.charm.sh/vhs/examples/type.gif" width="600" />

### Keys

Key commands take an optional `@time` and repeat `count`.
For example, the following presses the `Left` key 5 times with a 500
millisecond delay between each keystroke.

```elixir
Left@500ms 5
```

#### Backspace

Press the backspace key with the `Backspace` command.

```elixir
Backspace 18
```

<img alt="Example of pressing the Backspace key 18 times" src="https://stuff.charm.sh/vhs/examples/backspace.gif" width="600" />

#### Ctrl

You can access the control modifier and send control sequences with the `Ctrl`
command.

```elixir
Ctrl+R
```

<img alt="Example of pressing the Ctrl+R key to reverse search" src="https://stuff.charm.sh/vhs/examples/ctrl.gif" width="600" />

#### Enter

Press the enter key with the `Enter` command.

```elixir
Enter 2
```

<img alt="Example of pressing the Enter key twice" src="https://stuff.charm.sh/vhs/examples/enter.gif" width="600" />

#### Arrow Keys

Press any of the arrow keys with the `Up`, `Down`, `Left`, `Right` commands.

```elixir
Up 2
Down 2
Left
Right
Left
Right
Type "B"
Type "A"
```

<img alt="Example of pressing the arrow keys to navigate text" src="https://stuff.charm.sh/vhs/examples/arrow.gif" width="600" />

#### Tab

Enter a tab with the `Tab` command.

```elixir
Tab@500ms 2
```

<img alt="Example of pressing the tab key twice for autocomplete" src="https://stuff.charm.sh/vhs/examples/tab.gif" width="600" />

#### Space

Press the space bar with the `Space` command.

```elixir
Space 10
```

<img alt="Example of pressing the space key" src="https://stuff.charm.sh/vhs/examples/space.gif" width="600" />

### Sleep

The `Sleep` command allows you to continue capturing frames without interacting
with the terminal. This is useful when you need to wait on something to
complete while including it in the recording like a spinner or loading state.
The command takes a number argument in seconds.

```elixir
Sleep 0.5   # 500ms
Sleep 2     # 2s
Sleep 100ms # 100ms
Sleep 1s    # 1s
```

### Hide

The `Hide` command allows you to exclude commands from the output.

```elixir
Hide
```

This command is helpful for performing any setup and cleanup required to record
a GIF, such as building the latest version of a binary and removing the binary
once the demo is recorded.

```elixir
Output example.gif

# Setup
Hide
Type "go build -o example . && clear"
Enter
Show

# Recording...
Type 'Running ./example'
...
Enter

# Cleanup
Hide
Type 'rm example'
```

### Show

The `Show` command allows you to specify that the following commands should be
un-hidden in the output. Since commands are shown by default, this command is
really only useful after using the `Hide` command.

```elixir
Hide
Type "You won't see this being typed."
Show
Type "You will see this being typed."
```

<img alt="Example of typing something while hidden" src="https://stuff.charm.sh/vhs/examples/hide.gif" width="600" />

***

## Testing

VHS can be used for integration testing by outputting an `.txt` or `.ascii`
file. This will result in a golden file that you can commit to your git repo
and view diffs between previous and future runs.

```elixir
Output golden.ascii
```

You can also integrate VHS into your CI pipeline with the
[VHS GitHub Action](https://github.com/charmbracelet/vhs-action)

## Syntax Highlighting

If your editor supports syntax highlighting from tree-sitter,
you can syntax highlight `.tape` files with [Tree-sitter VHS Grammar](https://github.com/charmbracelet/tree-sitter-vhs).

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.social/@charmcli)
* [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/vhs/raw/main/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/">
  <img
    alt="The Charm logo"
    width="400"
    src="https://stuff.charm.sh/charm-badge.jpg"
  />
</a>

Charm热爱开源 • Charm loves open source
