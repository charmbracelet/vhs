Output examples/settings/set-caption-key-style.gif
Set Width 800
Set Height 225
Set Padding 30

Set FontSize 18
Set TypingSpeed 0.1
Set Theme "Catppuccin Macchiato"

Set Shell zsh
Set CaptionKeyStyle "vim"

CaptionOn

# Set up line editing to demo control keys
Hide
Type "bindkey -e"
Enter
Ctrl+l
Show

Type "echo Goodbye, Planet"
Sleep 0.5s
Enter
Sleep 2s
Ctrl+p
Sleep 1.5s
Ctrl+w
Sleep 1.1s
Type "World"
Enter
Sleep 2s
