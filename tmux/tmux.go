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

func HasSession(sessionName string) (bool, error) {
	cmd := exec.Command("tmux", "has-session", "-t", quote(sessionName))
	output, err := runAndGetStderr(cmd)
	if err != nil {
		if strings.HasPrefix(output, "can't find session") {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed while checking for tmux session %s", sessionName)
	}
	return true, nil
}

func NewSession(sessionName string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", quote(sessionName))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to start tmux session %s", sessionName)
	}
	return nil
}

func AttachSession(sessionName string) error {
	cmd := exec.Command("tmux", "-2", "attach-session", "-t", quote(sessionName))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to attach to tmux session %s", sessionName)
	}
	return nil
}

func ListSessions() error {
	cmd := exec.Command("tmux", "list-sessions")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to list tmux sessions")
	}
	return nil
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
	cmd := exec.Command(
		"tmux",
		"new-window",
		"-c", quote(directory),
		"-t", fmt.Sprintf("\"%s:%d\"", sessionName, windowId),
		"-n", quote(windowName))
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to create tmux window")
	}
	return nil
}

func SelectWindow(sessionName string, windowId int32) error {
	cmd := exec.Command(
		"tmux",
		"select-window",
		"-t", fmt.Sprintf("\"%s:%d\"", sessionName, windowId))
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to select tmux window")
	}
	return nil
}
