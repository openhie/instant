package main

import (
	"testing"
)

func Test_executeCommand(t *testing.T) {
	tests := []struct {
		name           string
		CustomOptions  customOption
		deployCommands []string
	}{
		{
			name: "Test case assert startupPackages",
			CustomOptions: customOption{
				startupAction:   "down",
				startupPackages: []string{"core", "elastic-analytics"},
				instantVersion:  "latest",
				targetLauncher:  "docker",
			},
			deployCommands: []string{"down", "core", "elastic-analytics", "--instant-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert envVarFileLocation",
			CustomOptions: customOption{
				startupAction:      "init",
				envVarFileLocation: "./usr/bin",
				instantVersion:     "latest",
				targetLauncher:     "k8s",
			},
			deployCommands: []string{"init", "--env-file=./usr/bin", "--instant-version=latest", "-t=k8s"},
		},
		{
			name: "Test case assert envVars",
			CustomOptions: customOption{
				startupAction:  "up",
				envVars:        []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				instantVersion: "latest",
				targetLauncher: "docker",
			},
			deployCommands: []string{"up", "-e=NODE_ENV=DEV", "-e=DOMAIN_NAME=instant.com", "--instant-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert customPackageFileLocations",
			CustomOptions: customOption{
				startupAction:              "init",
				customPackageFileLocations: []string{"./local/cPack"},
				instantVersion:             "latest",
				targetLauncher:             "docker",
			},
			deployCommands: []string{"init", "-c=./local/cPack", "--instant-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert dev and only flags",
			CustomOptions: customOption{
				startupAction:  "destroy",
				onlyFlag:       true,
				instantVersion: "v1.02b",
				targetLauncher: "k8s",
				devMode:        true,
			},
			deployCommands: []string{"destroy", "--only", "--dev", "--instant-version=v1.02b", "-t=k8s"},
		},
		{
			name: "Test case assert all fields",
			CustomOptions: customOption{
				startupAction:              "init",
				startupPackages:            []string{"hmis", "mcsd"},
				envVarFileLocation:         "./home/bin",
				envVars:                    []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				customPackageFileLocations: []string{"./usr/local/cPack"},
				onlyFlag:                   true,
				instantVersion:             "v1.03a",
				targetLauncher:             "k8s",
				devMode:                    true,
			},
			deployCommands: []string{"init", "hmis", "mcsd", "--env-file=./home/bin", "-e=NODE_ENV=DEV",
				"-e=DOMAIN_NAME=instant.com", "-c=./usr/local/cPack", "--only", "--dev", "--instant-version=v1.03a", "-t=k8s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customOptions = tt.CustomOptions
			runDeployCommand = func(startupCommands []string) error {
				return nil
			}

			executeCommand()
			for i, dc := range DeployCommands {
				if dc != tt.deployCommands[i] {
					t.Errorf("DeployCommands variable error, got = %v, expected %v", dc, tt.deployCommands[i])
				}
			}
		})
	}
}
