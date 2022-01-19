package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	Instead of running the test below using the buttons (if you have extensions installed), run the test using this command
	cmd:	go test -timeout 2m -v -run ^TestRunDirectDockerCommand$ github.com/openhie/instant/goinstant
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
	loadConfig()

	testCases := []testingStruct{
		{
			cmds:                     []string{"docker", "core", "init"},
			testInfo:                 "Test 1: Attempt to init OpenHIM Core",
			heartbeatNotWantedBefore: true,
			heartbeatWantedAfter:     true,
		},
		{
			cmds:                    []string{"docker", "core", "down"},
			testInfo:                "Test 2: Attempt to bring OpenHIM Core down",
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
			testInfo:                "Test 4: Attempt to destroy OpenHIM Core.",
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
			os.Stdout = nil

			t.Run(tt.name, func(t *testing.T) {
				hbCheck := CheckOpenHIMheartbeat()
				if test.heartbeatWantedBefore {
					if !hbCheck {
						t.Fatal("Expected heartbeat and not found")
					}
				}
				if test.heartbeatNotWantedBefore {
					if hbCheck {
						t.Fatal("Heartbeat found when not expected")
					}
				}

				if err := RunDirectDockerCommand(tt.args.startupCommands); (err != nil) != tt.wantErr {
					t.Errorf("RunDirectDockerCommand() error = %v, wantErr %v", err, tt.wantErr)
				}

				hbCheck = CheckOpenHIMheartbeat()
				if test.heartbeatWantedAfter {
					if !hbCheck {
						t.Fatal("Expected heartbeat and not found")
					}
				}
				if test.heartbeatNotWantedAfter {
					if hbCheck {
						t.Fatal("Heartbeat found when not expected")
					}
				}

				t.Log(t.Name() + " passed!\n")
			})
		}
	}
}

func TestSliceContains(t *testing.T) {
	testCases := []struct{
		slice	   []string
		element  string
		result   bool
		testInfo string
	}{
		{
			testInfo: "SliceContain test - should return true when slice contains element",
			slice:    []string{"Optimus Prime", "Iron Hyde"},
			element:  "Optimus Prime",
			result:   true,
		},
		{
			testInfo: "SliceContain test - should return false when slice does not contain element",
			slice:    []string{"Optimus Prime", "Iron Hyde"},
			element:  "Megatron",
			result:   false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testInfo, func(t *testing.T) {
			ans := sliceContains(tt.slice, tt.element)

			if ans != tt.result {
				t.Fatal("SliceContains should return" + fmt.Sprintf("%t", tt.result) + "but returned" + fmt.Sprintf("%t", ans))
			}
			t.Log(tt.testInfo + " passed!")
		})
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

func Test_getPackagePaths(t *testing.T) {
	type args struct {
		inputArr []string
		flags    []string
	}
	tests := []struct {
		name             string
		args             args
		wantPackagePaths []string
	}{
		{
			name: "Test 1",
			args: args{
				inputArr: []string{"-c=../docs", "-c=./docs"},
				flags:    []string{"-c=", "--custom-package="},
			},
			wantPackagePaths: []string{"../docs", "./docs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPackagePaths := getPackagePaths(tt.args.inputArr, tt.args.flags); !assert.Equal(t, tt.wantPackagePaths, gotPackagePaths) {
				t.Errorf("getPackagePaths() = %v, want %v", gotPackagePaths, tt.wantPackagePaths)
			}
		})
	}
}

func Test_getEnvironmentVariables(t *testing.T) {
	type args struct {
		inputArr []string
		flags    []string
	}
	tests := []struct {
		name                     string
		args                     args
		wantEnvironmentVariables []string
	}{
		{
			name: "Test case environment variables found",
			args: args{
				inputArr: []string{"-e=NODE_ENV=PROD", "-e=DOMAIN_NAME=instant.com"},
				flags:    []string{"-e=", "--env-file="},
			},
			wantEnvironmentVariables: []string{"-e", "NODE_ENV=PROD", "-e", "DOMAIN_NAME=instant.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEnvironmentVariables := getEnvironmentVariables(tt.args.inputArr, tt.args.flags); !assert.Equal(t, tt.wantEnvironmentVariables, gotEnvironmentVariables) {
				t.Errorf("getEnvironmentVariables() = %v, want %v", gotEnvironmentVariables, tt.wantEnvironmentVariables)
			}
		})
	}
}
