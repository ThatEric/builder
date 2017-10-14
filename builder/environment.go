package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/fatih/color"
)

//
// Description: Builder will build and run the three packages necessary to run.
//
// Example:
//			go run environment.go
//

// RunEnvironment starts everything
func RunEnvironment(installPath string, readyChan chan error) {

	var err error
	quit := make(chan bool)
	c := color.New(color.FgRed)

	dat, _ := ioutil.ReadFile(installPath + "/builder_running.txt")
	if len(dat) != 0 {
		c.Println("Environment is already running")
		readyChan <- nil
		return
	}

	txt := []byte("running")
	err = ioutil.WriteFile(installPath+"/builder_running.txt", txt, 0777)

	if err != nil {
		readyChan <- err
		return
	}

	// Cleanup
	defer func() {
		if err = os.Remove(installPath + "/builder_running.txt"); err != nil {
			readyChan <- fmt.Errorf("Remove Failed %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for _ = range signalCh {
			quit <- true
		}
	}()

	packagesToBuild, err := getListOfPackagesToBuild(installPath)
	if err != nil {
		c.Println(err)
		readyChan <- nil
		return
	}

	fmt.Println("Package Builder Start")

	goPathSRC := fmt.Sprintf("%s\\src", os.Getenv("GOPATH"))

	// change directory then build then run!
	for _, pkg := range packagesToBuild {
		var localErr error
		packageDir := fmt.Sprintf("%s\\%s", goPathSRC, pkg)

		if localErr = buildPackage(packageDir); localErr != nil {
			c.Printf("Build:\t%s - Fail: %s\n", pkg, localErr)
			readyChan <- localErr
			break
		} else {
			fmt.Printf("Build:\t%s - Success\n", pkg)
		}

		if len(quit) > 0 {
			break
		}
		time.Sleep(time.Second * 2)

		go func(packageDir, pkg string) {
			fmt.Printf("Run:\t%s\n", pkg)

			if localErr := runProgram(packageDir); localErr != nil {
				c.Printf("Fail:\t%s - %s\n", pkg, localErr)
				readyChan <- localErr
			}
			quit <- true
		}(packageDir, pkg)

		if len(quit) > 0 {
			break
		}
		time.Sleep(time.Second)
	}

	fmt.Println("Waiting for interrupt")
	readyChan <- nil
	<-quit
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

	return strings.Split(fileContent, "\n"), nil
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
