Output paginator.gif
Set FontSize 12
Set Padding 3em

Hidden Begin
  Type "go build -o paginator ." Enter
  Type "clear" Enter
Hidden End

Type "./paginator" Enter
Sleep 1s

Right@500ms 10
Sleep 2s
Left@500ms 10
Sleep 2s

Hidden Begin
  Ctrl+C
  Type "rm paginator" Enter
Hidden End
