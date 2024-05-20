package cmd

import (
	"capten/pkg/clog"
	"fmt"

	"github.com/spf13/cobra"
)

func readAndValidCreateResourceFlags(cmd *cobra.Command) (resourceType string, err error) {
	resourceType, _ = cmd.Flags().GetString("type")
	if len(resourceType) == 0 {
		return "", fmt.Errorf("specify the resource type in the command line")
	}
	return
}

var resourceCreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "cluster resource create",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Resource Created, %s", resourceType)
	},
}

var resourceDeleteSubCmd = &cobra.Command{
	Use:   "delete",
	Short: "cluster resource delete",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Resource Deleted, %s", resourceType)
	},
}

var resourceListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "cluster resource list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Resources Listed, %s", resourceType)
	},
}

var resourceShowSubCmd = &cobra.Command{
	Use:   "show",
	Short: "cluster resource show",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Resources showed, %s", resourceType)
	},
}
