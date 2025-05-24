package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// marshal2OrderedMap reads YAML from an io.Reader and unmarshals it into a yaml.MapSlice.
// It uses goccy/go-yaml which should preserve key quoting.
// This function might not be needed if we directly manipulate AST.
func marshal2OrderedMap(r io.Reader) (*yaml.MapSlice, error) {
	var outMap yaml.MapSlice
	err := yaml.NewDecoder(r).Decode(&outMap)
	if err != nil {
		return nil, err
	}

	return &outMap, nil
}

// merge merges the source map into the destination map.
// Existing keys in dst are overwritten by src, and new keys from src are added to dst.
// This function might not be needed if we directly manipulate AST.
func merge(src, dst *yaml.MapSlice) {
	dstMap := make(map[interface{}]int)
	for i, item := range *dst {
		dstMap[item.Key] = i
	}

	for _, srcItem := range *src {
		if i, ok := dstMap[srcItem.Key]; ok {
			// Key exists in dst, update value
			(*dst)[i].Value = srcItem.Value
		} else {
			// Key does not exist in dst, append it
			*dst = append(*dst, srcItem)
		}
	}
}

// ProcessFiles reads, merges, and writes YAML files preserving destination key quoting.
func ProcessFiles(srcPath, dstPath string) error {
	// Read source file
	srcBytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Read destination file
	dstBytes, err := ioutil.ReadFile(dstPath)
	if err != nil {
		return fmt.Errorf("failed to read destination file: %w", err)
	}

	// Parse source and destination YAML into ASTs
	srcAST, err := parser.ParseBytes(srcBytes, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse source YAML: %w", err)
	}

	dstAST, err := parser.ParseBytes(dstBytes, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse destination YAML: %w", err)
	}

	// Get the root map node for both ASTs
	srcMap, ok := srcAST.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return fmt.Errorf("source YAML is not a map")
	}
	dstMap, ok := dstAST.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		return fmt.Errorf("destination YAML is not a map")
	}

	// Create a map to store destination keys and their original nodes for quoting preservation
	dstKeyNodes := make(map[string]*ast.MappingValueNode)
	for _, valueNode := range dstMap.Values {
		if keyNode, ok := valueNode.Key.(*ast.StringNode); ok {
			dstKeyNodes[keyNode.Value] = valueNode
		}
	}

	// Merge source into destination AST
	for _, srcValueNode := range srcMap.Values {
		if srcKeyNode, ok := srcValueNode.Key.(*ast.StringNode); ok {
			keyStr := srcKeyNode.Value

			// Check if the key exists in the destination
			if dstValueNode, exists := dstKeyNodes[keyStr]; exists {
				// Key exists, update the value in the destination AST
				dstValueNode.Value = srcValueNode.Value
			} else {
				// Key does not exist, add the new key-value pair to the destination AST
				dstMap.Values = append(dstMap.Values, srcValueNode)
			}
		}
	}

	// Write the modified destination AST back to the file
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open destination file for writing: %w", err)
	}
	defer dstFile.Close()

	// Use AST's String() method to get the YAML output string
	output := dstAST.String()

	_, err = dstFile.WriteString(output)
	if err != nil {
		return fmt.Errorf("failed to write merged YAML to file: %w", err)
	}

	return nil
}

// entrypoint contains the main logic of the program.
func entrypoint() int {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <source_yaml_file> <destination_yaml_file>")
		return 1
	}

	srcPath := args[0]
	dstPath := args[1]

	err := ProcessFiles(srcPath, dstPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		return 1
	}

	fmt.Println("YAML files merged successfully!")
	return 0
}

func main() {
	os.Exit(entrypoint())
}
