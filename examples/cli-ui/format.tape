Output format.gif

Set FontSize 32

Set Width 2200
Set Height 400

Hidden Begin
  Type irb Enter
  Type "require 'cli/ui'" Enter
  Type "CLI::UI::StdoutRouter.enable" Enter
  Ctrl+L
Hidden End

Type 'puts CLI::UI.fmt "{{red:Red}} {{green:Green}}"' Sleep 500ms Enter
Sleep 5s
