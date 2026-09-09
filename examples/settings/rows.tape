Output examples/settings/rows.gif

Set Width 475
Set Rows 8
Set FontSize 32

# Eight lines of "Line N", filling the frame edge to edge: exactly
# 8 lines for 8 rows, no more, no less. The last one skips its
# trailing newline so the returning prompt lands on the same line
# instead of scrolling everything up by one.
Hide
Type "clear; echo Line 1; echo Line 2; echo Line 3; echo Line 4; echo Line 5; echo Line 6; echo Line 7; printf 'Line 8'"
Show

Sleep 1s
Enter

Sleep 2s
