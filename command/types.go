package command

type (
	Command interface {
		Execute() error
	}

	Bus interface {
		// Add new command to bus for execution
		Add(Command) error

		// Run service and execute commands in FIFO
		Run() error
	}

	Launch struct {
	}

	Terminate struct {
	}

	Restart struct {
	}
)
