Set Output ./examples/ctrl.gif
Set FontSize 42
Set Height 250

Type@25ms "echo 'Press Ctrl+U'"
Sleep 250ms
Ctrl+U
Sleep 1s
Type@25ms "echo 'Press Ctrl+C'"
Sleep 500ms
Ctrl+C
Sleep 1s
