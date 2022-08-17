package dolly

import (
	"log"
)

// Run starts a dolly process, evaluates a set of Commands by executing them in
// order, and then cleans up the processes, and compiles the GIF.
func Run(commands []Command) {
	var settingCommands, normalCommands []Command

	// Split the commands into two groups: setting commands and normal commands.
	for _, command := range commands {
		if command.Type == Set {
			settingCommands = append(settingCommands, command)
		} else {
			normalCommands = append(normalCommands, command)
		}
	}

	d := New()

	// Apply the setting commands first, so that the options are correctly set.
	for _, command := range settingCommands {
		log.Println(command.Args)
		command.Execute(&d)
	}

	d.Start()
	defer d.Cleanup()

	// Execute the normal commands on the running instance.
	for _, command := range normalCommands {
		log.Println(command)
		command.Execute(&d)
	}
}
