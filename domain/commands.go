package domain

type (
	CommandSet map[Order]Command

	Command interface {
		Execute() error
	}

	ScaleCommand struct {
		Provider Provider
	}

	ScaleCommands []ScaleCommand

	Launch struct {
		ScaleCommand
	}

	Terminate struct {
		ScaleCommand
	}
)

func (lc *Launch) Execute() error {
	return nil
}

func (lc *Terminate) Execute() error {
	return nil
}
