package config

import (
	"log"

	"github.com/Masterminds/semver"
)

type SupportedClient struct {
	ServerVersion *semver.Version
	Constraints   *semver.Constraints
}

// When bumping the server version, ensure that it is either backwards
// compatible or it is properly reflected here in which client versions it supports
func buildSupportedClientList() []*SupportedClient {
	return []*SupportedClient{
		{
			ServerVersion: semver.MustParse("1.0.0"),
			Constraints:   mustParseNewConstraint("^1.0.0"),
		},
	}
}

func GetSupportedClient(version *semver.Version) *SupportedClient {
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
