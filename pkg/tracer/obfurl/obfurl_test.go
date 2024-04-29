//go:build unit

package obfurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObscureUrl(t *testing.T) {
	t.Parallel()

	type testCase struct {
		path         string
		expectedPath string
	}
	cases := []testCase{
		{path: "/entities", expectedPath: "/entities"},
		{path: "/entities/", expectedPath: "/entities/"},
		{path: "/entities/123", expectedPath: "/entities/{ID}"},
		{path: "/entities/123/", expectedPath: "/entities/{ID}/"},
		{path: "/entities/123/test", expectedPath: "/entities/{ID}/test"},
		{path: "/entities/123/test/", expectedPath: "/entities/{ID}/test/"},
		{path: "/entities/123/test/765", expectedPath: "/entities/{ID}/test/{ID}"},
		{path: "/entities/123/test/765/", expectedPath: "/entities/{ID}/test/{ID}/"},
		{path: "/entities/53500e10-a535-4d14-8cc0-846451f47f26", expectedPath: "/entities/{UUID}"},
		{path: "/entities/53500e10-a535-4d14-8cc0-846451f47f26/", expectedPath: "/entities/{UUID}/"},
		{path: "/entities/53500e10-a535-4d14-8cc0-846451f47f26/test", expectedPath: "/entities/{UUID}/test"},
		{path: "/entities/53500e10-a535-4d14-8cc0-846451f47f26/test/", expectedPath: "/entities/{UUID}/test/"},
		{
			path:         "/entities/53500e10-a535-4d14-8cc0-846451f47f26/test/53500e10-a535-4d14-8cc0-846451f47f26",
			expectedPath: "/entities/{UUID}/test/{UUID}",
		},
		{
			path:         "/entities/53500e10-a535-4d14-8cc0-846451f47f26/test/53500e10-a535-4d14-8cc0-846451f47f26/",
			expectedPath: "/entities/{UUID}/test/{UUID}/",
		},
		{
			path:         "/entities/53500e10/test/53500E10-A535-4D14-8CC0-846451F47F26/",
			expectedPath: "/entities/53500e10/test/{UUID}/",
		},
		{
			path:         "/entities/53500e10-a535-4d14-8cc0-846451f47f2g/test/",
			expectedPath: "/entities/53500e10-a535-4d14-8cc0-846451f47f2g/test/",
		},
	}

	for _, testCase := range cases {
		t.Run("Obscuring path", func(t *testing.T) {
			obscuredUrl := ObfuscateURL(testCase.path)

			assert.Equal(t, testCase.expectedPath, obscuredUrl)
		})
	}
}
