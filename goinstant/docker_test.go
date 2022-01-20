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

	type args struct {
		startupCommands          []string
		heartbeatWantedBefore    bool
		heartbeatNotWantedBefore bool
		heartbeatWantedAfter     bool
		heartbeatNotWantedAfter  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1: Attempt to init OpenHIM Core",
			args: args{
				startupCommands:          []string{"docker", "core", "init"},
				heartbeatNotWantedBefore: true,
				heartbeatWantedAfter:     true,
			},
			wantErr: false,
		},
		{
			name: "Test 2: Attempt to bring OpenHIM Core down",
			args: args{
				startupCommands:         []string{"docker", "core", "down"},
				heartbeatWantedBefore:   true,
				heartbeatNotWantedAfter: true,
			},
			wantErr: false,
		},
		{
			name: "Test 3: Attempt to bring OpenHIM Core up.",
			args: args{
				startupCommands:          []string{"docker", "core", "up"},
				heartbeatNotWantedBefore: true,
				heartbeatWantedAfter:     true,
			},
			wantErr: false,
		},
		{
			name: "Test 4: Attempt to destroy OpenHIM Core.",
			args: args{
				startupCommands:         []string{"docker", "core", "destroy"},
				heartbeatWantedBefore:   true,
				heartbeatNotWantedAfter: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		os.Stdout = nil

		t.Run(tt.name, func(t *testing.T) {
			hbCheck := CheckOpenHIMheartbeat()
			if tt.args.heartbeatWantedBefore {
				if !hbCheck {
					t.Fatal("Expected heartbeat and not found")
				}
			}
			if tt.args.heartbeatNotWantedBefore {
				if hbCheck {
					t.Fatal("Heartbeat found when not expected")
				}
			}

			if err := RunDirectDockerCommand(tt.args.startupCommands); (err != nil) != tt.wantErr {
				t.Errorf("RunDirectDockerCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			hbCheck = CheckOpenHIMheartbeat()
			if tt.args.heartbeatWantedAfter {
				if !hbCheck {
					t.Fatal("Expected heartbeat and not found")
				}
			}
			if tt.args.heartbeatNotWantedAfter {
				if hbCheck {
					t.Fatal("Heartbeat found when not expected")
				}
			}

			t.Log(t.Name() + " passed!\n")
		})
	}
}

func TestSliceContains(t *testing.T) {
	testCases := []struct {
		slice    []string
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

func Test_extractCommands(t *testing.T) {
	type resultStruct struct {
		environmentVariables []string
		deployCommand        string
		otherFlags           []string
		deployEnvironment    string
		packages             []string
		customPackagePaths   []string
		instantVersion       string
	}

	testCases := []struct {
		startupCommands []string
		expectedResults resultStruct
		testInfo        string
	}{
		{
			startupCommands: []string{"init", "docker", "--instant-version=v2.0.1", "-c=../test", "-c=../test1", "-e=NODE_ENV=dev", "-onlyFlag", "core"},
			expectedResults: resultStruct{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "init",
				otherFlags:           []string{"-onlyFlag"},
				deployEnvironment:    "docker",
				packages:             []string{"core"},
				customPackagePaths:   []string{"../test", "../test1"},
				instantVersion:       "v2.0.1",
			},
			testInfo: "Extract commands test 1 - should return the expected commands",
		},
		{
			startupCommands: []string{"up", "kubernetes", "--instant-version=v2.0.2", "-c=../test", "-c=../test1", "-e=NODE_ENV=dev", "-onlyFlag", "core"},
			expectedResults: resultStruct{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "up",
				otherFlags:           []string{"-onlyFlag"},
				deployEnvironment:    "kubernetes",
				packages:             []string{"core"},
				customPackagePaths:   []string{"../test", "../test1"},
				instantVersion:       "v2.0.2",
			},
			testInfo: "Extract commands test 2 - should return the expected commands",
		},
		{
			startupCommands: []string{"down", "k8s", "--instant-version=v2.0.2", "-c=../test", "-c=../test1", "--env-file=../test.env", "-onlyFlag", "core", "hapi-fhir"},
			expectedResults: resultStruct{
				environmentVariables: []string{"--env-file", "../test.env"},
				deployCommand:        "down",
				otherFlags:           []string{"-onlyFlag"},
				deployEnvironment:    "k8s",
				packages:             []string{"core", "hapi-fhir"},
				customPackagePaths:   []string{"../test", "../test1"},
				instantVersion:       "v2.0.2",
			},
			testInfo: "Extract commands test 3 - should return the expected commands",
		},
		{
			startupCommands: []string{"destroy", "swarm", "--instant-version=v2.0.2", "--custom-package=../test", "-c=../test1", "-e=NODE_ENV=dev", "--onlyFlag", "core", "hapi-fhir"},
			expectedResults: resultStruct{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "destroy",
				otherFlags:           []string{"--onlyFlag"},
				deployEnvironment:    "swarm",
				packages:             []string{"core", "hapi-fhir"},
				customPackagePaths:   []string{"../test", "../test1"},
				instantVersion:       "v2.0.2",
			},
			testInfo: "Extract commands test 4 - should return the expected commands",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testInfo, func(t *testing.T) {
			environmentVariables, deployCommand, otherFlags, deployEnvironment, packages, customPackagePaths, instantVersion := extractCommands(tt.startupCommands)

			if !assert.Equal(t, environmentVariables, tt.expectedResults.environmentVariables) {
				t.Fatal("ExtractCommands should return the correct environment variables")
			}
			if !assert.Equal(t, deployCommand, tt.expectedResults.deployCommand) {
				t.Fatal("ExtractCommands should return the correct deploy command")
			}
			if !assert.Equal(t, otherFlags, tt.expectedResults.otherFlags) {
				t.Fatal("ExtractCommands should return the correct 'otherFlags'")
			}
			if !assert.Equal(t, deployEnvironment, tt.expectedResults.deployEnvironment) {
				t.Fatal("ExtractCommands should return the correct deployEnvironment")
			}
			if !assert.Equal(t, packages, tt.expectedResults.packages) {
				t.Fatal("ExtractCommands should return the correct packages")
			}
			if !assert.Equal(t, customPackagePaths, tt.expectedResults.customPackagePaths) {
				t.Fatal("ExtractCommands should return the correct custom package paths")
			}
			if !assert.Equal(t, instantVersion, tt.expectedResults.instantVersion) {
				t.Fatal("ExtractCommands should return the correct instant version")
			}
			t.Log(tt.testInfo + " passed!")
		})
	}
}
