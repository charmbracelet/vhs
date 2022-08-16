package dolly

// Run starts a dolly process, evaluates a set of Commands by executing them in
// order, and then cleans up the processes, and compiles the GIF.
func Run(commands []Command) {
	d := New(DefaultDollyOptions())
	defer d.Cleanup()

	for _, command := range commands {
		command.Execute(&d)
	}
}
