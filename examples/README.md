# Examples

In this directory you can find some different configuration examples:

Base:

In this example we are going to create a new alias named `stretchy-author` and the corresponding index `stretchy-author-{timestamp}`.

### YAML 
```bash
stretchy apply --elasticsearch-host=http://localhost:9200 \
    --index-prefix=stretchy \
    --path=./configs-yaml \
    --format=yaml
```

### JSON
```bash
stretchy apply --elasticsearch-host=http://localhost:9200 \
    --index-prefix=stretchy \
    --path=./configs-json \
    --format=json
```

The command is idemponent so if you don't change anything in the configuration nothing should happen on consecutive run.
