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
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/projectriff/riff/riff-cli/global"
)

type config struct {
	file string
}

func Root() *cobra.Command {
	var cfg config


	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "riff",
		Short: "Commands for creating and managing function resources",
		Long: `riff is for functions

the riff tool is used to create and manage function resources for the riff FaaS platform https://projectriff.io/`,
		SilenceErrors: true,
		DisableAutoGenTag: true,
	}

	cobra.OnInitialize(cfg.initConfig, initGlobal)
	rootCmd.PersistentFlags().StringVar(&cfg.file, "config", "", "config file (default is $HOME/.riff.yaml)")
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func (c config) initConfig() {

	if c.file != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.file)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".riff" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".riff")
	}

	viper.SetEnvPrefix("RIFF") // use RIFF_ as prefix
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

func initGlobal() {

	if os.Getenv("RIFF_INVOKER_VERSION") != "" {
		global.INVOKER_VERSION = os.Getenv("RIFF_INVOKER_VERSION")
	} else {
		if viper.GetString("invokerVersion") != "" {
			global.INVOKER_VERSION = viper.GetString("invokerVersion")
		} else {
			if viper.GetString("invoker-version") != "" {
				global.INVOKER_VERSION = viper.GetString("invoker-version")
			}
		}
	}
}
