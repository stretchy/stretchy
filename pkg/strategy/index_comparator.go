package strategy

import (
	"github.com/stretchy/stretchy/pkg/configuration"
)

type IndexAction int

const (
	IndexDecisionNone IndexAction = iota
	IndexDecisionCreate
	IndexDecisionMigrate
	IndexDecisionUpdate
)

func (id IndexAction) String() string {
	return []string{"None", "Create", "Migrate", "Update"}[id]
}

type IndexActionVoter struct {
	allowSoftUpdate bool
}

func NewIndexActionVoter(allowSoftUpdate bool) *IndexActionVoter {
	return &IndexActionVoter{
		allowSoftUpdate: allowSoftUpdate,
	}
}

type IndexVoterResult struct {
	action  IndexAction
	changes configuration.ChangeCollection
}

func NewIndexVoterResult(
	action IndexAction,
	changes configuration.ChangeCollection,
) IndexVoterResult {
	return IndexVoterResult{
		action:  action,
		changes: changes,
	}
}

func (indexIndexComparatorResult IndexVoterResult) Action() IndexAction {
	return indexIndexComparatorResult.action
}

func (indexIndexComparatorResult IndexVoterResult) Changes() configuration.ChangeCollection {
	return indexIndexComparatorResult.changes
}

func (ic *IndexActionVoter) Compare(
	currentConfiguration *configuration.Index,
	newConfiguration *configuration.Index,
) (IndexVoterResult, error) {
	if currentConfiguration == nil {
		return IndexVoterResult{
			action:  IndexDecisionCreate,
			changes: nil, // Maybe we should add here the whole new mapping
		}, nil
	}

	changes, err := currentConfiguration.Diff(*newConfiguration)

	if err != nil {
		return IndexVoterResult{}, err
	}

	if len(changes) == 0 {
		return NewIndexVoterResult(
			IndexDecisionNone,
			nil,
		), nil
	}

	if ic.allowSoftUpdate && ic.canBeASoftUpdate(changes) {
		return NewIndexVoterResult(
			IndexDecisionUpdate,
			changes,
		), nil
	}

	return NewIndexVoterResult(
		IndexDecisionMigrate,
		changes,
	), nil
}

// A soft update should be possible only when there are only properties addition
func (ic IndexActionVoter) canBeASoftUpdate(changes configuration.ChangeCollection) bool {
	for _, c := range changes {
		if c.Path[0] == "settings" {
			return false
		}

		// This is meant to check if the change is a full new property added on:
		//  - Base Document
		//  - Inside a property with type object
		//  - Inside a property with type nested
		if c.Path[0] == "mappings" {
			if c.Path[len(c.Path)-2] != "properties" {
				return false
			} else if c.Type != configuration.ChangeTypeCreate {
				return false
			}
		}
	}

	return true
}
