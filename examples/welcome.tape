# Render the output GIF to demo.gif
Output examples/welcome.gif
Output examples/frames/

# Set up a 1200x600 frame with 46px font size.
Set FontSize 46
Set Height 600
Set Width 1200
Set TypingSpeed 40ms

Type "Welcome to VHS!"
Sleep 0.5s
Ctrl+C
Sleep 0.5s

Type "A tool for writing GIFs as code."
Sleep 0.5s
Left 14
Sleep 0.5s
Type " terminal"
Sleep 1s
Ctrl+E
Sleep 0.5s
Ctrl+U
Sleep 0.5s

Type 'echo "Demo your CLIs and TUIs."'
Sleep 0.5s
Enter
Sleep 0.5s

Type "gum spin -s line -- sleep 2.5"
Sleep 0.5s
Enter

Sleep 3s

Type "clear"
Sleep 0.5s
Enter
Sleep 1s