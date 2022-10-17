# VHS documentation
#
# Output:
#   Output <path>.gif               Create a GIF output at the given <path>
#   Output <path>.mp4               Create an MP4 output at the given <path>
#   Output <path>.webm              Create a WebM output at the given <path>
#
# Settings:
#   Set FontSize <Number>           Set the font size of the terminal
#   Set FontFamily <String>         Set the font family of the terminal
#   Set Height <Number>             Set the height of the terminal
#   Set Width <Number>              Set the width of the terminal
#   Set LetterSpacing <Float>       Set the font letter spacing (tracking)
#   Set LineHeight <Float>          Set the font line height
#   Set Theme <String>              Set the theme of the terminal (JSON)
#   Set Padding <Number>[em|px]     Set the padding of the terminal
#   Set Framerate <Number>          Set the framerate of the recording
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Backspace[@<time>] [number]     Press the Backspace key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output


# Render the output GIF to demo.gif
Output ./examples/demo.gif

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

# vhs < demo.tape
