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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var appConfig = config{DirectoryAliasMap: make(map[string]*aliasConfig)}

type aliasConfig struct {
	Path string `mapstructure:"path"`
}

type config struct {
	force             bool
	DirectoryAliasMap map[string]*aliasConfig `mapstructure:"directory-alias"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "dTail",
	Long: `dTail 是一个类似于 tail -f 的工具，专为文件夹设计。它可以持续监控文件夹内的最新文件，并将新写入的内容输出到终端。
特别适用于按小时等更细粒度切割生成新文件的日志系统。
更多信息：https://github.com/bynow2code/dtail`,
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dTail.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&appConfig.force, "force", "f", false, "强制操作")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dTail")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}

	err := viper.Unmarshal(&appConfig)
	cobra.CheckErr(err)
}