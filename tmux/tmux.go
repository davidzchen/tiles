package tmux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/davidzchen/tiles/tilespb"
	"github.com/pkg/errors"
)

func Start(session *tilespb.TmuxSession) error {
	has, err := HasSession(session.GetName())
	if err != nil {
		return err
	}
	if has {
		if err := AttachSession(session.GetName()); err != nil {
			return err
		}
	}

	if err := NewSession(session.GetName()); err != nil {
		return err
	}
	for i, window := range session.GetWindow() {
		if err := NewWindow(session.GetName(), window.GetName(), window.GetDir(), i); err != nil {
			return err
		}
	}
	if err := SelectWindow(session.GetName(), 0); err != nil {
		return err
	}
	return AttachSession(session.GetName())
}

func quote(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

var runAndGetStderr = func(cmd *exec.Cmd) (string, error) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stderr.String()
	return output, err
}

var runAndGetStdout = func(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	output := stdout.String()
	return output, err
}

var run = func(cmd *exec.Cmd, errMsg string) error {
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, errMsg)
	}
	return nil
}

func HasSession(sessionName string) (bool, error) {
	if sessionName == "" {
		return false, fmt.Errorf("sessionName cannot be empty")
	}
	cmd := exec.Command("tmux", "has-session", "-t", quote(sessionName))
	output, err := runAndGetStderr(cmd)
	if err != nil {
		if strings.HasPrefix(output, "can't find session") {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed while checking for tmux session %q", sessionName)
	}
	return true, nil
}

func NewSession(sessionName string) error {
	if sessionName == "" {
		return fmt.Errorf("sessionName cannot be empty")
	}
	cmd := exec.Command("tmux", "new-session", "-d", "-s", quote(sessionName))
	return run(cmd, fmt.Sprintf("failed to start tmux session %q", sessionName))
}

func AttachSession(sessionName string) error {
	if sessionName == "" {
		return fmt.Errorf("sessionName cannot be empty")
	}
	cmd := exec.Command("tmux", "-2", "attach-session", "-t", quote(sessionName))
	return run(cmd, fmt.Sprintf("failed to attach to tmux session %q", sessionName))
}

func ListSessions() error {
	cmd := exec.Command("tmux", "list-sessions")
	return run(cmd, "failed to list tmux sessions")
}

func stateToProto(output string) (*tilespb.TilesConfig, error) {
	sessionMap := make(map[string]*tilespb.TmuxSession)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("mailformed tmux output line: %s\noutput:\n%s", line, output)
		}
		window := &tilespb.TmuxSession_Window{
			Name: parts[1],
			Dir:  parts[2],
		}
		sessionName := parts[0]
		session, exists := sessionMap[sessionName]
		if !exists {
			session = &tilespb.TmuxSession{
				Name:   sessionName,
				Window: []*tilespb.TmuxSession_Window{window},
			}
			sessionMap[sessionName] = session
		} else {
			session.Window = append(session.Window, window)
		}
	}

	var config *tilespb.TilesConfig
	for _, session := range sessionMap {
		config.TmuxSession = append(config.TmuxSession, session)
	}
	return config, nil
}

func GetState() (*tilespb.TilesConfig, error) {
	cmd := exec.Command("tmux", "list-windows", "-a", "-F", quote("#S:#W:#{pane_current_path}"))
	output, err := runAndGetStdout(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dump tmux windows")
	}
	return stateToProto(output)
}

func NewWindow(sessionName, windowName, directory string, windowId int) error {
	if sessionName == "" {

	}
	if windowName == "" {

	}
	if directory == "" {

	}
	cmd := exec.Command(
		"tmux",
		"new-window",
		"-c", quote(directory),
		"-t", fmt.Sprintf("\"%s:%d\"", sessionName, windowId),
		"-n", quote(windowName))
	return run(cmd, fmt.Sprintf("tmux session %q: failed to create tmux window %q", sessionName, windowName))
}

func SelectWindow(sessionName string, windowId int) error {
	if sessionName == "" {

	}
	cmd := exec.Command(
		"tmux",
		"select-window",
		"-t", fmt.Sprintf("\"%s:%d\"", sessionName, windowId))
	return run(cmd, fmt.Sprintf("tmux session %q: failed to select tmux window %d", sessionName, windowId))
}
