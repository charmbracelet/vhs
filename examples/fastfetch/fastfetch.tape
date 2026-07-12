# Source code for the VHS fastfetch example.
#
# To run:
#
#     vhs < fastfetch.tape

Output examples/fastfetch/fastfetch.gif
Output examples/fastfetch/fastfetch.mp4
Output examples/fastfetch/fastfetch.webm
Output examples/fastfetch/fastfetch.webp
Output examples/fastfetch/frames/

Set TypingSpeed 75ms
Set FontSize 18
Set Width 1300
Set Height 550

Hide
Type "fastfetch() { command fastfetch -c examples/fastfetch/config.jsonc --pipe false; }"
Enter
Type "clear"
Enter
Show

Type "fastfetch"

Sleep 500ms
Enter
Sleep 2s

Type "Welcome to VHS!"
Sleep 1
Space
Type "A tool for generating terminal GIFs from code."

Sleep 5s
