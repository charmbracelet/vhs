Output split.gif

Require tmux

Hide
  Type "tmux -f /dev/null -L test new-session -- bash" Enter
  Type "tmux split-window -d -h -- bash && \" Enter
  Type "tmux set status && \" Enter
  Type 'tmux setw pane-border-style "fg=0" && \' Enter
  Type 'tmux setw pane-active-border-style "fg=0"' Enter
  Sleep 0.5
  Ctrl+L
  Sleep 1
Show

Type "Window 1"
Sleep 1

# tmux bind to switch to next pane
Ctrl+B
Type o

Sleep 1
Type "Window 2"
Sleep 1
