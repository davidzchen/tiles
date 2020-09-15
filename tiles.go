package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidzchen/tiles/tilespb"
	"github.com/davidzchen/tiles/tmux"
	"google.golang.org/protobuf/encoding/prototext"
)

const (
	startSubcommand    = "start"
	attachSubcommand   = "attach"
	lsSubcommand       = "ls"
	saveSubcommand     = "save"
	helpSubcommand     = "help"
	defaultConfigPath  = "~/.tiles"
	defaultSessionName = "default"
)

var (
	startCmd   = flag.NewFlagSet(startSubcommand, flag.ExitOnError)
	attachCmd  = flag.NewFlagSet(attachSubcommand, flag.ExitOnError)
	lsCmd      = flag.NewFlagSet(lsSubcommand, flag.ExitOnError)
	saveCmd    = flag.NewFlagSet(saveSubcommand, flag.ExitOnError)
	saveOutput = saveCmd.String("o", "", "File to save the current tmux configuration to")
	helpCmd    = flag.NewFlagSet(helpSubcommand, flag.ExitOnError)

	subcommands = []string{
		startSubcommand,
		attachSubcommand,
		lsSubcommand,
		saveSubcommand,
		helpSubcommand,
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected \"save\", \"attach\", \"ls\", or \"save\" subcommands")
		os.Exit(1)
	}

	args := os.Args[2:]
	switch os.Args[1] {
	case startSubcommand:
		runStart(args)
	case attachSubcommand:
		runAttach(args)
	case lsSubcommand:
		runLs(args)
	case saveSubcommand:
		runSave(args)
	case helpSubcommand:
		runHelp(args)
	}
}

func readConfig(path string) (*tilespb.TilesConfig, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read config file: %v", path, err)
	}
	config := &tilespb.TilesConfig{}
	if err := prototext.Unmarshal(in, config); err != nil {
		return nil, fmt.Errorf("%s: failed to parse config file: %v", path, err)
	}
	return config, nil
}

func writeConfig(config *tilespb.TilesConfig, path string) error {
	opts := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	buf, err := opts.Marshal(config)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal config: %v", path, err)
	}
	if err := ioutil.WriteFile(path, buf, 0644); err != nil {
		return fmt.Errorf("%s: failed to write config file: %v", path, err)
	}
	return nil
}

func getSessionName(args []string) string {
	sessionName := defaultSessionName
	if len(args) > 0 {
		sessionName = args[0]
	}
	return sessionName
}

func runStart(args []string) {
	path := defaultConfigPath
	sessionName := getSessionName(args)
	config, err := readConfig(path)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(2)
	}

	var tmuxSession *tilespb.TmuxSession
	for _, s := range config.GetTmuxSession() {
		if s.GetName() == sessionName {
			tmuxSession = s
			break
		}
	}
	if tmuxSession == nil {
		fmt.Printf("%s: tmux session not found in config: %s", path, sessionName)
		os.Exit(2)
	}

	tmux.Start(tmuxSession)
}

func runAttach(args []string) {
	sessionName := getSessionName(args)
	tmux.AttachSession(sessionName)
}

func runLs(args []string) {
	tmux.ListSessions()
}

func runSave(args []string) {
	path := defaultConfigPath
	state, err := tmux.GetState()
	if err != nil {
		fmt.Printf("failed to get tmux state: %v", err)
		os.Exit(3)
	}
	if err := writeConfig(state, path); err != nil {
		fmt.Printf("failed to write tiles config: %v", err)
		os.Exit(4)
	}
}

func runHelp(args []string) {

}
