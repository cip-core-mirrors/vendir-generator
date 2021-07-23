package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed templates/vendir.tpl
var vendirTemplate string

type CipConfig struct {
	Components []struct {
		BaseDestinationPath string `yaml:"baseDestinationPath"`
		Source              []struct {
			SourcePath string
			Modules    []struct {
				SourceName      string `yaml:"sourceName"`
				DestinationName string `yaml:"destinationName"`
			}
			GitConfig struct {
				Url string
				Ref string
			} `yaml:"gitConfig"`
		}
	}
}

type VendirDirectory struct {
	Path    string
	Content []VendirDirectoryContent
}

type VendirDirectoryContent struct {
	Path         string
	NewRootPath  string
	Git          VendirDirectoryContentGit
	IncludePaths string
}

type VendirDirectoryContentGit struct {
	Url string
	Ref string
}

func main() {
	args := os.Args

	log.SetFlags(log.Ldate | log.Ltime)

	if len(args) < 2 {
		log.Fatal("An argument must be provide and specifiy the path to the config file")
	}

	errConfig, config := loadConfig(os.Args[1])

	if errConfig != nil {
		log.Fatal(errConfig.Error())
	}

	_, vendir := configToVendirDirectories(config)

	errTemplate := vendirToTemplate(vendir)

	if errTemplate != nil {
		log.Fatal(errTemplate.Error())
	}
}

func loadConfig(filename string) (error, CipConfig) {
	config := CipConfig{}
	file, errOpen := ioutil.ReadFile(filename)

	if errOpen != nil {
		return errors.New("Cannot open your file"), config
	}

	errorParse := yaml.Unmarshal(file, &config)

	if errorParse != nil {
		return errors.New("Cannot parse your file, check format"), config
	}

	return nil, config
}

func configToVendirDirectories(config CipConfig) (error, []VendirDirectory) {
	allVendirDirectory := []VendirDirectory{}

	for _, component := range config.Components {
		vendirDirectory := VendirDirectory{}

		vendirDirectory.Path = component.BaseDestinationPath

		for _, source := range component.Source {

			gitConfig := VendirDirectoryContentGit{
				Url: source.GitConfig.Url,
				Ref: source.GitConfig.Ref,
			}

			for _, module := range source.Modules {

				content := VendirDirectoryContent{}

				content.NewRootPath = filepath.Join(source.SourcePath, module.SourceName)
				content.IncludePaths = filepath.Join(source.SourcePath, module.SourceName, "**", "*")
				content.Path = module.DestinationName
				content.Git = gitConfig

				vendirDirectory.Content = append(vendirDirectory.Content, content)
			}
		}

		allVendirDirectory = append(allVendirDirectory, vendirDirectory)
	}

	return nil, allVendirDirectory
}

func vendirToTemplate(allDirectory []VendirDirectory) error {
	template, errLoad := template.New("main").Parse(vendirTemplate)

	if errLoad != nil {
		return errors.New("Failed to load result template")
	}

	f, errCreateResultFile := os.Create("../vendir.yml")

	if errCreateResultFile != nil {
		return errors.New("Failed to create result file")
	}

	return template.ExecuteTemplate(f, "vendir", allDirectory)
}
