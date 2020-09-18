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

func TestStateToProto(t *testing.T) {
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
		cmdOutput   string
		cmdError    error
		hasSession  bool
	}{
		{
			name:        "not-found",
			sessionName: "default",
			cmdOutput:   "can't find session default",
			cmdError:    fmt.Errorf("unused"),
			hasSession:  false,
		},
	}

	// XXX: Change this to get the command constructed and check that.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevRunAndGetStderr := runAndGetStderr
			defer func() {
				runAndGetStderr = prevRunAndGetStderr
			}()
			runAndGetStderr = func(cmd *exec.Cmd) (string, error) {
				return test.cmdOutput, test.cmdError
			}

			has, err := HasSession(test.sessionName)
			if err != nil {

			}
			if has != test.hasSession {
				t.Fatalf("HasSession() want: %t; got: %t", test.hasSession, has)
			}
		})
	}
}

func TestGetState(t *testing.T) {
	tests := []struct {
		name string
	}{}
}
