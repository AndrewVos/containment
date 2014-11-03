package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var executer Executer

func main() {
	executer = &SSHExecuter{}

	if _, err := os.Stat("containment.yml"); err != nil {
		log.Fatalln("Couldn't find containment.yml. You're gonna need one of these.")
	}
	b, err := ioutil.ReadFile("containment.yml")
	if err != nil {
		log.Fatalf("Couldn't read containment.yml\n%v\n", err)
	}
	var configuration Configuration
	err = yaml.Unmarshal(b, &configuration)
	if err != nil {
		log.Fatalf("Couldn't decode containment.yml\n%v\n", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "update":
			if len(os.Args) > 2 {
				update(configuration, os.Args[2])
				return
			}
		}
	}
	fmt.Println(`Usage:
  containment status
  containment update IMAGE
  containment start IMAGE
  containment stop IMAGE
  containment restart IMAGE
  containment logs IMAGE
	`)

	fmt.Println("Images:")
	for _, container := range configuration.Containers {
		fmt.Println("  " + container.Image)
	}

	fmt.Println()

	fmt.Println("Clusters:")
	for _, cluster := range configuration.Clusters {
		fmt.Println("  " + cluster.Name)
	}
}

func status() error {
	return nil
}

func update(configuration Configuration, image string) error {
	container, exists := configuration.FindContainerByImageName(image)
	if !exists {
		return errors.New(fmt.Sprintf("Couldn't find container %q\n", image))
	}
	clusters := configuration.FindClustersThatShouldHaveContainer(container)
	if len(clusters) == 0 {
		return errors.New(fmt.Sprintf("Couldn't find a cluster for container %q\n", image))
	}
	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			b, err := executer.Execute(
				host.Address,
				host.Port,
				host.User,
				fmt.Sprintf("sudo docker pull %v", container.Image),
			)
			if err == nil {
				fmt.Println(string(b))
			} else {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func start(configuration Configuration, image string) error {
	return nil
}

func stop(configuration Configuration, image string) error {
	return nil
}

func restart(configuration Configuration, image string) error {
	return nil
}

func logs(configuration Configuration, image string) error {
	return nil
}
