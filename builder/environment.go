package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

//
// Description: Builder will build and run the three packages necessary to run.
//
// Example:
//			go run environment.go
//

// RunEnvironment starts everything
func RunEnvironment(installPath string, skipPackageBuild bool) {
	var err error
	c := color.New(color.FgRed)

	packagesToBuild, err := getListOfPackagesToBuild(installPath)
	if err != nil {
		c.Println(err)
		return
	}

	goPathSRC := fmt.Sprintf("%s\\src", os.Getenv("GOPATH"))

	var wg sync.WaitGroup
	// change directory then build then run!
	for _, pkg := range packagesToBuild {
		wg.Add(1)
		packageDir := fmt.Sprintf("%s\\%s", goPathSRC, pkg)

		if localErr := buildPackage(packageDir); localErr != nil {
			c.Printf("Build:\t%s - Fail: %s\n", pkg, localErr)
			break
		} else {
			fmt.Printf("Build:\t%s - Success\n", pkg)
		}

		go func(packageDir, pkg string) {
			fmt.Printf("Run:\t%s\n", pkg)
			if localErr := runProgram(packageDir); localErr != nil {
				c.Printf("Fail:\t%s - %s\n", pkg, localErr)
			}
			wg.Done()
		}(packageDir, pkg)
	}

	wg.Wait()
	fmt.Println("Package Builder End")
}

func currentPath() string {
	dir, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func getListOfPackagesToBuild(installPath string) (packages []string, err error) {
	pkgFilePath := installPath + "/packages_to_build.txt"
	dat, err := ioutil.ReadFile(pkgFilePath)
	if err != nil {
		return packages, err
	}

	fileContent := string(dat)
	if len(fileContent) == 0 {
		return packages, fmt.Errorf("file content is empty for %s", pkgFilePath)
	}

	// TODO: builder, update to better solution of splitting
	return strings.Split(fileContent, "\r\n"), nil
}

func buildPackage(path string) (err error) {
	if err = os.Chdir(path); err != nil {
		return err
	}

	if err = exec.Command("go", "build").Run(); err != nil {
		return err
	}

	return nil
}

func runProgram(path string) (err error) {
	elements := strings.Split(path, "\\")
	executable := fmt.Sprintf("%s\\%s.exe", path, elements[len(elements)-1])

	cmd := exec.Command(executable, "")
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
