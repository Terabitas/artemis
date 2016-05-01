package domain

type (
	CommandSet map[Order]Command

	Command interface {
		Execute() error
	}

	CMDError struct {
		Code    int
		Message string
	}

	BaseCommand struct {
		Provider Provider
		State    CommandState
		Error    *CMDError
	}

	BaseCommands []BaseCommand

	Launch struct {
		BaseCommand
	}

	Terminate struct {
		BaseCommand

		NodeID ID
	}
)

func (lc *Launch) Execute() error {

	// API call

	return nil
}

func (lc *Terminate) Execute() error {

	// API call

	return nil
}
