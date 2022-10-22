Output status-widget.gif

Set FontSize 32

Set Width 2200
Set Height 400

Hide
Type "irb --noautocomplete"
Enter
Type "require 'cli/ui'"
Enter
Type "CLI::UI::StdoutRouter.enable"
Enter
Ctrl+L
Show

Type 'CLI::UI::Spinner.spin("Building: {{@widget/status:1:2:3:4}}") { |spinner| sleep 3 }'
Sleep .5
Enter
Sleep 5
