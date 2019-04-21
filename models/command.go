package models

type Command struct {
	AddURL, RedirectionURLOfAddURL, RemoveURL, Port string
	ListRedirections, Help                          bool
	Args                                            []string
}

func (cmd *Command) IsHelpCommand() bool {
	return cmd.Help
}

func (cmd *Command) IsListRedirectionCommand() bool {
	return cmd.ListRedirections
}

func (cmd *Command) IsAddURLCommand() bool {
	return len(cmd.RedirectionURLOfAddURL) > 0 && len(cmd.Args) > 0 && cmd.Args[0] == "configure"
}

func (cmd *Command) IsRemoveURLCommand() bool {
	return len(cmd.RemoveURL) > 0
}

func (cmd *Command) IsStartServerInPort() bool {
	return len(cmd.Port) > 0 && len(cmd.Args) > 0 && cmd.Args[0] == "run"
}
