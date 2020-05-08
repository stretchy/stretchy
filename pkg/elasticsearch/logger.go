package elasticsearch

import (
	"log"
	"os"
)

func newLogger() *log.Logger {
	return log.New(os.Stdout, "ElasticSearchClient ", log.LstdFlags)
}
