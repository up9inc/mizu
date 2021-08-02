package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/up9inc/mizu/cli/mizu"
	"github.com/up9inc/mizu/cli/uiUtils"
	"io/ioutil"
)

var outputFileName string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate example config file to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().Visit(mizu.InitFlag)

		template := mizu.GetTemplateConfig()
		if outputFileName != "" {
			data := []byte(template)
			_ = ioutil.WriteFile(outputFileName, data, 0644)
			mizu.Log.Infof(fmt.Sprintf("Template File written to %s", fmt.Sprintf(uiUtils.Purple, outputFileName)))
		} else {
			mizu.Log.Debugf("Writing template config.\n%v", template)
			fmt.Printf("%v", template)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&outputFileName, "file", "f", "", "Save content to local file")
}
