/*
Copyright © 2024 changqq <https://github.com/bynow2code/dtail>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"dTail/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type config struct {
	force           bool
	FolderShorthand map[string]shorthandConfig `mapstructure:"folder_shorthand" yaml:"folder_shorthand"`
}

type shorthandConfig struct {
	FolderPath string `mapstructure:"folder_path" yaml:"folder_path"`
}

var cfgFile string
var appCfg = config{FolderShorthand: make(map[string]shorthandConfig)}

var rootCmd = &cobra.Command{
	Use:   "dtail",
	Short: `Dtail is a command-line tool similar to the unix command tail -f, with the difference that dtail is designed for folder.`,
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dtail.yaml).")
	rootCmd.PersistentFlags().BoolVarP(&appCfg.force, "force", "f", false, "force")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dtail")
	}

	if err := viper.ReadInConfig(); err == nil {
		util.PrintInfo("Use config file:", viper.ConfigFileUsed(), ".")
	} else {
		cobra.CheckErr(err)
	}

	err := viper.Unmarshal(&appCfg)
	cobra.CheckErr(err)
}
