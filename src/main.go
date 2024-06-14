package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

const storageDir string = "storage"

type Path struct {
	BasePath string `yaml:"basepath"`
	Filename string `yaml:"filename"`
}

type Source struct {
	Name    string `yaml:"name"`
	Remote  bool   `yaml:"remote"`
	KeyPath string `yaml:"keypath"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	Paths   []Path `yaml:"paths"`
}

type Config struct {
	Port     int      `yaml:"port"`
	Name     string   `yaml:"name"`
	Interval int      `yaml:"interval"`
	Sources  []Source `yaml:"sources"`
}

var config = Config{}

func init() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Error opening config file")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.Decode(&config)
}

func main() {
	for _, source := range config.Sources {
		go process(source)
	}

	server := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: http.HandlerFunc(handleRequest),
	}

	log.Print("Listening on port " + strconv.Itoa(config.Port))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error")
	}
}
