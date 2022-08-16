package dolly

func Run(commands []Command) {
	d := New(DefaultDollyOptions())
	defer d.Cleanup()

	for _, command := range commands {
		command.Execute(&d)
	}
}
