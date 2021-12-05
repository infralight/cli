package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/infralight/cli/client"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/spf13/cobra"
)

var (
	planFile string
)

var implicationsCmd = &cobra.Command{
	Use:   "implications [cmd]",
	Short: "implications",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, args []string) error {
		var count = 0

		resources, err := getIdsFromPlan(planFile)
		if err != nil {
			return err
		}

		for index, resource := range resources {
			input, err := c.Inventory.NewSearchInput(resource.ID)
			if err != nil {
				return err
			}

			result, err := c.Inventory.SearchInventory(input)
			if err != nil {
				return fmt.Errorf("failed to search in the inventory, error: %w", err)
			}

			for _, responseObj := range result.ResponseObjects {
				if responseObj.AssetID == resource.ARN {
					continue
				}

				count += 1
				resources[index].EffectedResources = append(resources[index].EffectedResources, responseObj)
			}

		}

		return render(client.Implication{
			FormatVersion: client.FormatVersion,
			Resources:     resources,
			TotalObjects:  count,
		})
	},
}

func getIdsFromPlan(path string) ([]client.ResourceDetails, error) {
	resources := make([]client.ResourceDetails, 0)

	f, err := os.Open(path)
	if err != nil {
		return resources, err
	}
	defer f.Close() // nolint: errcheck

	var plan *tfjson.Plan
	err = json.NewDecoder(f).Decode(&plan)
	if err != nil {
		return resources, err
	}

	err = plan.Validate()
	if err != nil {
		return resources, err
	}

	for _, item := range plan.ResourceChanges {
		if item.ProviderName != "registry.terraform.io/hashicorp/aws" {
			continue
		}

		arn, id := getArnAndIdFromChange(item.Change.Before)
		if id == "" {
			continue
		}
		resources = append(resources, client.ResourceDetails{
			ARN:               arn,
			ID:                id,
			Name:              item.Name,
			Type:              item.Type,
			ModuleAddress:     item.ModuleAddress,
			Address:           item.Address,
			ProviderName:      item.ProviderName,
			EffectedResources: []client.Resource{},
		})
	}

	return resources, nil
}

func getArnAndIdFromChange(in interface{}) (string, string) {
	v := reflect.ValueOf(in)
	arn := ""
	id := ""

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			attribute := v.MapIndex(key)

			if key.Kind() == reflect.String {
				if key.String() == "arn" {
					arn = fmt.Sprintf("%v", attribute.Interface())
				}

				if key.String() == "id" {
					id = fmt.Sprintf("%v", attribute.Interface())
				}
			}
		}
	}

	return arn, id
}

func init() {
	implicationsCmd.Flags().StringVarP(&planFile, "plan-file", "f", "plan.json", "path for terraform plan json file")
	rootCmd.AddCommand(implicationsCmd)
}
