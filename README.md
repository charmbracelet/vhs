# Dolly

Programmatically record GIFs on the terminal with Go 🎬

## Usage

```go
// First, create a new `dolly` with any options you want.
d := dolly.New(dolly.WithOutput("demo.gif"), dolly.WithFontSize(32))

// Next, `defer` the `Cleanup` function to tear down the spawned processes, remove
// the generated frames, and compile the GIF.
defer d.Cleanup()

// Type anything you want (with the desired delay per keystroke):
d.Type("echo 'Hello, Demo!'", dolly.WithSpeed(50))
d.Enter()

// You can sleep at any point in the process and the delay will show up in your
// video because `dolly` captures frames automatically in the background.
// Our demo.gif will have a one second delay at the end of the GIF
time.Sleep(time.Second)
```

<img width="400" src="./demo.gif" alt="Automatic GIF recording" />

Update the `main.go` file to perform the actions you need to screenshot and
then run the file. This will automatically spawn
[`ttyd`](https://github.com/tsl0922/ttyd) (which must be installed) and then
perform the actions via [`go-rod`](https://github.com/go-rod/rod).

```bash
go run .
```
