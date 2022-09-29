Output exec.gif

Hidden Begin
  Type "go build -o exec ." Enter
  Type "clear" Enter
Hidden End

Type "EDITOR=nano ./exec" Enter
Sleep 1s

Type@1s "aaaa"
Sleep 2s

Type@1s "e"
Type "Hello, EDITOR!"
Sleep 2s
Ctrl+X
Sleep 1s
Type "n"
Sleep 2s

Hidden Begin
  Type "q"
  Type "rm exec" Enter
Hidden End
