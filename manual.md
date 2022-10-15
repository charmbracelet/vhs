# Manual

## Introduction

VHS let's you write terminal GIFs as code.

Below is an example of a `.tape` file which types and executes a command:

```haskell
Output demo.gif

Set FontSize 22
Set Height 600
Set Width 1200

Type "echo 'Hello, VHS!'"
Enter

Sleep 1
```

To render the output of the script, run:

```bash
vhs < demo.tape
```

A file called `demo.gif` (Output) will be generated.

## Reference

A tape file is made up of a sequence of commands.

There are several types of commands in the *Cassette* language, such as:

* Output
* Set
* Sleep
* Type
* Enter
* Backspace
* Tab
* Up
* Down
* Left
* Right

**Output**

The Output command instructs VHS where to save the output of the recording to
file names with the extension .gif, .webm, .mp4 will have those respective file
types. Multiple outputs can be specified to produce multiple outputs types.

```haskell
Output demo.gif
Output demo.webm
Output demo.mp4
```

**Set**

The Set command allows VHS to adjust settings in the terminal, such as fonts,
dimensions, and themes.

```haskell
Set FontSize 22
Set FontFamily "SF Mono"
Set Height 600
Set Width 1200
Set LetterSpacing 1
Set LineHeight 1.2
Set TypingSpeed 100ms
Set Theme { ... }
Set Padding 5em
Set Framerate 60
```

**Sleep**

The Sleep command instructs VHS to continue recording frames without
interacting with the terminal. The default units are in seconds.

```haskell
Sleep 5
Sleep 10
Sleep 0.5
```

**Type**

The Type command instructs VHS to input a string of characters into the
terminal. The characters will be input with a small delay defined by the
TypingSpeed, which defaults to 100 milliseconds after each keypress. This can
be adjusted by setting the default TypingSpeed to the desired delay between
characters or adjusted on a per-command basis with the `@<time>` syntax.

```haskell
Type "Hello, VHS!"
Type@500ms "500ms delay between each character."
```


**Enter**

The Enter command instructs VHS to press the Enter key.
The command accepts an optional speed (delay between repetitions) and an
optional repeat count.

```haskell
Enter[@<time>] [repeat]
```

**Backspace**

The Backspace command instructs VHS to press the Backspace key. The command
accepts an optional speed and optional repeat count.

```haskell
Backspace[@<time>] [repeat]
```

**Tab**

The Tab command instructs VHS to press the Tab key. The command accepts an
optional speed and optional repeat count.

```haskell
Tab[@<time>] [repeat]
```

**Up**

The Up command instructs VHS to press the Up arrow key. The command accepts an
optional speed and optional repeat count.

```haskell
Up[@<time>] [repeat]
```

**Down**

The Down command instructs VHS to press the Down arrow key. The command accepts
an optional speed and optional repeat count.

```haskell
Down[@<time> [repeat]
```

**Left**

The Left command instructs VHS to press the Left arrow key. The command accepts
an optional speed and optional repeat count.

```haskell
Left[@<time>] [repeat]
```

**Right**

The Right command instructs VHS to press the Right arrow key. The command
accepts an optional speed and optional repeat count.


```haskell
Right[@<time>] [repeat]
```

