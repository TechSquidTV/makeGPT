package config

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

//go:embed assets/avatar.png
var avatar embed.FS

type FileStructure struct {
	BaseDir     string
	Directories map[string]string
	Files       map[string]FileDetails
}

type FileDetails struct {
	Path    string
	Content []byte
}

type GPTConfigFile struct {
	Name         string   `yaml:"name" json:"name"`
	ID           string   `yaml:"id" json:"id"`
	Description  string   `yaml:"description" json:"description"`
	Instructions string   `yaml:"instructions" json:"instructions"`
	Starters     []string `yaml:"starters" json:"starters"`
}

type GPTConfig struct {
	Config        GPTConfigFile
	FileStructure *FileStructure
}

var defaultKnowledgeArticle = `## Knowledge Directory
Use the knowledge directory (knowledge/) to store files that will be accessible to your GPT. 
This allows you to provide additional specialized or more current information on topics to help improve responses.

### What to include
- Documents
	- Text files
	- Markdown (.md) files
	- PDFs
- Data
  - CSV files
  - JSON files
  - XML files

### Tips!
- Pre-process your data to remove unnessecary information, styling markup, etc. This will help your GPT understand the core information faster.
- Vectorized data relates data based on "distance".
	- If your data is relatively flat, markdown tables or CSV files may work even better than JSON.
- For data that should be live, consider using an Action, for relatively static data sharing as a knowledge source will be faster.

## Actions
Actions in GPTs connect an OpenAPI (Swagger) spec to your GPT. This allows your GPT to respond to live API requests and queries.
Currently, OpenAI only allows for a single 'action'.
Add your OpenAPI spec file to the 'actions/' directory and it will be automatically included.

## avatar.png
In the root of your project should be an image named "avatar.png" that will be used as the avatar/profile picture for your GPT on the OpenAI website and API.

## makeGPT
makeGPT is the CLI tool used to generate this GPT project.
makeGPT provides a way for developers to manage and develop their custom GPTs with git, and deploy with CI/CD pipelines.

### Commands
- 'makegpt init [path]' - Initializes a new GPT project with the necessary files and configuration
- 'makegpt deploy' - Deploys the current GPT to OpenAI

### Using Git
makeGPT can be downloaded as a part of your CI pipeline and used to deploy new versions to OpenAI.
In your CI/CD pipeline, configure a job named 'deploy' to run on either tags or releases depending on your preferences.
Within the deploy job, run the 'makegpt deploy' command to deploy the currently checked-out version of the GPT to OpenAI.
`

var defaultAction = `{
  "openapi": "3.1.0",
  "info": {
    "title": "Get weather data",
    "description": "Retrieves current weather data for a location.",
    "version": "v1.0.0"
  },
  "servers": [
    {
      "url": "https://api.example-weather-app.com"
    }
  ],
  "paths": {
    "/location": {
      "get": {
        "description": "Get temperature for a specific location",
        "operationId": "GetCurrentWeather",
        "parameters": [
          {
            "name": "location",
            "in": "query",
            "description": "The city and state to retrieve the weather for",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "deprecated": false
      }
    }
  },
  "components": {
    "schemas": {}
  }
}
`

// Create new config object
func NewConfig(path string) *GPTConfig {
	config := &GPTConfig{
		Config: GPTConfigFile{
			Name:         "",
			ID:           "",
			Description:  "",
			Instructions: "",
			Starters:     []string{"how do manage my GPT with git?", "explain the `makegpt init` command.", "What are best practices for adding knowledge?", "how do I add an API action?"},
		},
	}
	config.FileStructure = config.NewFileStructure(path)
	return config
}

// Create file structure object
func (c *GPTConfig) NewFileStructure(baseDir string) *FileStructure {
	return &FileStructure{
		BaseDir: baseDir,
		Directories: map[string]string{
			"knowledge": "knowledge",
			"actions":   "actions",
		},
		Files: map[string]FileDetails{
			"ActionJSON": {
				Path:    "actions/weather.json.example",
				Content: []byte(defaultAction),
			},
			"KnowledgeArticle": {
				Path:    "knowledge/makeGPT.md",
				Content: []byte(defaultKnowledgeArticle),
			},
			"Avatar": {
				Path:    "avatar.png",
				Content: readAvatar(),
			},
			"Config": {
				Path:    "gpt.yaml",
				Content: c.ToYaml(),
			},
		},
	}
}

// Create directories
func (c *GPTConfig) CreateDirectories() {
	for _, dir := range c.FileStructure.Directories {
		err := os.MkdirAll(filepath.Join(c.FileStructure.BaseDir, dir), os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create directory", err)
		}
	}
}

// Create files
func (c *GPTConfig) CreateFiles() {
	for _, file := range c.FileStructure.Files {
		f, err := os.Create(filepath.Join(c.FileStructure.BaseDir, file.Path))
		if err != nil {
			log.Fatalf("Failed to create to file \"%s\"\n%s", file, err)
		}
		defer f.Close()
		_, err = f.Write(file.Content)
		if err != nil {
			log.Fatalf("Failed to write to file \"%s\"\n%s", file, err)
		}
	}
}

func readAvatar() []byte {
	avatar, err := avatar.ReadFile("assets/avatar.png")
	if err != nil {
		log.Fatalf("Error reading internal default avatar image: \n%v\n", err)
	}
	return avatar
}

func (c *GPTConfig) ToYaml() []byte {
	yaml, err := yaml.Marshal(c.Config)
	if err != nil {
		log.Fatalf("Failed to marshal config to yaml: \n%v\n", err)
	}
	return yaml
}

func LoadConfig(path string) (*GPTConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config GPTConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
