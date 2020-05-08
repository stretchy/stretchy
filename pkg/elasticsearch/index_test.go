package elasticsearch_test

import (
	"fmt"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
)

func TestCreateIndexName(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		time              time.Time
		aliasName         string
		expectedIndexName string
	}{
		{
			time:              now,
			aliasName:         "alias-a",
			expectedIndexName: fmt.Sprintf("alias-a-%d", now.Unix()),
		},
		{
			time:              time.Date(1989, 12, 18, 21, 17, 00, 91, time.UTC),
			aliasName:         "alias-b",
			expectedIndexName: "alias-b-630019020",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.aliasName, func(t *testing.T) {
			patch := monkey.Patch(time.Now, func() time.Time { return testCase.time })
			defer func() {
				fmt.Println(testCase.aliasName)
				patch.Unpatch()
			}()

			assert.Equal(
				t,
				testCase.expectedIndexName,
				elasticsearch.CreateIndexName(testCase.aliasName),
			)
		})
	}
}
func TestResolveAliasName(t *testing.T) {
	testCases := []struct {
		name              string
		indexName         string
		prefix            string
		expectedAliasName string
	}{
		{
			name:              "without prefix",
			indexName:         "my-index",
			prefix:            "",
			expectedAliasName: "my-index",
		},
		{
			name:              "with prefix",
			indexName:         "my-index",
			prefix:            "my-prefix",
			expectedAliasName: "my-prefix-my-index",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(
				t,
				testCase.expectedAliasName,
				elasticsearch.ResolveAliasName(
					testCase.prefix,
					testCase.indexName,
				),
			)
		})
	}
}
