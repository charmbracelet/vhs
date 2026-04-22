Output examples/cli-ui/nested-frames.gif

Set FontSize 32

Set Width 2000
Set Height 750

Hide
Type "irb --noautocomplete"
Enter
Type "require 'cli/ui'"
Enter
Type "CLI::UI::StdoutRouter.enable"
Enter
Ctrl+L
Show

Type "CLI::UI::Frame.open('Frame 1') do"
Enter
Type "  CLI::UI::Frame.open('Frame 2') { puts 'inside frame 2' }"
Enter
Type "  puts 'inside frame 1'"
Enter
Type "end"

Sleep 1

Enter

Sleep 3
