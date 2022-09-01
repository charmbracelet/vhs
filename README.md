# VHS

GIFs as code. Record GIFs for terminal applications with a just few lines of code 🎬.

<img width="400" src="./out.gif" alt="Automatic GIF recording with vhs" />

The above example is generated from a single vhs file: ([demo.vhs](./examples/demo.vhs)).

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

Save the above text to a file (`demo.vhs`) and generate the GIF with `vhs`:

```bash
vhs demo.vhs
open out.gif
```

## Commands

The following is an exhaustive list of all the commands that can be used with `vhs`.

```
Backspace[@<time>] [count]
Enter[@<time>] [count]
Sleep <time>
Space[@<time>] [count]
Type[@<time>] <characters...>

Down[@<time>] [count]
Left[@<time>] [count]
Right[@<time>] [count]
Up[@<time>] [count]

Alt+<character>
Ctrl+<character>
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
