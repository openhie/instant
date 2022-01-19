package main

import (
	"log"
	"net/http"
	"testing"
)

/*
	Instead of running the test below using the buttons, run the test using this command instead
	cmd:	go test -timeout -v 2m -run ^TestRunDirectDockerCommand$ github.com/openhie/instant/goinstant
	so that the timeout can be set manually
*/

type testingStruct struct {
	cmds                     []string
	testInfo                 string
	heartbeatWantedBefore    bool
	heartbeatNotWantedBefore bool
	heartbeatWantedAfter     bool
	heartbeatNotWantedAfter  bool
}

func TestRunDirectDockerCommand(t *testing.T) {
	testCases := []testingStruct{
		{
			cmds:                     []string{"docker", "core", "init"},
			testInfo:                 "Test 1: Attempt to init OpenHIM Core",
			heartbeatNotWantedBefore: true,
			heartbeatWantedAfter:     true,
		},
		{
			cmds:                    []string{"docker", "core", "down"},
			testInfo:                "Test 2: Attempt to bring OpenHIM down",
			heartbeatWantedBefore:   true,
			heartbeatNotWantedAfter: true,
		},
		{
			cmds:                     []string{"docker", "core", "up"},
			testInfo:                 "Test 3: Attempt to bring OpenHIM Core up.",
			heartbeatNotWantedBefore: true,
			heartbeatWantedAfter:     true,
		},
		{
			cmds:                    []string{"docker", "core", "destroy"},
			testInfo:                "Test 4: Attempt to destroy OpenHIM.",
			heartbeatWantedBefore:   true,
			heartbeatNotWantedAfter: true,
		},
	}

	type args struct {
		startupCommands []string
	}
	for _, test := range testCases {
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{
				name: test.testInfo,
				args: args{
					startupCommands: test.cmds,
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			hbCheck := CheckOpenHIMheartbeat()
			if test.heartbeatWantedBefore {
				if !hbCheck {
					t.Fatal("Heartbeat not found")
				}
			}

			if test.heartbeatNotWantedBefore {
				if hbCheck {
					t.Fatal("Heartbeat found when not wanted")
				}
			}

			t.Run(tt.name, func(t *testing.T) {
				if err := RunDirectDockerCommand(tt.args.startupCommands); (err != nil) != tt.wantErr {
					t.Errorf("RunDirectDockerCommand() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			hbCheck = CheckOpenHIMheartbeat()
			if test.heartbeatWantedAfter {
				if !hbCheck {
					t.Fatal("Heartbeat not found")
				}
			}

			if test.heartbeatNotWantedAfter {
				if hbCheck {
					t.Fatal("Heartbeat found when not wanted")
				}
			}
		}
	}
}

func CheckOpenHIMheartbeat() bool {
	resp, err := http.Get("http://localhost:9000")
	if resp == nil || resp.StatusCode != 200 {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return true
}
