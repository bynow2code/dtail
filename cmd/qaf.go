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
	"github.com/bynow2code/dtail/util"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var qafCmd = &cobra.Command{
	Use:   "qaf",
	Short: "Quick Access to Folders.",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all shortcut configurations.",
	Run: func(cmd *cobra.Command, args []string) {
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"shortcut", "folder path"})
		for shortcut, shortcutEntry := range appCfg.Qaf {
			tw.AppendRow(table.Row{shortcut, shortcutEntry.FolderPath})
		}
		fmt.Println(tw.Render())
	},
}

var addCmd = &cobra.Command{
	Use:   "add [shortcut] [folder absolute path]",
	Short: "Add a shortcut to the Folder's Path.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		shortcut := args[0]
		folderPath := args[1]
		if !appCfg.force {
			_, ok := appCfg.Qaf[shortcut]
			if ok {
				util.PrintError("duplicate shortcut, use -f to overwrite.")
			}
		}

		appCfg.Qaf[shortcut] = qafEntry{FolderPath: folderPath}
		viper.Set("qaf", appCfg.Qaf)

		err := viper.WriteConfig()
		cobra.CheckErr(err)

		util.PrintInfo("configuration file has been written.")
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [shortcut1] [shortcut2] ...",
	Short: "Remove the shortcut configuration for the folder, if deleting multiple shortcuts at once, any non-existent shortcuts will be ignored.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			_, ok := appCfg.Qaf[args[0]]
			if !ok {
				util.PrintError("shortcut does not exist.")
			}
		}

		for _, v := range args {
			delete(appCfg.Qaf, v)
		}

		viper.Set("qaf", appCfg.Qaf)

		err := viper.WriteConfig()
		cobra.CheckErr(err)

		util.PrintInfo("configuration file has been written.")
	},
}

func init() {
	qafCmd.AddCommand(listCmd, addCmd, removeCmd)
	rootCmd.AddCommand(qafCmd)
}
