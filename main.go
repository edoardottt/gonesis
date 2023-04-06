/*
gonesis

Create Golang project template ready to be pushed on GitHub with a single command

https://github.com/edoardottt/gonesis

edoardottt, https://www.edoardoottavianelli.it

Under GNU-GPL3: https://github.com/edoardottt/gonesis/blob/main/LICENSE
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	Permission0755 = 0755
	Permission0775 = 0775
	MDBashInit     = "```bash"
	Version        = "1.0.1"
	Banner         = "gonesis v" + Version + "\n\thttps://github.com/edoardottt/gonesis\n\n"
)

var (
	ErrProjectName = errors.New("the project name can contains only alphanumeric characters, _ and -")
	ErrGoModExists = errors.New("go.mod already exists")
)

func main() {
	fmt.Print(Banner)

	projectName, err := ProjectName()
	if err != nil {
		log.Fatal(err)
	}

	rootDir := "." + string(os.PathSeparator) + projectName
	err = CreateDir(".", projectName)

	if err != nil {
		log.Fatal(err)
	}

	// description.
	description := Description()

	// cmd.
	err = CreateDir(rootDir, "cmd")
	if err != nil {
		log.Fatal(err)
	}

	// main.
	CreateMain(rootDir, projectName)

	// go.mod.
	name := GithubHandle()
	cmd := exec.Command("go", "mod", "init", "github.com/"+name+"/"+projectName)

	cmd.Dir = "./" + projectName

	err = cmd.Run()
	if err != nil {
		log.Fatal(ErrGoModExists)
	}

	oldLocation := "./go.mod"
	newLocation := rootDir + string(os.PathSeparator) + "go.mod"

	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Fatal(err)
	}

	// pkg.
	err = CreateDir(rootDir, "pkg")
	if err != nil {
		log.Fatal(err)
	}

	// docs.
	err = CreateDir(rootDir, "docs")
	if err != nil {
		log.Fatal(err)
	}

	// internal.
	err = CreateDir(rootDir, "internal")
	if err != nil {
		log.Fatal(err)
	}

	// examples.
	err = CreateDir(rootDir, "examples")
	if err != nil {
		log.Fatal(err)
	}

	// api.
	if AskUser("Will you need APIs?") {
		err = CreateDir(rootDir, "api")
		if err != nil {
			log.Fatal(err)
		}
	}

	// server.
	if AskUser("Will you need a server?") {
		err = CreateDir(rootDir, "server")
		if err != nil {
			log.Fatal(err)
		}
	}

	// db.
	if AskUser("Will you need a database?") {
		err = CreateDir(rootDir, "db")
		if err != nil {
			log.Fatal(err)
		}
	}

	// README.md.
	Readme(rootDir, projectName, description, name)

	// .gitignore.
	Gitignore(rootDir)
}

// ProjectName takes as input from stdin the name of the
// project and checks if it's a valid name.
func ProjectName() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Name of the project: ")

	project, _ := reader.ReadString('\n')
	if len(project) > 0 && project[len(project)-1] == '\n' {
		project = project[:len(project)-1]
	}

	isProjectNameOkay := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !isProjectNameOkay.Match([]byte(project)) {
		return project, fmt.Errorf("%w", ErrProjectName)
	}

	return project, nil
}

// CreateMain creates the main.go file.
func CreateMain(rootDir string, projectName string) {
	err := CreateFile(rootDir+string(os.PathSeparator)+"cmd", projectName+".go")
	if err != nil {
		log.Fatal(err)
	}

	main := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
`
	err = WriteFile(rootDir+string(os.PathSeparator)+"cmd"+string(os.PathSeparator)+projectName+".go", main)

	if err != nil {
		log.Fatal(err)
	}
}

// GithubHandle takes as input from stdin the github profile name.
func GithubHandle() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Github username: ")

	name, _ := reader.ReadString('\n')

	if len(name) > 0 && name[len(name)-1] == '\n' {
		name = name[:len(name)-1]
	}

	return name
}

// Description takes as input from stdin the project description.
func Description() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Project description: ")

	desc, _ := reader.ReadString('\n')
	if len(desc) > 0 && desc[len(desc)-1] == '\n' {
		desc = desc[:len(desc)-1]
	}

	return desc
}

// AskUser prints the question taken as input and if
// the input is y/Y or yes/Yes/YES etc. returns true,
// false otherwise.
func AskUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[ ? ] " + question + " ")

	answer, _ := reader.ReadString('\n')
	if len(answer) > 0 && answer[len(answer)-1] == '\n' {
		answer = answer[:len(answer)-1]
	}

	answer = strings.ToLower(answer)
	if answer == "y" || answer == "yes" {
		return true
	}

	return false
}

// Gitignore creates the .gitignore file.
func Gitignore(rootDir string) {
	err := CreateFile(rootDir, ".gitignore")
	if err != nil {
		log.Fatal(err)
	}

	gitignore := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
	
# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work`

	err = WriteFile(rootDir+string(os.PathSeparator)+".gitignore", gitignore)
	if err != nil {
		log.Fatal(err)
	}
}

// Readme creates the README.md file.
func Readme(rootDir string, projectName string, description string, name string) {
	err := CreateFile(rootDir, "README.md")
	if err != nil {
		log.Fatal(err)
	}

	readme := "# " + projectName
	readme += "\n" + description
	readme += "\n\nInstallation 📡\n"
	readme += "-------\n"
	readme += "**Go 1.17+**\n"
	readme += MDBashInit + "\n"
	readme += "go install -v github.com/" + name + "/" + projectName + "@latest\n"
	readme += "```\n"
	readme += "**otherwise**\n"
	readme += MDBashInit + "\n"
	readme += "go get -v github.com/" + name + "/" + projectName + "\n"
	readme += "```\n\n"
	readme += "Usage 💻\n"
	readme += "-------\n"
	readme += MDBashInit + "\n"
	readme += projectName + "\n"
	readme += "```\n\n"
	readme += "Created with [gonesis](https://github.com/edoardottt/gonesis)❤️"

	err = WriteFile(rootDir+string(os.PathSeparator)+"README.md", readme)
	if err != nil {
		log.Fatal(err)
	}
}

//----------------------------------------
//---------------- helpers ---------------
//----------------------------------------

// CreateDir creates the directory with the name
// taken as input.
func CreateDir(path string, name string) error {
	err := os.MkdirAll(path+string(os.PathSeparator)+name, Permission0775)
	return err
}

// CreateFile creates the file with the name
// taken as input.
func CreateFile(path string, name string) error {
	err := os.WriteFile(path+string(os.PathSeparator)+name, []byte(""), Permission0755)
	return err
}

// WriteFile writes the content string into the file taken as input.
func WriteFile(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), Permission0755)
	return err
}
