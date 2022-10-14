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
#   Sleep <time>                    Sleep for a set amount of <time>
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
Output ./renders/demo.gif

# Set up a 1200x600 frame with 42px font size.
Set FontSize 42
Set Height 600
Set Width 1200

Hide
Type "Welcome to VHS!"
Ctrl+C
Show

Sleep 1

Type "VHS let's you write terminal GIFs as code."
Sleep .5
Ctrl+C
Sleep .5

Type "Let's take it for a spin."
Sleep .5
Ctrl+C
Sleep .5

Type "gum spin -s line -- sleep 5"
Enter
Sleep 5

# vhs < demo.tape
