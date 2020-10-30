package tmux

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/davidzchen/tiles/tilespb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		wantCmd     string
		hasErr      bool
	}{
		{
			name:        "valid",
			sessionName: "default",
			wantCmd:     "tmux new-session -d -s \"default\"",
			hasErr:      false,
		},
		{
			name:        "invalid-no-session-name",
			sessionName: "",
			wantCmd:     "unused",
			hasErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevRun := run
			defer func() {
				run = prevRun
			}()

			var gotCmd string
			run = func(cmd *exec.Cmd) error {
				gotCmd = cmd.String()
				return nil
			}
			err := NewSession(test.sessionName)
			if err != nil {
				t.Fatalf("NewSession(%q) got error: %v", test.sessionName, err)
			}
			if gotCmd != test.wantCmd {
				t.Fatalf("NewSession(%q)\n  got: %q\n  want: %q", gotCmd, test.wantCmd)
			}
		})
	}
}

func TestAttachSession(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		wantCmd     string
		hasError    bool
	}{
		{
			name:        "valid",
			sessionName: "default",
			wantCmd:     "tmux -2 attach-session -t \"default\"",
			hasError:    false,
		},
		{
			name:        "invalid-no-session-name",
			sessionName: "",
			wantCmd:     "unused",
			hasError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevRun := run
			defer func() {
				run = prevRun
			}()

			var gotCmd string
			run = func(cmd *exec.Cmd) error {
				gotCmd = cmd.String()
				return nil
			}
			err := AttachSession(test.sessionName)
			if err != nil {
				t.Fatalf("NewSession(%q) got error: %v", test.sessionName, err)
			}
			if gotCmd != test.wantCmd {
				t.Fatalf("NewSession(%q)\n  got: %q\n  want: %q", gotCmd, test.wantCmd)
			}
		})
	}
}

func TestStateToProto(t *testing.T) {
	output1 := "default\n"
	output2 := "default:src\n"
	output3 := "default:logs:/var/log:\n"
	output4 := "default:home:/home/dzc\n" +
		"default:notes:/home/dzc/notes\n" +
		"default:src:/home/dzc/go/src\n"

	tests := []struct {
		name       string
		tmuxOutput string
		want       *tilespb.TilesConfig
		wantError  string
	}{
		{
			name:       "empty",
			tmuxOutput: "",
			want:       &tilespb.TilesConfig{},
			wantError:  "",
		},
		{
			name:       "invalid-session-only",
			tmuxOutput: output1,
			want:       nil,
			wantError:  "malformed tmux output line",
		},
		{
			name:       "invalid-session-window-only",
			tmuxOutput: output2,
			want:       nil,
			wantError:  "malformed tmux output line",
		},
		{
			name:       "invalid-extra-colon",
			tmuxOutput: output3,
			want:       nil,
			wantError:  "malformed tmux output line",
		},
		{
			name:       "valid",
			tmuxOutput: output4,
			want: &tilespb.TilesConfig{
				TmuxSession: []*tilespb.TmuxSession{
					&tilespb.TmuxSession{
						Name: "default",
						Window: []*tilespb.TmuxSession_Window{
							&tilespb.TmuxSession_Window{
								Name: "home",
								Dir:  "/home/dzc",
							},
							&tilespb.TmuxSession_Window{
								Name: "notes",
								Dir:  "/home/dzc/notes",
							},
							&tilespb.TmuxSession_Window{
								Name: "src",
								Dir:  "/home/dzc/src",
							},
						},
					},
				},
			},
			wantError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := stateToProto(test.tmuxOutput)
			if test.wantError != "" {
				if err == nil {
					t.Fatalf("expected error \"%s\" but got no error", test.wantError)
				}
				if !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("want error \"%s\" but got \"%v\"", test.wantError, err)
				}
			}
			if err != nil {
				t.Fatalf("stateToProto() got error: %v", err)
			}
			if diff := cmp.Diff(test.want, got, protocmp.Transform()); diff != "" {
				t.Fatalf("stateToProto() got diff:\n%s", diff)
			}
		})
	}
}

func TestHasSession(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		stderr      string
		cmdError    error
		hasSession  bool
	}{
		{
			// If the session does not exist, tmux has-session exits with a non-zero exit code
			// and prints the line "can't find session <session-name>" to stderr.
			name:        "session-does-not-exist",
			sessionName: "oss",
			stderr:      "can't find session default",
			cmdError:    fmt.Errorf("unused"),
			hasSession:  false,
		},
		{
			// If the session exists, tmux has-session exits with status EXIT_SUCCESS with no
			// output to stdout or stderr.
			name:        "session-exists",
			sessionName: "default",
			stderr:      "",
			cmdError:    nil,
			hasSession:  true,
		},
		{
			name:        "command-error",
			sessionName: "util",
			stderr:      "usage: has-session [-t target-session]",
			cmdError:    fmt.Errorf("error message unused"),
			hasSession:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevRunAndGetStderr := runAndGetStderr
			defer func() {
				runAndGetStderr = prevRunAndGetStderr
			}()

			// Mock out runAndGetStderr to return the stderr and the error returned by
			// cmd.Run from the stderr and cmdError fields of the test case respectively.
			runAndGetStderr = func(cmd *exec.Cmd) (string, error) {
				return test.stderr, test.cmdError
			}

			has, err := HasSession(test.sessionName)
			if err != nil && test.cmdError == nil {
				t.Fatalf("HasSession(%q) got error %v but expected no error", test.sessionName, err)
			} else if err == nil && test.cmdError != nil {
				t.Fatalf("HasSession(%q) expected error but got no error", test.sessionName)
			}

			if has != test.hasSession {
				t.Fatalf("HasSession() want: %t; got: %t", test.hasSession, has)
			}
		})
	}
}

func TestNewWindow(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		windowName  string
		directory   string
		windowId    int
		wantCmd     string
		hasError    bool
	}{
		{
			name:        "valid",
			sessionName: "default",
			windowName:  "tiles",
			directory:   "~/go/src/github.com/davidzchen/tiles",
			windowId:    1,
			wantCmd:     "",
			hasError:    false,
		},
		{
			name:        "invalid-no-session-name",
			sessionName: "",
			windowName:  "tiles",
			directory:   "~/go/src/github.com/davidzchen/tiles",
			windowId:    1,
			wantCmd:     "unused",
			hasError:    true,
		},
		{
			name:        "invalid-no-window-name",
			sessionName: "default",
			windowName:  "",
			directory:   "~/go/src/github.com/davidzchen/tiles",
			windowId:    1,
			wantCmd:     "unused",
			hasError:    true,
		},
		{
			name:        "invalid-no-directory",
			sessionName: "default",
			windowName:  "tiles",
			directory:   "",
			windowId:    1,
			wantCmd:     "unused",
			hasError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func TestSelectWindow(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		windowId    int
		wantCmd     string
		hasError    bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
