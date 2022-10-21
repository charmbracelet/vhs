# Examples

### Gum

Example of recording a demo of [Gum](https://github.com/charmbracelet/gum)
with VHS.

#### Gum File

<img alt="gum file demo with VHS" src="./gum/file.gif" width="600" />

```
Output file.gif

Type "gum file ./src"
Sleep 1s
Enter
Sleep 2s
Down@500ms 4
Up@500ms 1
Sleep 1s
Enter
Sleep 1s
Down@500ms 4
Sleep 1s
Up@500ms 2
Sleep 2s
```

#### Gum Pager

<img alt="gum pager demo with VHS" src="./gum/pager.gif" width="600" />

```
Output pager.gif

Set Padding 2em
Set FontSize 16
Set Height 600

Type "gum pager < ~/src/gum/README.md --border normal"
Sleep 1s
Enter
Sleep 2s
Down@25ms 40
Sleep 1s
Up@25ms 30
Sleep 1s
Down@25ms 20
Sleep 3s
```

#### Gum Table

<img alt="gum table demo with VHS" src="./gum/table.gif" width="600" />

```
Output table.gif

Type "gum table < superhero.csv -w 2,12,5,6,6,8,4,20 --height 10"
Enter
Sleep 1s
Down@200ms 10
Sleep 1s
Down@200ms 10
Sleep 1s
Up@200ms 10
Sleep 1s
Enter
Sleep 3s
```

### GitHub CLI

Examples recorded with VHS for the GitHub CLI (`gh`):

#### Issues

<img alt="Simple gh issue demo" src="./gh-cli/gh-issue.gif" width="600" />

#### Pull Requests

<img alt="Simple gh pr demo" src="./gh-cli/gh-pr.gif" width="600" />

### Bubble Tea

Examples recorded with VHS for Bubble Tea.

* [GIFS Renders](https://github.com/charmbracelet/bubbletea/tree/master/examples)
* [Tape Files](./bubbletea)

### jqp

Example of recording a demo of [`jqp`](https://github.com/noahgorstein/jqp)
with VHS.

<img alt="Simple jqp demo with VHS" src="./jqp/jqp.gif" width="600" />

### Glow

Example of recording a demo of [Glow](https://github.com/charmbracelet/glow)
with VHS.

#### Glow Simple

<img alt="Simple glow demo with VHS" src="./glow/glow-simple.gif" width="600" />

```
Output glow-simple.ascii
Output glow-simple.gif

Set Width 1000
Set Height 1000

Type "glow"
Enter
Sleep 1s
Enter
Sleep 1s
Escape
Sleep 1s
Type "q"
Sleep 1s
```

#### Glow

<img alt="Glow demo with VHS" src="./glow/vhs-glow.gif" />

```
Output vhs-glow.gif
Output glow.ascii

Set Width 1600
Set Height 1040

Sleep 1s

Type "glow"

Sleep 100ms

Enter

Sleep 1s

Hide
Tab
Type "/artichoke"
Enter
Down 2
Show

Sleep 0.5s

Down 20

Hide
Escape
Type "l"
Down 5
Show

Sleep 1s
Up@400ms 5

Hide
Type "/ulysses"
Enter
Show

Sleep 0.5s

Down@200ms 20

Hide
Escape
Type "/"
Show

Sleep 0.5s

Type@500ms "todo"
Sleep 1

Hide
Escape
Type "/ulysses"
Enter
Show

Sleep 0.5s

Type@750ms "????"

Hide
Escape
Type "/artichoke"
Enter
Type "m"
Ctrl+A
Right 4
Show

Sleep 1s
Type@250ms "Tasty "
Sleep 1s

Hide
Escape
Down 5
Type "m"
Ctrl+U
Show

Sleep 1s
Type@150ms "Your new internet thing"
Sleep 3s

Hide
Ctrl+C
Show
```
