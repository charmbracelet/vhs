Output renders/jqp.gif
Output renders/jqp.mp4
Set Width 2400
Set Height 1600

Hidden Begin
Type "curl https://dummyjson.com/products | jqp"
Hidden End

Sleep 1s
Enter
Sleep 1s

Type@200ms '[ .products[] | select(.category=="smartphones") ]'
Sleep 1s
Enter

Sleep 1s
Tab

Sleep 1s
Down@25ms 20

Sleep 500ms
Down@25ms 20

Tab

Sleep 500ms
Down@25ms 30

Sleep 1s
Ctrl+S

Sleep 3s

Escape

Sleep 3s
