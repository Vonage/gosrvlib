package jirasrv

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Vonage/gosrvlib/pkg/timeutil"
	"github.com/stretchr/testify/require"
)

func TestComment_MarshalJSON(t *testing.T) {
	t.Parallel()

	dt, err := time.Parse(timeutil.TimeJiraFormat, "2024-06-26T10:20:31.144+0000")
	require.NoError(t, err)

	tests := []struct {
		name string
		data Comment
	}{
		{
			name: "with all fields",
			data: Comment{
				Self: "http://example.invalid/jira/rest/api/2/comment/10000",
				ID:   "10000",
				Author: &User{
					Self:        "http://example.invalid/jira/rest/api/2/user?username=john",
					Name:        "john",
					Key:         "john",
					Email:       "john@example.invalid",
					AvatarURLs:  map[string]string{"48x48": "http://example.invalid/jira/secure/useravatar?size=small&ownerId=john"},
					DisplayName: "John Doe",
					Active:      true,
					TimeZone:    "Australia/Sydney",
				},
				Body:         "This is a comment",
				RenderedBody: "<p>This is a comment</p>",
				UpdateAuthor: &User{
					Self:        "http://example.invalid/jira/rest/api/2/user?username=john",
					Name:        "john",
					Key:         "john",
					Email:       "john@example.invalid",
					AvatarURLs:  map[string]string{"48x48": "http://example.invalid/jira/secure/useravatar?size=small&ownerId=john"},
					DisplayName: "john Doe",
					Active:      true,
					TimeZone:    "Australia/Sydney",
				},
				Created: timeutil.DateTime[timeutil.TJira](dt),
				Updated: timeutil.DateTime[timeutil.TJira](dt),
				Visibility: &Visibility{
					Type:  "role",
					Value: "Administrators",
				},
				Properties: []EntityProperty{
					{
						Key:   "property1",
						Value: "value1",
					},
					{
						Key:   "property2",
						Value: map[string]any{"subkey": "subvalue"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(tt.data)
			require.NoError(t, err)

			var unmarshaled Comment

			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			require.Equal(t, tt.data, unmarshaled)
		})
	}
}
