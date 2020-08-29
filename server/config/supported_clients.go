package config

import (
	"log"

	"github.com/Masterminds/semver"
)

type SupportedClient struct {
	ServerVersion *semver.Version
	Constraints   *semver.Constraints
}

func buildSupportedClientList() []*SupportedClient {
	return []*SupportedClient{
		{
			ServerVersion: semver.MustParse("1.0.0"),
			Constraints:   mustParseNewConstraint("^1.0.0"),
		},
	}
}

func GetSupportedClients(version *semver.Version) *SupportedClient {
	scList := buildSupportedClientList()

	for _, sc := range scList {
		if sc.ServerVersion.Equal(version) {
			return sc
		}
	}

	return nil
}

func mustParseNewConstraint(c string) *semver.Constraints {
	constraint, err := semver.NewConstraint(c)
	if err != nil {
		log.Fatalf("Failed making new constraint: %s", c)
	}

	return constraint
}
