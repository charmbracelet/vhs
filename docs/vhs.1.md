% VHS(1) Version 0.2 | "Video Home System" Documentation

NAME
====

**vhs** — Terminal GIFs as code.

SYNOPSIS
========

  **vhs** file.tape

  **vhs** [**-h**|**--help**]

  **vhs** [**-v**|**--version**]

DESCRIPTION
===========

VHS let's you write terminal GIFs as code.
VHS reads `.tape` files and renders GIFs (videos).
A tape file is a script made up of commands describing what actions to perform
in the render.

The following is a list of all possible commands in VHS:

* `Output` &lt;path&gt;.(gif|webm|mp4)
* `Set` &lt;setting&gt; value
* `Sleep` &lt;time&gt;
* `Type` "&lt;string&gt;"
* `Ctrl`+&lt;key&gt;
* `Backspace` [repeat]
* `Down` [repeat]
* `Enter` [repeat]
* `Left` [repeat]
* `Right` [repeat]
* `Tab` [repeat]
* `Up` [repeat]
* `Hide`
* `Show`

OUTPUT
======

The Output command instructs VHS where to save the output of the recording.

File names with the extension `.gif`, `.webm`, `.mp4` will have the respective
file types.

SETTINGS
========

The Set command allows VHS to adjust settings in the terminal,
such as fonts, dimensions, and themes.

The following is a list of all possible setting commands in VHS:

* Set `FontSize` &lt;number&gt;
* Set `FontFamily` &lt;string&gt;
* Set `Height` &lt;number&gt;
* Set `Width` &lt;number&gt;
* Set `LetterSpacing` &lt;float&gt;
* Set `LineHeight` &lt;float&gt;
* Set `TypingSpeed` &lt;time&gt;
* Set `Theme` &lt;json&gt;
* Set `Padding` &lt;number&gt;
* Set `Framerate` &lt;number&gt;

BUGS
====

See GitHub Issues: <https://github.com/charmbracelet/vhs/issues>

AUTHOR
======

Charm <vt100@charm.sh>

SEE ALSO
========

**vhs(1)**
