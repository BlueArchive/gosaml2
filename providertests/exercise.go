// +build go1.7

package providertests

import (
	"testing"

	saml2 "github.com/bluearchive/gosaml2"
	"github.com/stretchr/testify/require"
)

func ExerciseProviderTestScenarios(t *testing.T, scenarios []ProviderTestScenario) {
	for _, scenario := range scenarios {
		t.Run(scenario.ScenarioName, func(t *testing.T) {
			_, err := saml2.DecodeUnverifiedBaseResponse(scenario.Response)
			// DecodeUnverifiedBaseResponse is more permissive than RetrieveAssertionInfo.
			// If an error _is_ returned it should match, but it is OK for no error to be
			// returned even when one is expected during full validation.
			if err != nil {
				scenario.CheckError(t, err)
			}

			assertionInfo, err := scenario.ServiceProvider.RetrieveAssertionInfo(scenario.Response)
			if scenario.CheckError != nil {
				scenario.CheckError(t, err)
			} else {
				require.NoError(t, err)
			}

			if err == nil {
				if scenario.CheckWarningInfo != nil {
					scenario.CheckWarningInfo(t, assertionInfo.WarningInfo)
				} else {
					require.False(t, assertionInfo.WarningInfo.InvalidTime)
					require.False(t, assertionInfo.WarningInfo.NotInAudience)
				}
			}
		})
	}
}
