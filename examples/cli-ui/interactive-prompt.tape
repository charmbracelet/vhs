Output interactive-prompt.gif

Set FontSize 32

Set Width 2200
Set Height 500

Hide
Type irb Enter
Type "require 'cli/ui'" Enter
Type "CLI::UI::StdoutRouter.enable" Enter
Ctrl+L
Show

Type "CLI::UI.ask('What do you use?', options: %w(rails go ruby python))"
Sleep 500ms
Enter
Sleep 1s

Down@500ms 3
Sleep 1s
Up@500ms 2
Sleep 1s

Enter

Sleep 3s
