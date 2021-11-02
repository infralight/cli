package client

import (
	"encoding/json"
	"fmt"
)

const (
	FormatVersion = "0.1"
)

type InventorySearchResult struct {
	Aggregations    map[string]interface{} `json:"aggregations"`
	ResponseObjects []Resource             `json:"responseObjects"`
	TotalObjects    int                    `json:"totalObjects"`
}

type Resource struct {
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	ProviderID    string `json:"providerId"`
	AssetID       string `json:"assetId"`
	AssetType     string `json:"assetType"`
	Region        string `json:"region"`
	EnvironmentID string `json:"environmentId"`
	InventoryItem string `json:"inventoryItem"`
	IacType       string `json:"iacType"`
}

type ResourceDetails struct {
	ARN               string     `json:"arn,omitempty"`
	Name              string     `json:"name,omitempty"`
	Type              string     `json:"type,omitempty"`
	ModuleAddress     string     `json:"module_address,omitempty"`
	Address           string     `json:"address,omitempty"`
	ProviderName      string     `json:"provider_name,omitempty"`
	EffectedResources []Resource `json:"effected_resources,omitempty"`
}

type Implication struct {
	FormatVersion string            `json:"format_version,omitempty"`
	Resources     []ResourceDetails `json:"resources,omitempty"`
	TotalObjects  int               `json:"total_objects,omitempty"`
}

type InventoryClient struct {
	*baseClient
}

func (c *InventoryClient) NewSearchInput(query string) (map[string]interface{}, error) {
	inventorySearchJson := `{
  "queryPage": 0,
  "queryAdditionalFilters": [],
  "queryAdditionalFilterOut": [],
  "aggregations": {
    "is_managed_by_env": {
      "terms": {
        "field": "state.keyword",
        "size": 10000
      },
      "aggs": {
        "environment": {
          "terms": {
            "field": "environmentId.keyword",
            "size": 10000
          }
        }
      }
    },
    "ismanaged": {
      "terms": {
        "field": "state.keyword",
        "size": 10000
      }
    },
    "assetType": {
      "terms": {
        "field": "assetType.keyword",
        "size": 10000
      }
    },
    "environment": {
      "terms": {
        "field": "environmentId.keyword",
        "size": 10000
      }
    },
    "providerId": {
      "terms": {
        "field": "providerId.keyword",
        "size": 10000
      }
    },
    "region": {
      "terms": {
        "field": "region.keyword",
        "size": 10000
      }
    },
    "assetsVolume": {
      "terms": {
        "field": "resourceCreationDate",
        "size": 10000,
        "order": {
          "_key": "desc"
        },
        "min_doc_count": 1
      }
    }
  },
  "freeTextSearch": {
    "query_string": {
      "query": "\"%s\""
    }
  },
  "isDateSortedDescending": true,
  "isExcludedReturns": false
}`

	newQuery := fmt.Sprintf(inventorySearchJson, query)

	// Declared an empty map interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(newQuery), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *InventoryClient) SearchInventory(input map[string]interface{}) (result InventorySearchResult, err error) {
	err = c.httpc.NewRequest("POST", "/inventoryV2").
		JSONBody(input).
		Into(&result).
		Run()
	return result, err
}
