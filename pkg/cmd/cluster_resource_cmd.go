package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/clog"
	"capten/pkg/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func readAndValidResourceIdentfierFlags(cmd *cobra.Command) (resourceType, id string, err error) {
	resourceType, _ = cmd.Flags().GetString("resource-type")
	if len(resourceType) == 0 {
		return "", "", fmt.Errorf("specify the resource type in the command line")
	}

	id, _ = cmd.Flags().GetString("id")
	if len(id) == 0 {
		return "", "", fmt.Errorf("specify the resource id in the command line")
	}
	return
}

func readCloudTypeAttributesFlags(cmd *cobra.Command, cloudType string) (attributes map[string]string, err error) {
	attributes = map[string]string{}
	switch cloudType {
	case "aws":
		accessKey, _ := cmd.Flags().GetString("access-key")
		if len(accessKey) == 0 {
			return nil, fmt.Errorf("specify the access key in the command line")
		}

		secretKey, _ := cmd.Flags().GetString("secret-key")
		if len(secretKey) == 0 {
			return nil, fmt.Errorf("specify the secret key in the command line")
		}

		attributes["access-key"] = accessKey
		attributes["secret-key"] = secretKey

	case "azure":
		clientId, _ := cmd.Flags().GetString("client-id")
		if len(clientId) == 0 {
			return nil, fmt.Errorf("specify the client id in the command line")
		}

		clientSecret, _ := cmd.Flags().GetString("client-secret")
		if len(clientSecret) == 0 {
			return nil, fmt.Errorf("specify the client secret in the command line")
		}

		attributes["client-id"] = clientId
		attributes["client-secret"] = clientSecret
	default:
		return nil, fmt.Errorf("invalid cloud type: %s", cloudType)
	}
	return
}

func readAndValidResourceDataFlags(cmd *cobra.Command, resourceType string) (attributes map[string]string, err error) {
	labels, _ := cmd.Flags().GetString("labels")
	attributes = map[string]string{
		"labels": labels,
	}

	switch resourceType {
	case "git-project":
		gitProjectUrl, _ := cmd.Flags().GetString("git-project-url")
		if len(resourceType) == 0 {
			return nil, fmt.Errorf("specify the git project url in the command line")
		}

		accessToken, _ := cmd.Flags().GetString("access-token")
		if len(accessToken) == 0 {
			return nil, fmt.Errorf("specify the access token in the command line")
		}

		userId, _ := cmd.Flags().GetString("user-id")
		if len(userId) == 0 {
			return nil, fmt.Errorf("specify the user id in the command line")
		}

		attributes["git-project-url"] = gitProjectUrl
		attributes["access-token"] = accessToken
		attributes["user-id"] = userId
	case "cloud-provider":
		cloudType, _ := cmd.Flags().GetString("cloud-type")
		if len(cloudType) == 0 {
			return nil, fmt.Errorf("specify the cloud type in the command line")
		}
		attributes["cloud-type"] = cloudType

		cloudAttributes, err := readCloudTypeAttributesFlags(cmd, cloudType)
		if err != nil {
			return nil, err
		}

		for key, value := range cloudAttributes {
			attributes[key] = value
		}
		log.Printf(" Cloud attributes %v \n", cloudAttributes)
	case "container-registry":
		registryUrl, _ := cmd.Flags().GetString("registry-url")
		if len(registryUrl) == 0 {
			return nil, fmt.Errorf("specify the registry url in the command line")
		}
		registryType, _ := cmd.Flags().GetString("registry-type")
		if len(registryType) == 0 {
			return nil, fmt.Errorf("specify the registry type in the command line")
		}
		registryUserName, _ := cmd.Flags().GetString("registry-username")
		if len(registryUserName) == 0 {
			return nil, fmt.Errorf("specify the registry username in the command line")
		}
		registryPassword, _ := cmd.Flags().GetString("registry-password")
		if len(registryPassword) == 0 {
			return nil, fmt.Errorf("specify the registry password in the command line")
		}

		attributes["registry-url"] = registryUrl
		attributes["registry-type"] = registryType
		attributes["registry-username"] = registryUserName
		attributes["registry-password"] = registryPassword
	default:
		return nil, fmt.Errorf("invalid resource type: %s", resourceType)
	}
	return
}

func readAndValidCreateResourceFlags(cmd *cobra.Command) (resourceType string, attributes map[string]string, err error) {
	resourceType, _ = cmd.Flags().GetString("resource-type")
	if len(resourceType) == 0 {
		return "", nil, fmt.Errorf("specify the resource type in the command line")
	}

	attributes, err = readAndValidResourceDataFlags(cmd, resourceType)
	if err != nil {
		return "", nil, err
	}
	return
}

func readAndValidUpdateResourceFlags(cmd *cobra.Command) (resourceType, id string, attributes map[string]string, err error) {
	resourceType, id, err = readAndValidResourceIdentfierFlags(cmd)
	if err != nil {
		return "", "", nil, err
	}

	attributes, err = readAndValidResourceDataFlags(cmd, resourceType)
	if err != nil {
		return "", "", nil, err
	}
	return
}

var resourceCreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "cluster resource create",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, attributes, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.AddClusterResource(captenConfig, resourceType, attributes)
		if err != nil {
			clog.Logger.Errorf("failed to create cluster resource, %v", err)
		}
	},
}

var resourceUpdateSubCmd = &cobra.Command{
	Use:   "update",
	Short: "cluster resource update",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, id, attributes, err := readAndValidUpdateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.UpdateClusterResource(captenConfig, resourceType, id, attributes)
		if err != nil {
			clog.Logger.Errorf("failed to update cluster resource, %v", err)
		}
	},
}

var resourceDeleteSubCmd = &cobra.Command{
	Use:   "delete",
	Short: "cluster resource delete",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, id, err := readAndValidResourceIdentfierFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}
		err = agent.RemoveClusterResource(captenConfig, resourceType, id)
		if err != nil {
			clog.Logger.Errorf("failed to delete cluster resource, %v", err)
		}
	},
}

var resourceListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "cluster resource list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		if len(resourceType) == 0 {
			clog.Logger.Error(fmt.Errorf("specify the resource type in the command line"))
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}
		err = agent.ListClusterResources(captenConfig, resourceType)
		if err != nil {
			clog.Logger.Errorf("failed to list cluster resources, %v", err)
		}
	},
}
