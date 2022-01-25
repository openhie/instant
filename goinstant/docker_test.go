package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sliceContains(t *testing.T) {
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
			name: "Test 1 - '-c' flag",
			args: args{
				inputArr: []string{"-c=../docs", "-c=./docs"},
				flags:    []string{"-c=", "--custom-package="},
			},
			wantPackagePaths: []string{"../docs", "./docs"},
		},
		{
			name: "Test 2 - '--custom-package' flag",
			args: args{
				inputArr: []string{"--custom-package=../docs", "--custom-package=./docs"},
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
		{
			name: "Test case environment variables file found",
			args: args{
				inputArr: []string{"--env-file=../test.env"},
				flags:    []string{"-e=", "--env-file="},
			},
			wantEnvironmentVariables: []string{"--env-file", "../test.env"},
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

func Test_createZipFile(t *testing.T) {
	var reader io.Reader

	type args struct {
		file    string
		content io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		osCreate func(name string) (*os.File, error)
		ioCopy   func(dst io.Writer, src io.Reader) (written int64, err error)
	}{
		{
			name: "Test case create zip file no errors",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: false,
			osCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, nil
			},
		},
		{
			name: "Test case create zip file with errors from OsCreate",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			osCreate: func(name string) (*os.File, error) {
				return &os.File{}, errors.New("Test error")
			},
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, nil
			},
		},
		{
			name: "Test case create zip file with errors from IoCopy",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			osCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, errors.New("Test error")
			},
		},
		{
			name: "Test case create empty zip",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			osCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			ioCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OsCreate = tt.osCreate
			IoCopy = tt.ioCopy

			if err := createZipFile(tt.args.file, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("createZipFile() error = %v, wantErr %v", err, tt.wantErr)
				log.Println(tt.name, "failed!")
			} else {
				log.Println(tt.name, "passed!")
			}
		})
	}
}

func Test_runCommand(t *testing.T) {
	testCases := []struct {
		commandName    string
		suppressErrors []string
		commandSlice   []string
		pathToPackage  string
		errorString    error
		testInfo       string
	}{
		{
			commandName:    "docker",
			suppressErrors: nil,
			commandSlice:   []string{"ps"},
			pathToPackage:  "",
			errorString:    nil,
			testInfo:       "runCommand - run basic docker ps test",
		},
		{
			commandName:    "docker",
			suppressErrors: nil,
			commandSlice:   []string{"volume", "rm", "test-volume"},
			pathToPackage:  "",
			errorString:    fmt.Errorf("Error waiting for Cmd. Error: No such volume: test-volume\n: exit status 1"),
			testInfo:       "runCommand - remove nonexistant volume should return error",
		},
		{
			commandName:    "docker",
			suppressErrors: []string{"Error: No such volume: test-volume"},
			commandSlice:   []string{"volume", "rm", "test-volume"},
			pathToPackage:  "",
			errorString:    nil,
			testInfo:       "runCommand - remove nonexistant volume and suppress error",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testInfo, func(t *testing.T) {
			pathToPackage, err := runCommand(tt.commandName, tt.suppressErrors, tt.commandSlice...)
			if !assert.Equal(t, pathToPackage, tt.pathToPackage) {
				t.Fatal("RunCommand failed - path to package returned is incorrect")
			}
			if err != nil && tt.errorString != nil && !assert.Equal(t, err.Error(), tt.errorString.Error()) {
				t.Fatal("RunCommand failed - error returned incorrect")
			}

			if (err != nil && tt.errorString == nil) || (err == nil && tt.errorString != nil) {
				log.Fatal("RunCommand failed - error returned incorrect")
			}

			t.Log(tt.testInfo + " passed!")
		})
	}
}

func Test_unzipPackage(t *testing.T) {
	// OsOpenFile and IoCopy are intentionally not tested here since they
	// will read/write to and from the file system.
	type args struct {
		zipContent io.ReadCloser
	}
	tests := []struct {
		name              string
		args              args
		wantPathToPackage string
		wantErr           bool
		zipOpenReader     func(name string) (*zip.ReadCloser, error)
		osMkdirAll        func(path string, perm fs.FileMode) error
		filepathJoin      func(elem ...string) string
		osRemove          func(name string) error
	}{
		{
			name: "Test case no errors",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           false,
			zipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zReader := new(zip.ReadCloser)
				defer zReader.Close()

				file := new(zip.File)
				file.Name = "./"

				zReader.File = append(zReader.File, file)
				return zReader, nil
			},
			osMkdirAll: func(path string, perm fs.FileMode) error {
				return nil
			},
			filepathJoin: func(elem ...string) string {
				return ""
			},
			osRemove: func(name string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from ZipOpenReader",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			zipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zReader := new(zip.ReadCloser)
				defer zReader.Close()

				file := new(zip.File)
				file.Name = "./"

				zReader.File = append(zReader.File, file)
				return zReader, errors.New("ZipOpenReader error")
			},
			osMkdirAll: func(path string, perm fs.FileMode) error {
				return nil
			},
			filepathJoin: func(elem ...string) string {
				return ""
			},
			osRemove: func(name string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from OsMkdirAll",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			zipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zReader := new(zip.ReadCloser)
				defer zReader.Close()

				file := new(zip.File)
				file.Name = "./"

				zReader.File = append(zReader.File, file)
				return zReader, nil
			},
			osMkdirAll: func(path string, perm fs.FileMode) error {
				return errors.New("OsMkdirAll error")
			},
			filepathJoin: func(elem ...string) string {
				return ""
			},
			osRemove: func(name string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from OsRemove",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			zipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zReader := new(zip.ReadCloser)
				defer zReader.Close()

				file := new(zip.File)
				file.Name = "./"

				zReader.File = append(zReader.File, file)
				return zReader, nil
			},
			osMkdirAll: func(path string, perm fs.FileMode) error {
				return nil
			},
			filepathJoin: func(elem ...string) string {
				return ""
			},
			osRemove: func(name string) error {
				return errors.New("OsRemove error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultMockVariables()

			ZipOpenReader = tt.zipOpenReader
			OsMkdirAll = tt.osMkdirAll
			FilepathJoin = tt.filepathJoin
			OsRemove = tt.osRemove

			gotPathToPackage, err := unzipPackage(tt.args.zipContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("unzipPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPathToPackage != tt.wantPathToPackage {
				t.Errorf("unzipPackage() = %v, want %v", gotPathToPackage, tt.wantPathToPackage)
				log.Println(tt.name, "failed!")
			} else {
				log.Println(tt.name, "passed!")
			}
		})
	}
}

func defaultMockVariables() {
	OsCreate = func(name string) (*os.File, error) {
		return &os.File{}, nil
	}
	IoCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
		return 1, nil
	}
}
