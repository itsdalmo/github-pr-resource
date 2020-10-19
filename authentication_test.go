package resource_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	resource "github.com/telia-oss/github-pr-resource"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {

	validInstallationResponse := `
{
  "id": 9912873,
  "account": {
    "login": "github",
    "id": 1,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
    "avatar_url": "https://github.com/images/error/hubot_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/orgs/github",
    "html_url": "https://github.com/github",
    "followers_url": "https://api.github.com/users/github/followers",
    "following_url": "https://api.github.com/users/github/following{/other_user}",
    "gists_url": "https://api.github.com/users/github/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/github/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/github/subscriptions",
    "organizations_url": "https://api.github.com/users/github/orgs",
    "repos_url": "https://api.github.com/orgs/github/repos",
    "events_url": "https://api.github.com/orgs/github/events",
    "received_events_url": "https://api.github.com/users/github/received_events",
    "type": "Organization",
    "site_admin": false
  },
  "repository_selection": "all",
  "access_tokens_url": "https://api.github.com/installations/1/access_tokens",
  "repositories_url": "https://api.github.com/installation/repositories",
  "html_url": "https://github.com/organizations/github/settings/installations/1",
  "app_id": 1,
  "target_id": 1,
  "target_type": "Organization",
  "permissions": {
    "checks": "write",
    "metadata": "read",
    "contents": "read"
  },
  "events": [
    "push",
    "pull_request"
  ],
  "created_at": "2018-02-09T20:51:14Z",
  "updated_at": "2018-02-09T20:51:14Z",
  "single_file_name": null
}
`

	validTokenResponse := `
{
  "token": "v1.b71be873ad96e64a84025ae7bee7694a99cb4ba9",
  "expires_at": "2020-06-21T00:03:29Z",
  "permissions": {
    "checks": "write",
    "contents": "read",
    "metadata": "read",
    "pull_requests": "write"
  },
  "repository_selection": "selected"
}
`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/repos/itsdalmo/test-repository/installation" {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/vnd.github.machine-man-preview+json", r.Header.Get("Accept"))
			assert.Equal(t, "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI2OTE5ODIsImlhdCI6MTU5MjY5MTQ0MiwiaXNzIjoiNjk1OTIifQ.H_a6i7TpaGOsaoliH_i7AT5UMwM9LqO21lEFYiZ96_H15cEF6D_kyZrcHyinP2fSC8rX_OQ-DvPsehTNTtfOhgM-nsdgg-gTzdCSlASgc00sGhw_pjBFwHpD6V1NQojc82L8SAR9Bg75g0xVlQ_dAR_Lbtmk252X_AlabRAfdSchK_GVdc3kSEbp28lc87EF7J_lFdRCNuVi1xcLFOXPmMeu4epSBf1ZMtuts28C7iqaI4QJ9keaGFug1wpL-WLDcFbvmB2nJhBYN9tArGM0ZHZ5i4EhFyFjGpBwTyo5P7WY7P3zYtz36gwgntYRtPPivcFQ-wUWuvMpL6vKd-Pp8w", r.Header.Get("Authorization"))
			_, _ = fmt.Fprintln(w, validInstallationResponse)
		}

		if r.RequestURI == "/app/installations/9912873/access_tokens" {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/vnd.github.machine-man-preview+json", r.Header.Get("Accept"))
			assert.Equal(t, "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI2OTE5ODIsImlhdCI6MTU5MjY5MTQ0MiwiaXNzIjoiNjk1OTIifQ.H_a6i7TpaGOsaoliH_i7AT5UMwM9LqO21lEFYiZ96_H15cEF6D_kyZrcHyinP2fSC8rX_OQ-DvPsehTNTtfOhgM-nsdgg-gTzdCSlASgc00sGhw_pjBFwHpD6V1NQojc82L8SAR9Bg75g0xVlQ_dAR_Lbtmk252X_AlabRAfdSchK_GVdc3kSEbp28lc87EF7J_lFdRCNuVi1xcLFOXPmMeu4epSBf1ZMtuts28C7iqaI4QJ9keaGFug1wpL-WLDcFbvmB2nJhBYN9tArGM0ZHZ5i4EhFyFjGpBwTyo5P7WY7P3zYtz36gwgntYRtPPivcFQ-wUWuvMpL6vKd-Pp8w", r.Header.Get("Authorization"))
			_, _ = fmt.Fprintln(w, validTokenResponse)
		}
	}))
	defer ts.Close()

	tests := []struct {
		description         string
		source              resource.Source
		expectedAccessToken string
	}{
		{
			description: "return given access token",
			source: resource.Source{
				Repository:  "itsdalmo/test-repository",
				AccessToken: "oauthtoken",
			},
			expectedAccessToken: "oauthtoken",
		},
		{
			description: "create access token",
			source: resource.Source{
				Repository: "itsdalmo/test-repository",
				AppId:      "69592",
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAv6oJaa+0BCOT0dITtJK6oPfrE2R0Iynofj/a/vq4rfAw2MJp\nXA5l6WbwmQhevSm8sIYNb2T32qLyRsVXwBo1J8ovIoTXiPdsV41D19tytjctnf+u\n0LncTo2JJR3ik7/Ynu5Id4zjnsn1pYqwLX4EFxgCuEKwdZutPyiY2J7wATfsTtOJ\nsLF8idQijG23i5Obs6AWCZcHOhvdgfYAUOxLv2WRkCG5O1aXYa6nqVn0AgRgKjaJ\nnAENEoG7O9OEWIcUiG30riQouxMHfj0bATCvYoj7a2tvn4CUqk+SBODyh3Fvi0oC\na/NCBYFeK/5uUNyXQHqWy7xEVRef5lC0XbrUVwIDAQABAoIBAFDdqhkIQ/iXFjAp\n5ZyDZ/CwiWNmN8X6UZiq0nhQSolA1SsvY4qunHsMrqiyql4/dNg5xwNf419A7t3D\nN5HavOCr4pU63UFxuyl5dc1mTpDo2PtXvGdec8BE4T9iy40xHXF48eRW8la1uUn+\nKPUYvRsNS2B46sDETSVfuJV1AahRN6aD2WnzQ7wB+S/mqsPXqy+S2zobnU70Wzmt\nhYQzsuIY+BkfOCS565guYvJt66wRGi/NsnybC6z3iRZxZygtKorMKXPjmtdoyCZ3\njkOHhXV5XH8Ldut+1mycg+6c+dTZ9RTrAzo4ouofptm9+ZNlv1aK7HrDDKYQ/x/z\n70hhIfkCgYEA7aNt7db3S6sTEA+yC5DLxxkFIK7K8qxdP6vE087fdHO+ik1LvyI+\ni5Wqj/fT2d/lYR+cbH/zy2JBy5RWfdVJ9HBW/eZ9RqQx8ry7lE5MqbwU6y+d2bCF\nAjffx3yJK1aljuVLdGeu8abYsusUQUhbslZJNg6Ar7z+FUHaChJI0KsCgYEAznk4\nx+5PtvIWj6aXTJiqtofnOtHXXxrq5alzBErvDBHHzP9/ZCucHxJZ4zqHWduj2+9b\nJUgk2BH3+wFMVKpcdTld2iTFWGaFsArTJkE4SBQGx3zcdsoESOnu/DrOTFNu+LiV\nhLzb9CHKM91fIet8MYYxKiyH09+Mi7xBtw8WQwUCgYAH/tCrCOmHHTll9/E4nGWO\nzFO01syzP4NfqgrUSYiRJXfKtXEP/Dn4fk+fymnRUcwo6WRc7i0osaSfEd2bHDsB\nw2nZ3xBl+Q5JKXpyMfQ4XcCibRa1hU/kVDbuQk1nLOIjHandP8POE5wE4Q3saF/V\nbzvFWtWPlB9EXdPVNOpIQwKBgHkJEsIQ72XdUGBxVewu6pQJ4wDWFhzIWL68sJHp\no2w92BRSCkmcTu7gARV1L/b7DHlXPOUD/6UyE15vCmHvZDfLozrHp3AE2YWzMsgQ\nH4ARTVAP3+U603wytkfh6SFRH5JqEiw30fCxBimVMbleo/UcJyID7LPFLkyT1SoM\njA5JAoGACEhaHTXWFYXv+eTJUXFcwhDZ5sRvQYRCTLGSv746lr+SpZ02bEcRvaTH\nv0K1Hph6OzhcCO27VdimngnzsoXoa2OXicWpdjX3hqhXMvwspe9a1y9u7LPZwqJH\n9Kc7VaZP1iS4EIxEc+qx38HqdeiUBcqTMRam6k3mSwactBKCKDI=\n-----END RSA PRIVATE KEY-----\n",
				V3Endpoint: ts.URL,
			},
			expectedAccessToken: "v1.b71be873ad96e64a84025ae7bee7694a99cb4ba9",
		},
		{
			description: "handle trailing slash",
			source: resource.Source{
				Repository: "itsdalmo/test-repository",
				AppId:      "69592",
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAv6oJaa+0BCOT0dITtJK6oPfrE2R0Iynofj/a/vq4rfAw2MJp\nXA5l6WbwmQhevSm8sIYNb2T32qLyRsVXwBo1J8ovIoTXiPdsV41D19tytjctnf+u\n0LncTo2JJR3ik7/Ynu5Id4zjnsn1pYqwLX4EFxgCuEKwdZutPyiY2J7wATfsTtOJ\nsLF8idQijG23i5Obs6AWCZcHOhvdgfYAUOxLv2WRkCG5O1aXYa6nqVn0AgRgKjaJ\nnAENEoG7O9OEWIcUiG30riQouxMHfj0bATCvYoj7a2tvn4CUqk+SBODyh3Fvi0oC\na/NCBYFeK/5uUNyXQHqWy7xEVRef5lC0XbrUVwIDAQABAoIBAFDdqhkIQ/iXFjAp\n5ZyDZ/CwiWNmN8X6UZiq0nhQSolA1SsvY4qunHsMrqiyql4/dNg5xwNf419A7t3D\nN5HavOCr4pU63UFxuyl5dc1mTpDo2PtXvGdec8BE4T9iy40xHXF48eRW8la1uUn+\nKPUYvRsNS2B46sDETSVfuJV1AahRN6aD2WnzQ7wB+S/mqsPXqy+S2zobnU70Wzmt\nhYQzsuIY+BkfOCS565guYvJt66wRGi/NsnybC6z3iRZxZygtKorMKXPjmtdoyCZ3\njkOHhXV5XH8Ldut+1mycg+6c+dTZ9RTrAzo4ouofptm9+ZNlv1aK7HrDDKYQ/x/z\n70hhIfkCgYEA7aNt7db3S6sTEA+yC5DLxxkFIK7K8qxdP6vE087fdHO+ik1LvyI+\ni5Wqj/fT2d/lYR+cbH/zy2JBy5RWfdVJ9HBW/eZ9RqQx8ry7lE5MqbwU6y+d2bCF\nAjffx3yJK1aljuVLdGeu8abYsusUQUhbslZJNg6Ar7z+FUHaChJI0KsCgYEAznk4\nx+5PtvIWj6aXTJiqtofnOtHXXxrq5alzBErvDBHHzP9/ZCucHxJZ4zqHWduj2+9b\nJUgk2BH3+wFMVKpcdTld2iTFWGaFsArTJkE4SBQGx3zcdsoESOnu/DrOTFNu+LiV\nhLzb9CHKM91fIet8MYYxKiyH09+Mi7xBtw8WQwUCgYAH/tCrCOmHHTll9/E4nGWO\nzFO01syzP4NfqgrUSYiRJXfKtXEP/Dn4fk+fymnRUcwo6WRc7i0osaSfEd2bHDsB\nw2nZ3xBl+Q5JKXpyMfQ4XcCibRa1hU/kVDbuQk1nLOIjHandP8POE5wE4Q3saF/V\nbzvFWtWPlB9EXdPVNOpIQwKBgHkJEsIQ72XdUGBxVewu6pQJ4wDWFhzIWL68sJHp\no2w92BRSCkmcTu7gARV1L/b7DHlXPOUD/6UyE15vCmHvZDfLozrHp3AE2YWzMsgQ\nH4ARTVAP3+U603wytkfh6SFRH5JqEiw30fCxBimVMbleo/UcJyID7LPFLkyT1SoM\njA5JAoGACEhaHTXWFYXv+eTJUXFcwhDZ5sRvQYRCTLGSv746lr+SpZ02bEcRvaTH\nv0K1Hph6OzhcCO27VdimngnzsoXoa2OXicWpdjX3hqhXMvwspe9a1y9u7LPZwqJH\n9Kc7VaZP1iS4EIxEc+qx38HqdeiUBcqTMRam6k3mSwactBKCKDI=\n-----END RSA PRIVATE KEY-----\n",
				V3Endpoint: ts.URL + "/",
			},
			expectedAccessToken: "v1.b71be873ad96e64a84025ae7bee7694a99cb4ba9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := resource.GenerateAccessToken(&tt.source, time.Unix(1592691442, 0))
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedAccessToken, got)
		})
	}
}
