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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasCmd represents the add command
var aliasCmd = &cobra.Command{
	Use:   "alias [目录别名] [目录绝对路径]",
	Short: "添加目录快捷操作",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirAlias := args[0]
		dirPath := args[1]
		if !appConfig.force {
			_, ok := appConfig.DirectoryAliasMap[dirAlias]
			if ok {
				util.PrintFatalError(dirAlias, "已存在，使用全局标志 -f 强制覆盖")
			}
		}

		appConfig.DirectoryAliasMap[dirAlias] = &aliasConfig{Path: dirPath}
		viper.Set("directory-alias-map", appConfig.DirectoryAliasMap)

		err := viper.WriteConfig()
		cobra.CheckErr(err)

		fmt.Println("配置已更新")
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
