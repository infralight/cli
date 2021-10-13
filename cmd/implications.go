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
	RunE: func(_ *cobra.Command, _ []string) error {
		var count = 0

		resources, err := getIdsFromPlan(planFile)
		if err != nil {
			return err
		}

		for index, _ := range resources {
			input, err := c.Inventory.NewSearchInput(resources[index].ARN)
			if err != nil {
				return err
			}

			result, err := c.Inventory.SearchInventory(input)
			if err != nil {
				return fmt.Errorf("failed to search in the inventory, error: %w", err)
			}

			count += result.TotalObjects
			resources[index].EffectedResources = result.ResponseObjects
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
	defer f.Close()

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

		arn := getArnFromChange(item.Change.Before)
		if arn == "" {
			continue
		}
		resources = append(resources, client.ResourceDetails{
			ARN:               arn,
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

func getArnFromChange(in interface{}) string {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			attribute := v.MapIndex(key)

			if key.Kind() == reflect.String {
				if key.String() == "arn" {
					return fmt.Sprintf("%v", attribute.Interface())
				}
			}
		}
	}

	return ""
}

func init() {
	implicationsCmd.Flags().StringVarP(&planFile, "plan-file", "f", "plan.json", "path for terraform plan json file")
	rootCmd.AddCommand(implicationsCmd)
}
