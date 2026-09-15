Output examples/settings/rows-columns.gif

Set FontSize 32
Set Rows 10
Set Columns 40

# A count from 1 to 40 as the first line proves the column count,
# and it plus 9 more numbered lines below it fill the frame exactly
# top to bottom, proving all 10 rows. The last line skips its
# trailing newline so the returning prompt lands on the same line
# instead of scrolling everything up by one.
Hide
Type "clear; echo 1234567890123456789012345678901234567890; echo Line 2; echo Line 3; echo Line 4; echo Line 5; echo Line 6; echo Line 7; echo Line 8; echo Line 9; printf 'Line 10'"
Show

Sleep 1s
Enter

Sleep 2s
