package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	projectName, err := ProjectName()
	if err != nil {
		log.Fatal(err)
	}
	rootDir := "." + string(os.PathSeparator) + projectName
	CreateDir(".", projectName)

	//description
	description := Description()

	//cmd
	CreateDir(rootDir, "cmd")

	//main
	CreateMain(rootDir, projectName)

	//go.mod
	name := GithubHandle()
	cmd := exec.Command("go", "mod", "init", "github.com/"+name+"/"+projectName)
	err = cmd.Run()
	if err != nil {
		log.Fatal("go.mod already exists in this folder?")
	}
	oldLocation := "./go.mod"
	newLocation := rootDir + string(os.PathSeparator) + "go.mod"
	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Fatal(err)
	}

	//pkg
	CreateDir(rootDir, "pkg")

	//docs
	CreateDir(rootDir, "docs")

	//internal
	CreateDir(rootDir, "internal")

	//examples
	CreateDir(rootDir, "examples")

	//api
	if AskUser("Will you need API?") {
		CreateDir(rootDir, "api")
	}

	//web
	if AskUser("Will you need a web server?") {
		CreateDir(rootDir, "web")
	}

	//db
	if AskUser("Will you need a database?") {
		CreateDir(rootDir, "db")
	}

	//README.md
	Readme(rootDir, projectName, description, name)

	//.gitignore
	Gitignore(rootDir)

}

//ProjectName takes as input from stdin the name of the
//project and checks if it's a valid name.
func ProjectName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Name of the project: ")
	project, _ := reader.ReadString('\n')
	if len(project) > 0 && project[len(project)-1] == '\n' {
		project = project[:len(project)-1]
	}
	isProjectNameOkay := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !isProjectNameOkay.Match([]byte(project)) {
		return project, errors.New("the project name can contains only alphanumeric characters, _ and -")
	}

	return project, nil
}

//CreateMain creates the main.go file.
func CreateMain(rootDir string, projectName string) {
	CreateFile(rootDir+string(os.PathSeparator)+"cmd", projectName+".go")
	main := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
`
	WriteFile(rootDir+string(os.PathSeparator)+"cmd"+string(os.PathSeparator)+projectName+".go", main)
}

//GithubHandle takes as input from stdin the github profile name.
func GithubHandle() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Github username: ")
	name, _ := reader.ReadString('\n')
	if len(name) > 0 && name[len(name)-1] == '\n' {
		name = name[:len(name)-1]
	}
	return name
}

//Description takes as input from stdin the project description.
func Description() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Project description: ")
	desc, _ := reader.ReadString('\n')
	if len(desc) > 0 && desc[len(desc)-1] == '\n' {
		desc = desc[:len(desc)-1]
	}
	return desc
}

//AskUser prints the question taken as input and if
//the input is y/Y or yes/Yes/YES etc. returns true,
//false otherwise.
func AskUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " ")
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

//Gitignore creates the .gitignore file.
func Gitignore(rootDir string) {
	CreateFile(rootDir, ".gitignore")
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
	WriteFile(rootDir+string(os.PathSeparator)+".gitignore", gitignore)
}

//Readme creates the README.md file.
func Readme(rootDir string, projectName string, description string, name string) {
	CreateFile(rootDir, "README.md")
	readme := "# " + projectName
	readme = readme + "\n" + description
	readme = readme + "\n\nInstallation 📡\n"
	readme = readme + "-------\n"
	readme = readme + "**Go 1.17+**\n"
	readme = readme + "```bash\n"
	readme = readme + "go install -v github.com/" + name + "/" + projectName + "@latest\n"
	readme = readme + "```\n"
	readme = readme + "**otherwise**\n"
	readme = readme + "```bash\n"
	readme = readme + "go get -v github.com/" + name + "/" + projectName + "\n"
	readme = readme + "```\n\n"
	readme = readme + "Usage 💻\n"
	readme = readme + "-------\n"
	readme = readme + "```bash\n"
	readme = readme + projectName + "\n"
	readme = readme + "```\n\n"
	readme = readme + "Created with [gonesis](https://github.com/edoardottt/gonesis)❤️"
	WriteFile(rootDir+string(os.PathSeparator)+"README.md", readme)
}

//----------------------------------------
//---------------- helpers ---------------
//----------------------------------------

//CreateDir creates the directory with the name
//taken as input.
func CreateDir(path string, name string) (bool, error) {
	err := os.MkdirAll(path+string(os.PathSeparator)+name, 0775)
	success := false
	if err == nil {
		success = true
	}
	return success, err
}

//CreateFile creates the file with the name
//taken as input.
func CreateFile(path string, name string) (bool, error) {
	err := ioutil.WriteFile(path+string(os.PathSeparator)+name, []byte(""), 0755)
	success := false
	if err == nil {
		success = true
	}
	return success, err
}

//WriteFile writes the content string into the file taken as input.
func WriteFile(filename string, content string) (bool, error) {
	err := ioutil.WriteFile(filename, []byte(content), 0755)
	success := false
	if err == nil {
		success = true
	}
	return success, err
}
