Output examples/settings/set-caption-with-overlay-multiline.gif
Set Width 800
Set Height 225
Set Padding 30

Set FontSize 18
Set TypingSpeed 0.1
Set Theme "Catppuccin Macchiato"

Set Shell zsh
Set OverlayBoxOpacity 0
Set OverlayAlignment "top-right"

CaptionOn

# Set up line editing to demo control keys
Hide
Type "bindkey -e"
Enter
Ctrl+l
Show

Type "echo Hello, Planet"
Sleep 0.5s
Enter
Sleep 2s
Overlay@10s "Let's now edit\Nthe prior command\Nto show the world!"
Sleep 1s
Ctrl+p
Sleep 1.5s
Ctrl+w
Sleep 1.1s
Type "World"
Enter
Sleep 2s
