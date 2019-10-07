package resource_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		description string
		json        []byte
		request     resource.CheckRequest
	}{
		{
			description: "simple",
			json:        []byte(`{"source":{"access_token":"XXXXX","paths":["test/*","ci/*","README*"],"preview_schema":true,"repository":"digitalocean/github-pr-resource","v3_endpoint":"https://github.com/api/v3/", "v4_endpoint":"https://github.com/api/graphql"},"version":null}`),
			request: resource.CheckRequest{
				Source: resource.Source{
					AccessToken:   "XXXXX",
					Repository:    "digitalocean/github-pr-resource",
					PreviewSchema: true,
					Paths:         []string{"test/*", "ci/*", "README*"},
					V3Endpoint:    "https://github.com/api/v3/",
					V4Endpoint:    "https://github.com/api/graphql",
				},
				Version: resource.Version{},
			},
		},
		{
			description: "simple with version",
			json:        []byte(`{"source":{"access_token":"XXXXX","paths":["test/*","ci/*","README*"],"preview_schema":true,"repository":"digitalocean/github-pr-resource","v3_endpoint":"https://github.com/api/v3/", "v4_endpoint":"https://github.com/api/graphql"},"version":{"pr":"1","commit":"a4afe32","updated":"2019-08-20T00:14:16Z"}}`),
			request: resource.CheckRequest{
				Source: resource.Source{
					AccessToken:   "XXXXX",
					Repository:    "digitalocean/github-pr-resource",
					PreviewSchema: true,
					Paths:         []string{"test/*", "ci/*", "README*"},
					V3Endpoint:    "https://github.com/api/v3/",
					V4Endpoint:    "https://github.com/api/graphql",
				},
				Version: resource.Version{
					PR:          1,
					Commit:      "a4afe32",
					UpdatedDate: time.Date(2019, time.August, 20, 0, 14, 16, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			var output resource.CheckRequest

			err := json.Unmarshal(tc.json, &output)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.request.Source, output.Source)
				assert.Equal(t, tc.request.Version, output.Version)

				err = output.Source.Validate()
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		description string
		json        []byte
		versions    []resource.Version
	}{
		{
			description: "empty",
			json:        []byte(`[]`),
			versions:    []resource.Version{},
		},
		{
			description: "one version",
			json:        []byte(`[{"pr":"1","commit":"a4afe32","updated":"2019-08-20T00:14:16Z"}]`),
			versions: []resource.Version{
				{
					PR:          1,
					Commit:      "a4afe32",
					UpdatedDate: time.Date(2019, time.August, 20, 0, 14, 16, 0, time.UTC),
				},
			},
		},
		{
			description: "many versions",
			json:        []byte(`[{"pr":"1","commit":"a4afe32","updated":"2019-08-20T00:14:16Z"},{"pr":"2","commit":"a4afe33","updated":"2020-08-20T00:14:16Z"}]`),
			versions: []resource.Version{
				{
					PR:          1,
					Commit:      "a4afe32",
					UpdatedDate: time.Date(2019, time.August, 20, 0, 14, 16, 0, time.UTC),
				},
				{
					PR:          2,
					Commit:      "a4afe33",
					UpdatedDate: time.Date(2020, time.August, 20, 0, 14, 16, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			output, err := json.Marshal(tc.versions)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.json, output)
			}
		})
	}
}
