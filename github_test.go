package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resource "github.com/telia-oss/github-pr-resource"
)

func TestNewGithubClient(t *testing.T) {
	tests := []struct {
		description string
		source      resource.Source
		expect      struct {
			owner      string
			repository string
		}
	}{
		{
			description: "owner & repo set properly",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			expect: struct {
				owner      string
				repository string
			}{
				owner:      "itsdalmo",
				repository: "test-repository",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			client, err := resource.NewGithubClient(&tc.source)
			require.NoError(t, err)
			assert.Equal(t, tc.expect.owner, client.Owner)
			assert.Equal(t, tc.expect.repository, client.Repository)
		})
	}
}
