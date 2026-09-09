# Run from the repository root: vhs examples/sixel.tape
# Requires ttyd with sixel support; no image encoder is needed.
Output examples/sixel.gif

Set Shell bash
Set Width 800
Set Height 400
Set FontSize 18
Set Sixel true

Hide
Type "clear"
Enter
Show
Sleep 500ms

# A red 300x30 pixel block, followed by text.
Type `printf '\033Pq#0;2;100;0;0#0!300~-!300~-!300~-!300~-!300~\033\\\n'; echo after`
Enter
Sleep 1s
Screenshot examples/sixel.png

# Remove the image, then create another to exercise canvas replacement.
Type "clear"
Enter
Sleep 1s
Type `printf '\033Pq#0;2;0;100;0#0!300~-!300~-!300~-!300~-!300~\033\\\n'; echo again`
Enter
Sleep 1s
