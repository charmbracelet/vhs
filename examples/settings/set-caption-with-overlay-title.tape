Output examples/settings/set-caption-with-overlay-title.gif
Set Width 800
Set Height 225
Set Padding 30

Set FontSize 18
Set TypingSpeed 0.1
Set Theme "Catppuccin Macchiato"

Set Shell zsh
Set OverlayAlignment "top-right"
Set OverlayBoxColor "#ffffff"
Set OverlayBoxOpacity 1
Set OverlayFontColor "#000000"
Set OverlayMarginRight 10
Set OverlayMarginVertical 10

CaptionOn

# Set up line editing to demo control keys
Hide
Type "bindkey -e"
Enter
Ctrl+l
Show

Overlay@100s "Demo of Command Line Editing"
Type "echo Hello, Planet"
Sleep 0.5s
Enter
Sleep 2s
Sleep 1s
Ctrl+p
Sleep 1.5s
Ctrl+w
Sleep 1.1s
Type "World"
Enter
Sleep 2s
