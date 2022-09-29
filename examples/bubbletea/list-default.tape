Output list-default.gif
Set FontSize 16

Hidden Begin
  Type "go build -o list-default ." Enter
  Type "clear" Enter
Hidden End

Type "./list-default" Enter
Sleep 1s

Down@500ms 3
Right@500ms 3

Type "/"
Sleep 1s

Type "nutel" Enter
Sleep 1s
Type "?"
Sleep 1s
Type "?"
Sleep 1s
Escape@500ms 2

Hidden Begin
  Ctrl+C
  Type "rm list-default" Enter
Hidden End
