# Render the output GIF to demo.gif
Output examples/welcome.gif
Output examples/frames/

# Set up a 1200x600 frame with 46px font size.
Set FontSize 46
Set Height 600
Set Width 1200

Type "Welcome to VHS!"
Sleep 1s
Ctrl+C
Sleep 1s

Type "A tool for writing GIFs as code."
Sleep 1s
Left 14
Sleep 0.5s
Type " terminal"
Sleep 1s
Ctrl+E
Sleep 1s
Ctrl+U
Sleep 1s

Type 'echo "Demo your CLIs and TUIs."'
Sleep 1s
Enter
Sleep 1s

Type "gum spin -s line -- sleep 2.5"
Sleep 1s
Enter

Sleep 4s

Type "clear"
Sleep 0.5s
Enter
Sleep 1s