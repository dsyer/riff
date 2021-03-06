/*
 * Copyright 2018 the original author or authors.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/projectriff/riff/riff-cli/pkg/osutils"
	"github.com/spf13/cobra"
	"github.com/projectriff/riff/riff-cli/pkg/options"
	"github.com/projectriff/riff/riff-cli/pkg/docker"
	"github.com/projectriff/riff/riff-cli/pkg/kubectl"
)

func TestCreateCommandPathFromArg(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "--dry-run", "-v", "0.0.1-snapshot", "../test_data/command/echo"})

	err := rootCmd.Execute()
	as.NoError(err)

	as.NotEmpty(initOptions.FilePath)
	as.NotEmpty(initOptions.UserAccount)
	as.Equal("../test_data/command/echo", initOptions.FilePath)
}


func TestCreateCommandFromCWD(t *testing.T) {
	rootCmd, _, _, _:= setupCreateTest()
	currentdir := osutils.GetCWD()

	path := osutils.Path("../test_data/command/echo")
	os.Chdir(path)
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "--dry-run"})
	err := rootCmd.Execute()
	as.NoError(err)

	os.Chdir(currentdir)
}

func TestCreateCommandExplicitPath(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "--dry-run", "-f", osutils.Path("../test_data/command/echo"), "-v", "0.0.1-snapshot"})

	err := rootCmd.Execute()
	as.NoError(err)

	as.NotEmpty(initOptions.FilePath)
	as.NotEmpty(initOptions.UserAccount)
	as.Equal("../test_data/command/echo", initOptions.FilePath)
}

func TestCreateCommandWithUser(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "--dry-run", "-f", "../test_data/command/echo", "-u", "me"})

	err := rootCmd.Execute()
	as.NoError(err)

	as.NotEmpty(initOptions.FilePath)
	as.Equal("me", initOptions.UserAccount)
	as.Equal("../test_data/command/echo", initOptions.FilePath)
}

func TestCreateCommandExplicitPathAndLang(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "command", "--dry-run", "-f", osutils.Path("../test_data/command/echo"), "-v", "0.0.1-snapshot"})

	err := rootCmd.Execute()
	as.NoError(err)

	as.NotEmpty(initOptions.FilePath)
	as.Equal("../test_data/command/echo", initOptions.FilePath)
}

func TestCreateLanguageDoesNotMatchArtifact(t *testing.T) {
	rootCmd, _, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "command", "--dry-run", "-f", osutils.Path("../test_data/python/demo"), "-a", "demo.py"})

	err := rootCmd.Execute()
	as.Error(err)
	as.Equal("language command conflicts with artifact file extension demo.py", err.Error())
}

func TestCreatePythonCommand(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "python", "--dry-run", "-f", osutils.Path("../test_data/python/demo"), "-v", "0.0.1-snapshot", "--handler", "process"})

	err := rootCmd.Execute()
	as.NoError(err)
	as.NotEmpty(initOptions.UserAccount)
	as.Equal("process", initOptions.Handler)
}

func TestCreatePythonCommandWithDefaultHandler(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "python", "--dry-run", "-f", osutils.Path("../test_data/python/demo"), "-v", "0.0.1-snapshot"})

	err := rootCmd.Execute()
	as.NoError(err)
	as.Equal("demo", initOptions.Handler)
}

func TestCreateJavaWithVersion(t *testing.T) {
	rootCmd, initOptions, _, _:= setupCreateTest()
	currentdir := osutils.GetCWD()
	path := osutils.Path("../test_data/java")
	os.Chdir(path)
	as := assert.New(t)
	rootCmd.SetArgs([]string{"create", "java", "--dry-run", "-a", "target/upper-1.0.0.jar", "--handler", "function.Upper"})
	err := rootCmd.Execute()
	as.NoError(err)
	as.NotEmpty(initOptions.UserAccount)
	os.Chdir(currentdir)
}

func setupCreateTest() (*cobra.Command, *options.InitOptions, *BuildOptions, *ApplyOptions) {
	rootCmd, initOptions, initCommands := setupInitTest()

	buildCmd, buildOptions := Build(docker.RealDocker(), docker.DryRunDocker())

	applyCmd, applyOptions := Apply(kubectl.RealKubeCtl(), kubectl.DryRunKubeCtl())

	createCmd := Create(initCommands["init"], buildCmd, applyCmd)

	createNodeCmd := CreateNode(initCommands["node"], buildCmd, applyCmd)
	createJavaCmd := CreateJava(initCommands["java"], buildCmd, applyCmd)
	createPythonCmd := CreatePython(initCommands["python"], buildCmd, applyCmd)
	createShellCmd := CreateCommand(initCommands["command"], buildCmd, applyCmd)
	createGoCmd := CreateGo(initCommands["go"], buildCmd, applyCmd)

	createCmd.AddCommand(
		createNodeCmd,
		createJavaCmd,
		createPythonCmd,
		createShellCmd,
		createGoCmd,
	)

	rootCmd.AddCommand(
		createCmd,
	)


	return rootCmd, initOptions, buildOptions, applyOptions
}
