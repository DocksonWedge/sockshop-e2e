package sockshope2e_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

const CATALOGUE_URL = "http://35.222.229.152"

type Catalogue struct {
	Id          string   `json:"id"`
	Name        string   `json:"Name"`
	Description string   `json:"description"`
	ImageUrl    []string `json:"imageUrl"`
	Price       float32  `json:"price"`
	Count       int      `json:"count"`
	Tag         []string `json:"tag"`
}

func TestListcatalogues(t *testing.T) {

	respBody := []Catalogue{}
	apitest.New().
		EnableNetworking(cli).
		Get(fmt.Sprintf("%s/catalogue", CATALOGUE_URL)).
		Expect(t).
		Status(200).
		Assert(jsonpath.GreaterThan("$", 1)). // make sure multiple, we'd need to verify this is stable
		End().
		JSON(&respBody)

	for i, catalogue := range respBody {
		cJson, err := json.Marshal(catalogue)
		if err != nil {
			t.Errorf("Could not parse the Jsonfor catalog %d!", i)
		}
		if catalogue.Id == "" {
			t.Errorf("catalogue id was not set on catalogue %d! Response was %v", i, string(cJson)) // I like something like testify/assert, but will stick to native here
		}
		if catalogue.Name == "" {
			t.Errorf("catalogue name was not set on catalogue %d! Response was %v", i, string(cJson))
		}
		// I normally try to avoid asserting everything on an e2e test and keep it to sanity checks for the workflow, but we could change that if the contract is not checked elsewhere
	}
}
