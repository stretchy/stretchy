package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Masterminds/semver"
	"github.com/stretchy/stretchy/pkg/configuration"
)

// Client is an Elasticsearch client. Create one by calling New.
type Client interface {
	IndexExist(indexName string) (bool, error)
	AliasExist(aliasName string) (bool, error)

	CreateIndex(indexName string, mapping configuration.Index) error
	CreateAlias(aliasName string, indexName string) error

	UpdateAlias(aliasName string, newIndexName string) error

	GetAliasedIndex(aliasName string) (string, error)
	GetIndexConfiguration(indexName string) (configuration.Index, error)

	UpdateIndexConfiguration(indexName string, configuration configuration.Index) error

	Reindex(sourceIndexName string, targetIndexName string) error
}

const v6ClientMajor int64 = 6
const v7ClientMajor int64 = 7

// New creates a Client instance.
func New(options Options) (Client, error) {
	version, err := getElasticsearchVersion(options)
	if err != nil {
		return nil, err
	}

	switch version {
	case v6ClientMajor:
		return NewV6Client(options)
	case v7ClientMajor:
		return NewV7Client(options)
	}

	return nil, fmt.Errorf("version '%d' not supported yet", version)
}

type cluster struct {
	Version version `json:"version"`
}

type version struct {
	Number string `json:"number"`
}

func isAValidHost(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		return fmt.Errorf("missing scheme")
	}

	if u.Host == "" {
		return fmt.Errorf("missing host")
	}

	return nil
}

func getElasticsearchVersion(options Options) (int64, error) {
	if err := isAValidHost(options.Host); err != nil {
		return 0, fmt.Errorf("elasticsearch host: %s", err)
	}

	req, err := http.NewRequest("GET", options.Host, nil)
	if err != nil {
		return 0, err
	}

	if options.User != "" && options.Password != "" {
		req.SetBasicAuth(options.User, options.Password)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error retrieving elasticsearch version: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	clusterInfo := cluster{}

	if err := json.Unmarshal(body, &clusterInfo); err != nil {
		return 0, err
	}

	v, err := semver.NewVersion(clusterInfo.Version.Number)
	if err != nil {
		return 0, err
	}

	return v.Major(), nil
}
