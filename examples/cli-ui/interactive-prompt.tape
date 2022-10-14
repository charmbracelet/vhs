Output interactive-prompt.gif

Set FontSize 32

Set Width 2200
Set Height 500

Hide
Type irb
Enter
Type "require 'cli/ui'"
Enter
Type "CLI::UI::StdoutRouter.enable"
Enter
Ctrl+L
Show

Type "CLI::UI.ask('What do you use?', options: %w(rails go ruby python))"
Sleep .5
Enter
Sleep 1

Down@.5 3
Sleep 1
Up@.5 2
Sleep 1

Enter

Sleep 3
