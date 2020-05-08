# Stretchy
![CI](https://github.com/stretchy/stretchy/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/stretchy/stretchy/branch/master/graph/badge.svg)](https://codecov.io/gh/stretchy/stretchy)

Stretchy is an elasticsearch manager. It's main feature is zero-down-time remapping of data.

## Notice :warning:
Currently under development

## Compatibility
 - Elasticsearch 6.x
 - Elasticsearch 7.x
 
## Features
 - Zero-down time remapping of an index (including data-transfer)
 - Initialize indices based on config files
 - Show and compare mapping configurations

## Index and Mapping Configuration
For stretchy to work, your configuration needs to be defined in JSON or YAML format.
Each file should contain a valid structure for index creation (https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-create-index.html).
The name of the file will be used as the base alias name.

## Usage

```bash
stretchy apply --elasticsearch-host=http://localhost:9200 \
    --index-prefix=stretchy \ # A prefix that will be applied on each index and alias
    --path=./configs \ # Path where to search for configurations file
    --format=yaml \ # Format of the configurations file
    --enable-soft-update \ # Allows inplace remapping
    --dry-run # Do not apply changes
```

## Examples

[Some examples can be found here](examples)

## Road to v1
 - [ ] Refactor [configuration](pkg/configuration) package
 - [ ] Create a configuration file based on an existing index
 - [ ] Add metadata on indexes created with stretchy and keep only X old index versions (should be configurable)
 - [ ] Move [utils](pkg/utils) out from this project
 - [ ] Enrich documentation and examples
 - [ ] Create a new repo for homebrew's release
