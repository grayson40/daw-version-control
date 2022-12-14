/*
Copyright © 2022 Grayson Crozier <grayson40@gmail.com>
*/
package daw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/grayson40/daw/types"
)

// Write input files to staged file
func ExecuteAdd(input []string) {
	// Throw error if not an initialized repo
	if !IsInitialized() {
		fmt.Println("fatal: not a daw repository (or any of the parent directories): .daw")
		return
	}

	// Throw error if user credentials not configured
	if _, err := os.Stat("./.daw/credentials.json"); err != nil {
		fmt.Println("fatal: user credentials not configured\n  (use \"daw config --username <username> --email <email>\" to configure user credentials)")
		return
	}

	// Throw error if more than one file is inputted for staging
	if len(input) > 1 {
		fmt.Println("fatal: only one project file can be added at a time")
		return
	}

	// Get staged project
	stagedProject := GetStagedProject()

	// Get tracked projects
	trackedProjects, err := GetTracked()

	// Get project file input
	projectFile := input[0]

	// Get file name
	name := projectFile

	// Only want project files
	splitString := strings.Split(name, ".")
	if splitString[1] != "flp" {
		fmt.Printf("fatal: pathspec '%s' is not valid for tracking", name)
		return
	}

	// Check if file exists
	if _, err := os.Stat(name); err != nil {
		fmt.Printf("fatal: pathspec '%s' did not match any files", name)
	} else {
		// Get absolute file path
		path, err := filepath.Abs(name)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Get last modified time
		modTime := GetModifiedTime(name)

		// Add to tracked if untracked
		if !IsTrackedProject(projectFile) {
			trackedProjects = append(trackedProjects, types.File{
				Name:  name,
				Path:  path,
				Saved: modTime,
			})
		}

		// Append file for staging
		if !isStaged(path) {
			var changes []types.Change
			stagedProject = types.Project{
				Name:    name,
				Path:    path,
				Saved:   modTime,
				Changes: changes,
			}
		}
	}

	// Write to tracked json
	err = writeTracked(trackedProjects)
	if err != nil {
		panic(err)
	}

	// Write to staged json
	err = writeStaged(stagedProject)
	if err != nil {
		panic(err)
	}

	// Add to user project files if dne
	// userProjectFiles := requests.GetProjects(currentUserId)
	// userStagedFiles := GetStaged()
	// for _, stagedFile := range userStagedFiles {
	// 	if !projectFileInDb(userProjectFiles, stagedFile.Name) {
	// 		project := types.Project{
	// 			File:    stagedFile,
	// 			Commits: nil,
	// 		}
	// 		userProjectFiles = append(userProjectFiles, project)
	// 	}
	// }

	// Write updated project files to db
	// currentUserId := GetCurrentUser().ID.Hex()
	// requests.AddProject(stagedFiles, currentUserId)
}

// Returns true if file is already staged
func isStaged(filepath string) bool {
	stagedProject := GetStagedProject()
	if stagedProject.Path == filepath {
		return true
	}
	return false
}

// Writes commit array to json file, returns err
func writeStaged(stagedProject types.Project) error {
	file, err := json.MarshalIndent(stagedProject, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./.daw/staged.json", file, 0644)

	return err
}

// Write tracked array to json file
func writeTracked(files []types.File) error {
	file, err := json.MarshalIndent(files, "", "\t")
	if err != nil {
		panic(err)
	}

	writeErr := ioutil.WriteFile("./.daw/tracked.json", file, 0644)

	return writeErr
}
