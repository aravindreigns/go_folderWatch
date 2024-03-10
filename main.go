//Task - Go Program to keep a watch on a folder. Any new file changes are detected and based on the regular expression fetched from the config file the result is displayed in console.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var regexPattern *regexp.Regexp

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME") // add home directory as search path
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}

	regexPatternString := viper.GetString("regexPattern")
	var err error
	regexPattern, err = regexp.Compile(regexPatternString)
	if err != nil {
		log.Fatal("Error compiling regex pattern:", err)
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	folderToWatch := viper.GetString("folderToWatch")
	err = filepath.Walk(folderToWatch, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return watcher.Add(path)
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Watching folder: %s\n", folderToWatch)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("File %s has been modified.\n", event.Name)

					// Read the file content
					content, err := readFile(event.Name)
					if err != nil {
						log.Println("Error reading file:", err)
					} else {
						fmt.Println("File content:")
						matches := regexPattern.FindAllString(content, -1)
						for _, match := range matches {
							fmt.Println(match)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	<-done
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
