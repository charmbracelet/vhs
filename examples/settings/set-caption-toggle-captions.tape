Output examples/settings/set-caption-toggle-captions.gif
Set Width 800
Set Height 225
Set Padding 30

Set FontSize 18
Set TypingSpeed 0.1
Set Theme "Catppuccin Macchiato"

Set OverlayAlignment "top-right"
Set OverlayBoxColor "#ffffff"
Set OverlayBoxOpacity 1
Set OverlayFontColor "#000000"
Set OverlayMarginRight 10
Set OverlayMarginVertical 10

Overlay@100s "Selective Captioning"

CaptionOn
Type "echo captions are enabled"
Enter
Sleep 2s

CaptionOff
Type "echo captions are disabled"
Enter
Sleep 2s

CaptionOn
Type "echo captions are re-enabled"
Sleep 2s
