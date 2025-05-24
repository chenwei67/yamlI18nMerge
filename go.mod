module yamlI18nMerge

go 1.21.0

toolchain go1.24.3

require github.com/goccy/go-yaml v1.17.1 // Use the latest version

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Replace gopkg.in/yaml.v2 with goccy/go-yaml
// require gopkg.in/yaml.v2 v2.4.0 // indirect
