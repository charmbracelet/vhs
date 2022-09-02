# VHS

GIFs as code. Record GIFs for terminal applications with a just few lines of code 🎬.

<img width="400" src="./examples/out.gif" alt="Automatic GIF recording with vhs" />

The above example is generated from a single tape file: ([demo.tape](./examples/demo.tape)).

## Tutorial

Type anything into the terminal with the `Type` command.

```
Type "echo 'Hello, world!'"
```

Press the Enter key with the `Enter` command.

```
Enter
```

Press the Backspace key in the terminal with the `Backspace` command.

```
Backspace
```

Wait for a certain amount of time with the `Sleep` command.

```
Sleep 1s
```

Putting it all together...

```
Type "echo 'Hello World'"
Enter
Backspace 5
Sleep 1s
```

Save the above text to a file (`demo.tape`) and generate the GIF with `vhs`:

```bash
vhs < demo.tape
open out.gif
```

<img width="400" src="./examples/tutorial.gif" alt="Tutorial GIF recording with VHS" />

## Commands

The following is an exhaustive list of all the commands that can be used with `vhs`.

The `Set` command allows you to change settings from the defaults in order to style the output GIF how you want it to look.
See the [Settings](#Settings) section for a list of all possible settings.

```
Set <setting> <value>
```

Type allows you to type a string into the terminal at a given typing speed which can be specified with the `@time` modifier.

```
Type[@<time>] "<characters...>"
```

To press special keys such as `Enter`, `Backspace`, and arrow keys, you can use the following commands.
Each of these commands also takes an optional `@time` to specify typing speed with an optional repeat `count`.

```
Enter[@<time>] [count]
Backspace[@<time>] [count]
Down[@<time>] [count]
Left[@<time>] [count]
Right[@<time>] [count]
Up[@<time>] [count]
```

To allow the GIF to keep recording without any user input use the `Sleep` command.
This command continues recording the terminal for the specified duration so you can still view what is happening in the terminal.

```
Sleep <time>
```

To press special characters such as `Ctrl+C` or `Ctrl+L` you can use the `Ctrl+<character>` command.

```
Alt+<character>
Ctrl+<character>
```

## Settings

You can change some settings for your GIFs through the `Set` command.

The `Output` setting allows you to specify the name that the output file will be named.

```
Set Output out.gif
```

You can change certain font settings of the terminal with `FontFamily`, `FontSize` and `LineHeight`.

```
Set FontFamily "SF Mono"
Set FontSize 32
Set LineHeight 1.2
```

You can also change the `Framerate` (the rate at which screenshots are captured).

```
Set Framerate 60
```

Using the `Width` and `Height` settings, you can specify the dimensions of your terminal and output GIF.

```
Set Width 1200
Set Height 600
```

You can also set the padding of your terminal with the `Padding` command.

```
Set Padding 5em
```

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.technology/@charm)
* [Slack](https://charm.sh/slack)

## License

[MIT](https://github.com/charmbracelet/vhs/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
