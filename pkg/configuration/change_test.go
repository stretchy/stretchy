package configuration_test

import (
	"testing"

	"github.com/r3labs/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
)

func TestNewChangeTypeFromDiffType(t *testing.T) {
	testCases := []struct {
		name               string
		diffType           string
		expectedChangeType configuration.ChangeType
	}{
		{
			name:               "create",
			diffType:           diff.CREATE,
			expectedChangeType: configuration.ChangeTypeCreate,
		},
		{
			name:               "update",
			diffType:           diff.UPDATE,
			expectedChangeType: configuration.ChangeTypeUpdate,
		},
		{
			name:               "delete",
			diffType:           diff.DELETE,
			expectedChangeType: configuration.ChangeTypeDelete,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			changeType := configuration.NewChangeTypeFromDiffType(tc.diffType)
			assert.Equal(t, tc.expectedChangeType, changeType)
		})
	}
}

func TestChange_FullPath(t *testing.T) {
	change := configuration.Change{
		Path: []string{"a", "b", "c", "d"},
	}

	assert.Equal(t, "a.b.c.d", change.FullPath())
}

func TestChange_IsIndexMetadata(t *testing.T) {
	testCases := []struct {
		name            string
		change          configuration.Change
		isIndexMetadata bool
	}{
		{
			name: "settings.index.creation_date",
			change: configuration.Change{
				Path: []string{
					"settings",
					"index",
					"creation_date",
				},
			},
			isIndexMetadata: true,
		},
		{
			name: "settings.index.provided_name",
			change: configuration.Change{
				Path: []string{
					"settings",
					"index",
					"provided_name",
				},
			},
			isIndexMetadata: true,
		},
		{
			name: "settings.index.uuid",
			change: configuration.Change{
				Path: []string{
					"settings",
					"index",
					"uuid",
				},
			},
			isIndexMetadata: true,
		},
		{
			name: "settings.index.version",
			change: configuration.Change{
				Path: []string{
					"settings",
					"index",
					"version",
				},
			},
			isIndexMetadata: true,
		},
		{
			name: "should fail",
			change: configuration.Change{
				Path: []string{
					"mappings",
					"properties",
					"field",
					"type",
				},
			},
			isIndexMetadata: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			isIndexMetaData, err := tc.change.IsIndexMetadata()
			assert.NoError(t, err)
			assert.Equal(
				t,
				tc.isIndexMetadata,
				isIndexMetaData,
			)
		})
	}
}

func TestChange_ShouldBeReported(t *testing.T) {
	testCases := []struct {
		name             string
		change           configuration.Change
		shouldBeReported bool
	}{
		{
			name: "true test",
			change: configuration.Change{
				Path: []string{
					"mappings",
					"properties",
					"field",
					"type",
				},
			},
			shouldBeReported: true,
		},
		{
			name: "default ignore_above on keyword",
			change: configuration.Change{
				Type: configuration.ChangeTypeDelete,
				Path: []string{
					"mappings",
					"properties",
					"fieldName",
					"ignore_above",
				},
				From: float64(256),
				To:   nil,
			},
			shouldBeReported: false,
		},
		{
			name: "default ignore_above on keyword in object",
			change: configuration.Change{
				Type: configuration.ChangeTypeDelete,
				Path: []string{
					"mappings",
					"properties",
					"fieldName",
					"properties",
					"subFieldName",
					"ignore_above",
				},
				From: float64(256),
				To:   nil,
			},
			shouldBeReported: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.change.ShouldBeReported(),
				tc.shouldBeReported,
			)
		})
	}
}

func TestChange_String(t *testing.T) {
	testCases := []struct {
		name           string
		change         configuration.Change
		expectedString string
	}{
		{
			name: "with from to",
			change: configuration.Change{
				Path: []string{
					"mappings",
					"properties",
					"field",
					"type",
				},
				From: "keyword",
				To:   "text",
				Type: configuration.ChangeTypeUpdate,
			},
			expectedString: "UPDATE => mappings.properties.field.type [From: keyword - To: text]",
		},
		{
			name: "without from to",
			change: configuration.Change{
				Path: []string{
					"mappings",
					"properties",
					"field",
				},
				From: nil,
				To: map[string]interface{}{
					"type": "keyword",
				},
				Type: configuration.ChangeTypeCreate,
			},
			expectedString: "CREATE => mappings.properties.field",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedString, testCase.change.String())
		})
	}
}

func TestChangeType_String(t *testing.T) {
	testCases := []struct {
		changeType     configuration.ChangeType
		expectedString string
	}{
		{
			changeType:     configuration.ChangeTypeCreate,
			expectedString: "CREATE",
		},
		{
			changeType:     configuration.ChangeTypeUpdate,
			expectedString: "UPDATE",
		},
		{
			changeType:     configuration.ChangeTypeDelete,
			expectedString: "DELETE",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.expectedString, func(t *testing.T) {
			assert.Equal(
				t,
				tc.changeType.String(),
				tc.expectedString,
			)
		})
	}
}

func TestNewChangeTypeFromDiffType_Unknown(t *testing.T) {
	assert.Panics(
		t,
		func() {
			configuration.NewChangeTypeFromDiffType("8f251818-e652-4895-941a-75be3007509b")
		},
	)
}
