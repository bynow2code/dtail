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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shorthandCmd = &cobra.Command{
	Use:   "shorthand",
	Short: "Manage folder shorthand for dtail.",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing shorthands.",
	Run: func(cmd *cobra.Command, args []string) {
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"shorthand name", "folder path"})
		for shorthandName, shorthandConfig := range appCfg.FolderShorthand {
			tw.AppendRow(table.Row{shorthandName, shorthandConfig.FolderPath})
		}
		fmt.Println(tw.Render())
	},
}

var addCmd = &cobra.Command{
	Use:   "add [shorthand name] [folder absolute path]",
	Short: "Add shorthand name for a folder path.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		shorthandName := args[0]
		folderPath := args[1]
		if !appCfg.force {
			_, ok := appCfg.FolderShorthand[shorthandName]
			if ok {
				util.PrintError("duplicate shorthand names: ", shorthandName, ", use -f to override.")
			}
		}

		appCfg.FolderShorthand[shorthandName] = shorthandConfig{FolderPath: folderPath}
		viper.Set("folder_shorthand", appCfg.FolderShorthand)

		err := viper.WriteConfig()
		cobra.CheckErr(err)

		util.PrintInfo("configuration file has been written.")
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [shorthand name]",
	Short: "Remove one or more shorthand names, if the shorthand name does not exist, the error will be skipped.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			delete(appCfg.FolderShorthand, v)
		}

		viper.Set("folder_shorthand", appCfg.FolderShorthand)

		err := viper.WriteConfig()
		cobra.CheckErr(err)

		util.PrintInfo("configuration file has been written.")
	},
}

func init() {
	shorthandCmd.AddCommand(listCmd, addCmd, removeCmd)
	rootCmd.AddCommand(shorthandCmd)
}
