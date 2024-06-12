package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "time"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
	"context"
	redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

// Entry represents a file or directory entry in the YAML structure
type Entry struct {
	Name      string  `yaml:"-"`
	Extension string  `yaml:"extension,omitempty"`
	ID        string  `yaml:"id,omitempty"`
	Length    int     `yaml:"length,omitempty"`
	Path      string  `yaml:"path,omitempty"`
	Items     []*Entry `yaml:"items,omitempty"`
}

// UnmarshalYAML is a custom unmarshaller for Entry to handle nested structures
func (e *Entry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		// It's a directory with nested items
		for i := 0; i < len(value.Content); i += 2 {
			keyNode := value.Content[i]
			valueNode := value.Content[i+1]
			if keyNode.Value == "extension" {
				e.Extension = valueNode.Value
			} else if keyNode.Value == "id" {
				e.ID = valueNode.Value
			} else if keyNode.Value == "length" {
				fmt.Sscanf(valueNode.Value, "%d", &e.Length)
			} else if keyNode.Value == "path" {
				e.Path = valueNode.Value
			} else {
				e.Name = keyNode.Value
				var items []*Entry
				if err := valueNode.Decode(&items); err != nil {
					return err
				}
				e.Items = items
			}
		}
	} else if value.Kind == yaml.SequenceNode {
		for _, itemNode := range value.Content {
			var item Entry
			if err := itemNode.Decode(&item); err != nil {
				return err
			}
			e.Items = append(e.Items, &item)
		}
	}
	return nil
}

func flattenEntries(entries []*Entry, flattened *[]Entry) {
	for _, entry := range entries {
		*flattened = append(*flattened, *entry)
		if entry.Items != nil {
			flattenEntries(entry.Items, flattened)
		}
	}
}

func main() {
	DB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       6,
		Password: "",
	})
	defer DB.Close()

	app := fiber.New()

	if len(os.Args) < 2 {
		log.Fatal("Please provide the YAML file path")
	}
	yamlFilePath := os.Args[1]
	fmt.Println("Reading YAML file:", yamlFilePath)
	data, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	fmt.Println("YAML file content:", string(data))

	var entries []*Entry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	var flattenedEntries []Entry
	flattenEntries(entries, &flattenedEntries)

	global_key := "testing-go-fil-server"
	ds_key := "testing-go-fil-server-ts"

	var ctx = context.Background()

	DB.Del(ctx, ds_key)
	DB.Del(ctx, ds_key+".INDEX")

	for index, entry := range flattenedEntries {
		if entry.Path == "" {
			continue
		}
		if entry.ID == "" {
			continue
		}
		fmt.Printf("Index: %d , Name: %s, Path: %s , ID: %s\n", index, entry.Name, entry.Path, entry.ID)
		global_entry_key := fmt.Sprintf("%s.%s", global_key, entry.ID)
		DB.Set(ctx, global_entry_key, entry.Path, 0) // minimum need the path, could json blob store here instead
		circular_set.Add(DB, ds_key, entry.ID)
	}

	current := circular_set.Current(DB, ds_key)
	fmt.Printf("Current = %s\n", current)

	app.Get("/files/:uuid.:ext", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		var ctx = context.Background()
		global_entry_key := fmt.Sprintf("%s.%s", global_key, uuid)
		path, err := DB.Get(ctx, global_entry_key).Result()
		fmt.Println( uuid , global_entry_key , path )
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}
		return c.SendFile(path, false)
	})

	log.Fatal(app.Listen(":3000"))
}
