# Examples

### Gum

Example of recording a demo of [Gum](https://github.com/charmbracelet/gum)
with VHS.

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

<br />
<br />

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

<br />
<br />


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

### Bubble Tea

Examples recorded with VHS for Bubble Tea.

* [GIFS Renders](https://github.com/charmbracelet/bubbletea/tree/master/examples)
* [Tape Files](./bubbletea)

### Glow

Example of recording a demo of [Glow](https://github.com/charmbracelet/glow)
with VHS.

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

<br />
<br />

<img alt="Glow demo with VHS" src="./glow/vhs-glow.gif" />

```
Output vhs-glow.gif
Output glow.ascii

Set Width 1600
Set Height 1040

Sleep 1

Type glow

Sleep .1

Enter

Sleep 1

Hide
Tab
Type "/artichoke"
Enter
Down 2
Show

Sleep .5

Down 20

Hide
Escape
Type l
Down 5
Show

Sleep 1
Up@.4 5

Hide
Type "/ulysses"
Enter
Show

Sleep .5

Down@.2 20

Hide
Escape
Type "/"
Show

Sleep .5

Type@.5 todo
Sleep 1

Hide
Escape
Type "/ulysses"
Enter
Show

Sleep .5

Type@.75 "????"

Hide
Escape
Type "/artichoke"
Enter
Type "m"
Ctrl+A
Right 4
Show

Sleep 1
Type@.25 "Tasty "
Sleep 1

Hide
Escape
Down 5
Type "m"
Ctrl+U
Show

Sleep 1
Type@.15 "Your new internet thing"
Sleep 3

Hide
Ctrl+C
Show```
