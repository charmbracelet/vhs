package dolly

import (
	"log"
)

// Run starts a dolly process, evaluates a set of Commands by executing them in
// order, and then cleans up the processes, and compiles the GIF.
func Run(commands []Command) {
	d := New()
	defer d.Cleanup()

	var offset int

	// Apply the setting commands first, so that the options are correctly set.
	for i, command := range commands {
		if command.Type == Set {
			log.Println(command)
			command.Execute(&d)
		} else {
			offset = i
			break
		}
	}

	d.Start()

	// Execute the normal commands on the running instance.
	for _, command := range commands[offset:] {
		log.Println(command)
		command.Execute(&d)
	}
}
