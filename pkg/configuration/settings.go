package configuration

import (
	"github.com/imdario/mergo"
	"github.com/r3labs/diff"
	"github.com/stretchy/stretchy/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Settings map[string]interface{}

func (s *Settings) UnmarshalYAML(value *yaml.Node) error {
	yamlSettings := make(map[string]interface{})
	if err := value.Decode(&yamlSettings); err != nil {
		return err
	}

	*s = yamlSettings

	return nil
}

func (s Settings) CleanUp() {
	s.removeMetadataSettings()
	s.moveSettingsUnderIndexKey()
}

func (s Settings) GetIndexSettings() map[string]interface{} {
	indexSettings, exist := s["index"].(map[string]interface{})
	if !exist {
		return make(map[string]interface{})
	}

	return indexSettings
}

func (s Settings) removeMetadataSettings() {
	deleteFromIndex := []string{
		"creation_date",
		"provided_name",
		"uuid",
		"version",
	}

	indexSettings := s.GetIndexSettings()
	if len(indexSettings) == 0 {
		return
	}

	for _, key := range deleteFromIndex {
		if hasKey, _ := utils.MapHasKey(key, indexSettings); hasKey {
			delete(indexSettings, key)
		}
	}

	s["index"] = indexSettings
}

func (s Settings) moveSettingsUnderIndexKey() {
	toBeMovedUnderIndex := []string{
		"number_of_shards",
		"shard.check_on_startup",
		"codec",
		"routing_partition_size",
		"load_fixed_bitset_filters_eagerly",
		"number_of_replicas",
		"auto_expand_replicas",
		"refresh_interval",
		"max_result_window",
		"max_inner_result_window",
		"max_rescore_window",
		"max_docvalue_fields_search",
		"max_script_fields",
		"max_ngram_diff",
		"max_shingle_diff",
		"blocks.read_only",
		"blocks.read_only_allow_delete",
		"blocks.read",
		"blocks.write",
		"blocks.metadata",
		"max_refresh_listeners",
		"highlight.max_analyzed_offset",
		"max_terms_count",
		"routing.allocation.enable",
		"routing.rebalance.enable",
		"gc_deletes",
		"max_regex_length",
		"default_pipeline",
	}

	indexSettings := s.GetIndexSettings()
	if indexSettings == nil {
		indexSettings = make(map[string]interface{})
	}

	for _, key := range toBeMovedUnderIndex {
		if hasKey, _ := utils.MapHasKey(key, s); hasKey == true {
			indexSettings[key] = s[key]
			delete(s, key)
		}
	}

	if len(indexSettings) > 0 {
		s["index"] = indexSettings
	}
}

func (s Settings) Merge(settings Settings) error {
	return mergo.Map(
		&s,
		settings,
		mergo.WithOverride,
	)
}

func (s Settings) Diff(settings Settings) (ChangeCollection, error) {
	settingsChangeLogs, err := diff.Diff(s, settings)
	if err != nil {
		return nil, err
	}

	changes := ChangeCollection{}

	for _, c := range settingsChangeLogs {
		change := Change{
			Type: NewChangeTypeFromDiffType(c.Type),
			Path: append([]string{"settings"}, c.Path...),
			From: c.From,
			To:   c.To,
		}

		if change.ShouldBeReported() {
			changes = append(changes, change)
		}
	}

	return changes, nil
}
