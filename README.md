# VHS

<p>
  <img src="https://user-images.githubusercontent.com/42545625/198402537-12ca2f6c-0779-4eb8-a67c-8db9cb3df13c.png#gh-dark-mode-only" width="500" />
  <img src="https://user-images.githubusercontent.com/42545625/198402542-a305f669-a05a-4d91-b18b-ca76e72b655a.png#gh-light-mode-only" width="500" />
  <br>
  <a href="https://github.com/charmbracelet/vhs/releases"><img src="https://img.shields.io/github/release/charmbracelet/vhs.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/charmbracelet/vhs?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/charmbracelet/vhs/actions"><img src="https://github.com/charmbracelet/vhs/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Write terminal GIFs as code for integration testing and demoing your CLI tools.

<img alt="Welcome to VHS" src="https://stuff.charm.sh/vhs/examples/neofetch_3.gif" width="600" />

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

Tape files consist of a series of [commands](#vhs-command-reference). The commands are
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
vhs demo.tape
```

All done! You should see a new file called `demo.gif` (or whatever you named
the `Output`) in the directory.

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/demo.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/demo.gif">
  <img width="600" alt="A GIF produced by the VHS code above" src="https://stuff.charm.sh/vhs/examples/demo.gif">
</picture>

For more examples see the [`examples/`](https://github.com/charmbracelet/vhs/tree/main/examples) directory.

## Installation

> [!NOTE]
> VHS requires [`ttyd`](https://github.com/tsl0922/ttyd) and [`ffmpeg`](https://ffmpeg.org) to be installed and available on your `PATH`.

Use a package manager:

```sh
# macOS or Linux
brew install vhs

# macOS (via MacPorts)
sudo port install vhs

# Arch Linux (btw)
pacman -S vhs

# Nix
nix-env -iA nixpkgs.vhs

# Debian/Ubuntu
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
# Install ttyd from https://github.com/tsl0922/ttyd/releases
sudo apt update && sudo apt install vhs ffmpeg

# Fedora/RHEL
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
# Install ttyd from https://github.com/tsl0922/ttyd/releases
sudo yum install vhs ffmpeg

# Void Linux
sudo xbps-install vhs

# Windows using winget
# Note: Winget will install all the needed dependencies, which're FFmpeg and ttyd.
#       No need to do any prerequisites to install vhs using this method.
winget install charmbracelet.vhs

# Windows using scoop
scoop install vhs

```

Or, use Docker to run VHS directly, dependencies included:

```sh
docker run --rm -v $PWD:/vhs ghcr.io/charmbracelet/vhs <cassette>.tape
```

Or, download it:

* [Packages][releases] are available in Debian and RPM formats
* [Binaries][releases] are available for Linux, macOS, and Windows

Or, just install it with `go`:

```sh
go install github.com/charmbracelet/vhs@latest
```

[releases]: https://github.com/charmbracelet/vhs/releases

## Record Tapes

VHS has the ability to generate tape files from your terminal actions!

To record to a tape file, run:

```bash
vhs record > cassette.tape
```

Perform any actions you want and then `exit` the terminal session to stop
recording. You may want to manually edit the generated `.tape` file to add
settings or modify actions. Then, you can generate the GIF:

```bash
vhs cassette.tape
```

## Publish Tapes

VHS allows you to publish your GIFs to our servers for easy sharing with your
friends and colleagues. Specify which file you want to share, then use the
`publish` sub-command to host it on `vhs.charm.sh`. The output will provide you
with links to share your GIF via browser, HTML, and Markdown. 

```bash
vhs publish demo.gif
```

## The VHS Server

VHS has an SSH server built in! When you self-host VHS you can access it as
though it were installed locally. VHS will have access to commands and
applications on the host, so you don't need to install them on your machine.

To start the server run:

```sh
vhs serve
```

<details>
<summary>Configuration Options</summary>

* `VHS_PORT`: The port to listen on (`1976`)
* `VHS_HOST`: The host to listen on (`localhost`)
* `VHS_GID`: The Group ID to run the server as (current user's GID)
* `VHS_UID`: The User ID to run the server as (current user's UID)
* `VHS_KEY_PATH`: The path to the SSH key to use (`.ssh/vhs_ed25519`)
* `VHS_AUTHORIZED_KEYS_PATH`: The path to the authorized keys file (empty, publicly accessible)

</details>


Then, simply access VHS from a different machine via `ssh`:

```sh
ssh vhs.example.com < demo.tape > demo.gif
```

## VHS Command Reference

> [!NOTE]
> You can view all VHS documentation on the command line with `vhs manual`.

There are a few basic types of VHS commands:

* [`Output <path>`](#output): specify file output
* [`Require <program>`](#require): specify required programs for tape file
* [`Set <Setting> Value`](#settings): set recording settings
* [`Type "<characters>"`](#type): emulate typing
* [`Left`](#arrow-keys) [`Right`](#arrow-keys) [`Up`](#arrow-keys) [`Down`](#arrow-keys): arrow keys
* [`Backspace`](#backspace) [`Enter`](#enter) [`Tab`](#tab) [`Space`](#space): special keys
* [`Ctrl[+Alt][+Shift]+<char>`](#ctrl): press control + key and/or modifier
* [`Sleep <time>`](#sleep): wait for a certain amount of time
* [`Hide`](#hide): hide commands from output
* [`Show`](#show): stop hiding commands from output
* [`Screenshot`](#screenshot): screenshot the current frame
* [`Copy/Paste`](#copy--paste): copy and paste text from clipboard.
* [`Source`](#source): source commands from another tape

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

### Require

The `Require` command allows you to specify dependencies for your tape file.
These are useful to fail early if a required program is missing from the
`$PATH`, and it is certain that the VHS execution will not work as expected.

Require commands must be defined at the top of a tape file, before any non-
setting or non-output command.

```elixir
# A tape file that requires gum and glow to be in the $PATH
Require gum
Require glow
```

### Settings

The `Set` command allows you to change global aspects of the terminal, such as
the font settings, window dimensions, and GIF output location.

Setting must be administered at the top of the tape file. Any setting (except
`TypingSpeed`) applied after a non-setting or non-output command will be
ignored.

#### Set Shell

Set the shell with the `Set Shell <shell>` command

```elixir
Set Shell fish
```

#### Set Font Size

Set the font size with the `Set FontSize <number>` command.

```elixir
Set FontSize 10
Set FontSize 20
Set FontSize 40
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/font-size-10.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/font-size-10.gif">
  <img width="600" alt="Example of setting the font size to 10 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-10.gif">
</picture>

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/font-size-20.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/font-size-20.gif">
  <img width="600" alt="Example of setting the font size to 20 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-20.gif">
</picture>

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/font-size-40.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/font-size-40.gif">
  <img width="600" alt="Example of setting the font size to 40 pixels" src="https://stuff.charm.sh/vhs/examples/font-size-40.gif">
</picture>

#### Set Font Family

Set the font family with the `Set FontFamily "<font>"` command

```elixir
Set FontFamily "Monoflow"
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/font-family.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/font-family.gif">
  <img width="600" alt="Example of changing the font family to Monoflow" src="https://stuff.charm.sh/vhs/examples/font-family.gif">
</picture>

#### Set Width

Set the width of the terminal with the `Set Width` command.

```elixir
Set Width 300
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/width.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/width.gif">
  <img width="300" alt="Example of changing the width of the terminal" src="https://stuff.charm.sh/vhs/examples/width.gif">
</picture>

#### Set Height

Set the height of the terminal with the `Set Height` command.

```elixir
Set Height 1000
```
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/height.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/height.gif">
  <img width="300" alt="Example of changing the height of the terminal" src="https://stuff.charm.sh/vhs/examples/height.gif">
</picture>

#### Set Letter Spacing

Set the spacing between letters (tracking) with the `Set LetterSpacing`
Command.

```elixir
Set LetterSpacing 20
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/letter-spacing.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/letter-spacing.gif">
  <img width="600" alt="Example of changing the letter spacing to 20 pixels between characters" src="https://stuff.charm.sh/vhs/examples/letter-spacing.gif">
</picture>

#### Set Line Height

Set the spacing between lines with the `Set LineHeight` Command.

```elixir
Set LineHeight 1.8
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/line-height.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/line-height.gif">
  <img width="600" alt="Example of changing the line height to 1.8" src="https://stuff.charm.sh/vhs/examples/line-height.gif">
</picture>

#### Set Typing Speed

```elixir
Set TypingSpeed 500ms # 500ms
Set TypingSpeed 1s    # 1s
```

Set the typing speed of seconds per key press. For example, a typing speed of
`0.1` would result in a `0.1s` (`100ms`) delay between each character being typed.

This setting can also be overwritten per command with the `@<time>` syntax.

```elixir
Set TypingSpeed 0.1
Type "100ms delay per character"
Type@500ms "500ms delay per character"
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/typing-speed.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/typing-speed.gif">
  <img width="600" alt="Example of using the Type command in VHS" src="https://stuff.charm.sh/vhs/examples/typing-speed.gif">
</picture>

#### Set Theme

Set the theme of the terminal with the `Set Theme` command. The theme value
should be a JSON string with the base 16 colors and foreground + background.

```elixir
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
```

<img alt="Example of changing the theme to Whimsy" src="https://stuff.charm.sh/vhs/examples/theme.gif" width="600" />

You can also set themes by name:

```elixir
Set Theme "Catppuccin Frappe"
```

See the full list by running `vhs themes`, or in [THEMES.md](./THEMES.md).

#### Set Padding

Set the padding (in pixels) of the terminal frame with the `Set Padding`
command.

```elixir
Set Padding 0
```


<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/padding.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/padding.gif">
  <img width="600" alt="Example of setting the padding" src="https://stuff.charm.sh/vhs/examples/padding.gif">
</picture>


#### Set Margin

Set the margin (in pixels) of the video with the `Set Margin` command.

```elixir
Set Margin 60
Set MarginFill "#6B50FF"
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://vhs.charm.sh/vhs-1miKMtNHenh7O4sv76TMwG.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://vhs.charm.sh/vhs-1miKMtNHenh7O4sv76TMwG.gif">
  <img width="600" alt="Example of setting the margin" src="https://vhs.charm.sh/vhs-1miKMtNHenh7O4sv76TMwG.gif">
</picture>

#### Set Window Bar

Set the type of window bar (Colorful, ColorfulRight, Rings, RingsRight) on the terminal window with the `Set WindowBar` command.

```elixir
Set WindowBar Colorful
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://vhs.charm.sh/vhs-4VgviCu38DbaGtbRzhtOUI.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://vhs.charm.sh/vhs-4VgviCu38DbaGtbRzhtOUI.gif">
  <img width="600" alt="Example of setting the margin" src="https://vhs.charm.sh/vhs-4VgviCu38DbaGtbRzhtOUI.gif">
</picture>

#### Set Border Radius

Set the border radius (in pixels) of the terminal window with the `Set BorderRadius` command.

```elixir
# You'll likely want to add a Margin + MarginFill if you use BorderRadius.
Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://vhs.charm.sh/vhs-4nYoy6IsUKmleJANG7N1BH.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://vhs.charm.sh/vhs-4nYoy6IsUKmleJANG7N1BH.gif">
  <img width="400" alt="Example of setting the margin" src="https://vhs.charm.sh/vhs-4nYoy6IsUKmleJANG7N1BH.gif">
</picture>


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

#### Set Loop Offset

Set the offset for when the GIF loop should begin. This allows you to make the
first frame of the GIF (generally used for previews) more interesting.

```elixir
Set LoopOffset 5 # Start the GIF at the 5th frame
Set LoopOffset 50% # Start the GIF halfway through
```

#### Set Cursor Blink

Set whether the cursor should blink. Enabled by default.

```elixir
Set CursorBlink false
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://vhs.charm.sh/vhs-3rMCb80VEkaDdTOJMCrxKy.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://vhs.charm.sh/vhs-3rMCb80VEkaDdTOJMCrxKy.gif">
  <img width="600" alt="Example of setting the cursor blink." src="https://vhs.charm.sh/vhs-3rMCb80VEkaDdTOJMCrxKy.gif">
</picture>

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

Escape single and double quotes with backticks.

```elixir
Type `VAR="Escaped"`
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/type.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/type.gif">
  <img width="600" alt="Example of using the Type command in VHS" src="https://stuff.charm.sh/vhs/examples/type.gif">
</picture>

### Keys

Key commands take an optional `@time` and optional repeat `count` for repeating
the key press every interval of `<time>`.

```
Key[@<time>] [count]
```

#### Backspace

Press the backspace key with the `Backspace` command.

```elixir
Backspace 18
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/backspace.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/backspace.gif">
  <img width="600" alt="Example of pressing the Backspace key 18 times" src="https://stuff.charm.sh/vhs/examples/backspace.gif">
</picture>

#### Ctrl

You can access the control modifier and send control sequences with the `Ctrl`
command.

```elixir
Ctrl+R
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/ctrl.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/ctrl.gif">
  <img width="600" alt="Example of pressing the Ctrl+R key to reverse search" src="https://stuff.charm.sh/vhs/examples/ctrl.gif">
</picture>

#### Enter

Press the enter key with the `Enter` command.

```elixir
Enter 2
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/enter.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/enter.gif">
  <img width="600" alt="Example of pressing the Enter key twice" src="https://stuff.charm.sh/vhs/examples/enter.gif">
</picture>

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

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/arrow.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/arrow.gif">
  <img width="600" alt="Example of pressing the arrow keys to navigate text" src="https://stuff.charm.sh/vhs/examples/arrow.gif">
</picture>

#### Tab

Enter a tab with the `Tab` command.

```elixir
Tab@500ms 2
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/tab.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/tab.gif">
  <img width="600" alt="Example of pressing the tab key twice for autocomplete" src="https://stuff.charm.sh/vhs/examples/tab.gif">
</picture>

#### Space

Press the space bar with the `Space` command.

```elixir
Space 10
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/space.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/space.gif">
  <img width="600" alt="Example of pressing the space key" src="https://stuff.charm.sh/vhs/examples/space.gif">
</picture>

#### Page Up / Down

Press the Page Up / Down keys with the `PageUp` or `PageDown` commands.

```elixir
PageUp 3
PageDown 5
```

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

The `Hide` command instructs VHS to stop capturing frames. It's useful to pause
a recording to perform hidden commands.

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

The `Show` command instructs VHS to begin capturing frames, again. It's useful
after a `Hide` command to resume frame recording for the output.

```elixir
Hide
Type "You won't see this being typed."
Show
Type "You will see this being typed."
```

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://stuff.charm.sh/vhs/examples/hide.gif">
  <source media="(prefers-color-scheme: light)" srcset="https://stuff.charm.sh/vhs/examples/hide.gif">
  <img width="600" alt="Example of typing something while hidden" src="https://stuff.charm.sh/vhs/examples/hide.gif">
</picture>

### Screenshot

The `Screenshot` command captures the current frame (png format).

```elixir
# At any point...
Screenshot examples/screenshot.png
```

### Copy / Paste

The `Copy` and `Paste` copy and paste the string from clipboard.

```elixir
Copy "https://github.com/charmbracelet"
Type "open "
Sleep 500ms
Paste
```


### Source

The `source` command allows you to execute commands from another tape.

```elixir
Source config.tape
```

***

## Continuous Integration

You can hook up VHS to your CI pipeline to keep your GIFs up-to-date with
the official VHS GitHub Action:

> [⚙️ charmbracelet/vhs-action](https://github.com/charmbracelet/vhs-action)

VHS can also be used for integration testing. Use the `.txt` or `.ascii` output
to generate golden files. Store these files in a git repository to ensure there
are no diffs between runs of the tape file.

```elixir
Output golden.ascii
```

## Syntax Highlighting

There’s a tree-sitter grammar for `.tape` files available for editors that
support syntax highlighting with tree-sitter:

> [🌳 charmbracelet/tree-sitter-vhs](https://github.com/charmbracelet/tree-sitter-vhs)

It works great with Neovim, Emacs, and so on!

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
