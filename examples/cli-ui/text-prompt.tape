Output text-prompt.ascii
Output text-prompt.gif

Set FontSize 32

Set Width 2200
Set Height 500

Hide
Type irb Enter
Type "require 'cli/ui'" Enter
Type "CLI::UI::StdoutRouter.enable" Enter
Ctrl+L
Show

Type "CLI::UI.ask('Is CLI UI Awesome?', default: 'It is great!')" Sleep 500ms Enter
Sleep 1s
Type "I love it!" Sleep 1s Enter
Sleep 3s
