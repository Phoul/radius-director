package validation

import (
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/model"
)

func TestValidate(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"default": {},
			},
			AuthenticationProfiles: map[string]model.AuthenticationProfile{
				"default": {},
			},
			AccountingProfiles: map[string]model.AccountingProfile{
				"default": {},
			},
			MonitoringProfiles: map[string]model.MonitoringProfile{
				"default": {},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				RADIUSServers: map[string]model.RADIUSServer{
					"radius-1": {},
				},
				NASAssignments: map[string]model.NASAssignment{
					"core": {},
				},
			},
		},
	}

	if err := Validate(configuration); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestSortedKeys(t *testing.T) {
	values := map[string]int{
		"zulu":   1,
		"alpha":  2,
		"middle": 3,
	}

	if got, want := sortedKeys(values), []string{"alpha", "middle", "zulu"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("sortedKeys() = %v, want %v", got, want)
	}
}
