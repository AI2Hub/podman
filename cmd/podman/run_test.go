package main

import (
	"testing"

	ociv1 "github.com/opencontainers/image-spec/specs-go/v1"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/projectatomic/libpod/libpod"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

var (
	cmd         = []string{"podman", "test", "alpine"}
	CLI         *cli.Context
	testCommand = cli.Command{
		Name:   "test",
		Flags:  createFlags,
		Action: testCmd,
	}
)

// generates a mocked ImageData structure based on alpine
func generateAlpineImageData() *libpod.ImageData {
	config := &ociv1.ImageConfig{
		User:         "",
		ExposedPorts: nil,
		Env:          []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Entrypoint:   []string{},
		Cmd:          []string{"/bin/sh"},
		Volumes:      nil,
		WorkingDir:   "",
		Labels:       nil,
		StopSignal:   "",
	}

	data := &libpod.ImageData{
		ID:           "e21c333399e0aeedfd70e8827c9fba3f8e9b170ef8a48a29945eb7702bf6aa5f",
		RepoTags:     []string{"docker.io/library/alpine:latest"},
		RepoDigests:  []string{"docker.io/library/alpine@sha256:5cb04fce748f576d7b72a37850641de8bd725365519673c643ef2d14819b42c6"},
		Comment:      "Created:2017-12-01 18:48:48.949613376 +0000",
		Author:       "",
		Architecture: "amd64",
		Os:           "linux",
		Version:      "17.06.2-ce",
		Config:       config,
	}
	return data
}

// sets a global CLI
func testCmd(c *cli.Context) error {
	CLI = c
	return nil
}

// creates the mocked cli pointing to our create flags
// global flags like log-level are not implemented
func createCLI() cli.App {
	a := cli.App{
		Commands: []cli.Command{
			testCommand,
		},
	}
	return a
}

func getRuntimeSpec(c *cli.Context) *spec.Spec {
	runtime, _ := getRuntime(c)
	createConfig, _ := parseCreateOpts(c, runtime, "alpine", generateAlpineImageData())
	runtimeSpec, _ := createConfigToOCISpec(createConfig)
	return runtimeSpec
}

// TestPIDsLimit verifies the inputed pid-limit is correctly defined in the spec
func TestPIDsLimit(t *testing.T) {
	a := createCLI()
	args := []string{"--pids-limit", "22"}
	a.Run(append(cmd, args...))
	runtimeSpec := getRuntimeSpec(CLI)
	assert.Equal(t, runtimeSpec.Linux.Resources.Pids.Limit, int64(22))
}
