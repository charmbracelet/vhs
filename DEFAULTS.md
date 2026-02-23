# VHS Default Options

This document lists the default values for VHS configuration options.

## Terminal Settings

| Option | Default Value | Notes |
|--------|---------------|-------|
| FontFamily | JetBrains Mono, DejaVu Sans Mono, Menlo, Bitstream Vera Sans Mono, Inconsolata, Roboto Mono, Hack, Consolas, ui-monospace, monospace | Comma-separated font stack with Apple Symbols fallback |
| FontSize | 22 | |
| LetterSpacing | 1.0 | |
| LineHeight | 1.0 | |
| TypingSpeed | 50ms | |
| Shell | bash (Unix) / cmd (Windows) | Valid: bash, zsh, fish, powershell, pwsh, cmd, nu, osh, xonsh |
| CursorBlink | true | |

## Video Settings

| Option | Default Value | Notes |
|--------|---------------|-------|
| Width | 1200 | Pixels |
| Height | 600 | Pixels |
| Padding | 60 | Pixels |
| Framerate | 50 | Frames per second |
| PlaybackSpeed | 1.0 | Multiplier (0.5 = half speed, 2.0 = double speed) |
| MaxColors | 256 | For GIF output |
| LoopOffset | 0 | Percentage (0-100%) for GIF loop start position |

## Window Styling

| Option | Default Value | Notes |
|--------|---------------|-------|
| Margin | 0 | Pixels |
| MarginFill | #171717 | Hex color or image path |
| WindowBar | *(empty)* | Valid: Colorful, ColorfulRight, Rings, RingsRight |
| WindowBarSize | 30 | Pixels (when WindowBar is enabled) |
| WindowBarColor | #171717 | |
| BorderRadius | 0 | Pixels |
| BackgroundColor | #171717 | |

## Wait Command Settings

| Option | Default Value | Notes |
|--------|---------------|-------|
| WaitTimeout | 15s | Must be positive |
| WaitPattern | `>$` | Regex pattern |

## Default Theme Colors

| Option | Default Value | Notes |
|--------|---------------|-------|
| Background | #171717 | |
| Foreground | #dddddd | |
| Cursor | #dddddd | |
| CursorAccent | #171717 | |
| Black | #282a2e | ANSI 0 |
| BrightBlack | #4d4d4d | ANSI 8 |
| Red | #D74E6F | ANSI 1 |
| BrightRed | #FE5F86 | ANSI 9 |
| Green | #31BB71 | ANSI 2 |
| BrightGreen | #00D787 | ANSI 10 |
| Yellow | #D3E561 | ANSI 3 |
| BrightYellow | #EBFF71 | ANSI 11 |
| Blue | #8056FF | ANSI 4 |
| BrightBlue | #9B79FF | ANSI 12 |
| Magenta | #ED61D7 | ANSI 5 |
| BrightMagenta | #FF7AEA | ANSI 13 |
| Cyan | #04D7D7 | ANSI 6 |
| BrightCyan | #00FEFE | ANSI 14 |
| White | #bfbfbf | ANSI 7 |
| BrightWhite | #e6e6e6 | ANSI 15 |
