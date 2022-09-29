Output split-editors.gif
Set FontSize 16

Hidden Begin
  Type "go build -o split-editors ." Enter
  Type "clear" Enter
Hidden End

Type "./split-editors" Enter
Sleep 2s

Type "Hello, there!"
Sleep 1s
Tab
Sleep 1s
Type "Hi!"
Sleep 1s
Enter
Type "How are you?"
Sleep 1s
Tab
Sleep 1s
Ctrl+U
Sleep 1s
Type "I'm good! Thanks for asking :)"
Sleep 500ms
Enter
Sleep 1s
Ctrl+N
Sleep 1s
Tab@500ms 2
Type "Hello, world!"
Sleep 2s

Hidden Begin
  Escape
  Type "rm split-editors" Enter
Hidden End
