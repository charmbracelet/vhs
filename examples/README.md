# Examples

### Glow

<img alt="Glow demo with VHS" src="./glow/vhs-glow.gif" width="600" />

```elixir
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
