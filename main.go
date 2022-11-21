package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schema-cafe/go-types"
	"github.com/schema-cafe/go-types/filesystem"
)

func main() {
	dir := os.Getenv("TS_TYPES_DIR")
	c := &APIClient{
		Endpoint: os.Getenv("API_ENDPOINT"),
	}

	WriteNode(c, dir, "/")
}

func WriteNode(c *APIClient, dir, path string) {
	n, err := c.GetNode(path)
	if err != nil {
		panic(err)
	}

	if n.IsFolder {
		for _, entry := range n.Folder.Contents {
			if entry.IsFolder {
				// TODO: create folder
			}
			newPath := filepath.Join(path, entry.Name)
			WriteNode(c, filepath.Join(dir, newPath), newPath)
		}
	} else {
		// TODO: write .ts file
	}
}

type APIClient struct {
	Endpoint string
}

func (c *APIClient) GetNode(path string) (node *filesystem.Node[types.Schema], err error) {
	res, err := http.Get(c.Endpoint + "?path=" + path)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&node)
	if err != nil {
		return nil, fmt.Errorf("json error: %w", err)
	}

	return node, nil
}
