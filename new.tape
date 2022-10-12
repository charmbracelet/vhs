# VHS documentation
#
# Output:
#   Output <path>.gif
#   Output <path>.mp4
#   Output <path>.webm
#
# Settings:
#   Set FontSize <Number>
#   Set FontFamily <String>
#   Set Height <Number>
#   Set Width <Number>
#   Set LetterSpacing <Float>
#   Set LineHeight <Float>
#   Set Theme <String>
#   Set Padding <Number>[em|px]
#   Set Framerate <Number>
#
# Sleep:
#   Sleep <time>
#
# Type:
#   Type[@<time>] "<characters>"
#
# Keys:
#   Backspace[@<time>] [number]
#   Down[@<time>] [number]
#   Enter[@<time>] [number]
#   Space[@<time>] [number]
#   Tab[@<time>] [number]
#   Left[@<time>] [number]
#   Right[@<time>] [number]
#   Up[@<time>] [number]
#   Down[@<time>] [number]
#   Ctrl+<key>

Output ./renders/demo.gif

Set FontSize 42
Set Height 600
Set Width 1200

Type "Welcome to VHS!"

Sleep 1s

Ctrl+C
