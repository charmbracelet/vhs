# onefetch

<img width="600" src="onefetch.gif" />

```
# Source code for the VHS onefetch example.
#
# To run:
#
#     vhs < onefetch.tape

Output onefetch.gif
Output onefetch.mp4
Output onefetch.webm

Set TypingSpeed 75ms
Set FontSize 22
Set Width 1350
Set Height 900

Type "onefetch --ascii-input "$(cat vhs.ascii)" --ascii-colors 5"

Sleep 500ms

Enter

Sleep 2s

Type "Welcome to VHS!"

Sleep 1

Space

Type "A tool for generating terminal GIFs from code."

Sleep 5s
```
