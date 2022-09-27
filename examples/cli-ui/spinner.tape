Output spinner.gif

Set FontSize 32

Set Width 2200
Set Height 400

Hidden Begin
  Type irb Enter
  Type "require 'cli/ui'" Enter
  Type "CLI::UI::StdoutRouter.enable" Enter
  Ctrl+L
Hidden End

Type "CLI::UI::Spinner.spin('Spinning...') { sleep 3 }" Sleep 500ms Enter
Sleep 5s
