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
	configPathFlag = flag.String("c", defaultConfigPath, "Path to the .tiles config file. Set to "+defaultConfigPath+" by default.")
	saveCmd        = flag.NewFlagSet(saveSubcommand, flag.ExitOnError)
	saveOutputFlag = saveCmd.String("o", "", "File to save the current tmux configuration to")

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
	flag.Parse()

	switch os.Args[1] {
	case startSubcommand:
		runStart()
	case attachSubcommand:
		runAttach()
	case lsSubcommand:
		runLs()
	case saveSubcommand:
		runSave()
	case helpSubcommand:
		runHelp()
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

// getSessionName returns the session name passed in from the command line.
func getSessionName(args []string) string {
	sessionName := defaultSessionName
	if len(args) > 1 {
		sessionName = args[1]
	}
	return sessionName
}

func runStart() {
	sessionName := getSessionName(flag.Args())
	config, err := readConfig(*configPathFlag)
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
		fmt.Printf("%s: tmux session not found in config: %s", *configPathFlag, sessionName)
		os.Exit(2)
	}

	tmux.Start(tmuxSession)
}

func runAttach() {
	sessionName := getSessionName(flag.Args())
	tmux.AttachSession(sessionName)
}

func runLs() {
	tmux.ListSessions()
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

func runSave() {
	saveCmd.Parse(os.Args[2:])
	state, err := tmux.GetState()
	if err != nil {
		fmt.Printf("failed to get tmux state: %v", err)
		os.Exit(3)
	}
	if err := writeConfig(state, *saveOutputFlag); err != nil {
		fmt.Printf("failed to write tiles config: %v", err)
		os.Exit(4)
	}
}

func runHelp() {
	flag.PrintDefaults()
}
