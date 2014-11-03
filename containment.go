package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
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
			if len(os.Args) == 3 {
				err := update(configuration, os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		case "start":
			if len(os.Args) == 3 {
				err := start(configuration, os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		case "stop":
			if len(os.Args) == 3 {
				err := stop(configuration, os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
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

func findContainerAndClusters(configuration Configuration, image string) (Container, []Cluster, error) {
	container, exists := configuration.FindContainerByImageName(image)
	if !exists {
		return Container{}, nil, errors.New(fmt.Sprintf("Couldn't find container %q\n", image))
	}
	clusters := configuration.FindClustersThatShouldHaveContainer(container)
	if len(clusters) == 0 {
		return Container{}, nil, errors.New(fmt.Sprintf("Couldn't find a cluster for container %q\n", image))
	}
	return container, clusters, nil
}

func hostIdentifier(host Host) string {
	return fmt.Sprintf("[%v@%v] ", host.User, host.Address)
}

func update(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	outputs := make(chan string)
	var waitGroup sync.WaitGroup

	updateImage := func(container Container, host Host) {
		identifier := hostIdentifier(host)
		b, err := executer.Execute(
			host.Address,
			host.Port,
			host.User,
			fmt.Sprintf("sudo docker pull %v", container.Image),
		)

		if err == nil {
			for _, s := range strings.Split(string(b), "\n") {
				outputs <- fmt.Sprintf("%v%v", identifier, s)
			}
		} else {
			outputs <- fmt.Sprintf("%v%v", identifier, err.Error())
		}
		waitGroup.Done()
	}

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			waitGroup.Add(1)
			go updateImage(container, host)
		}
	}
	go func() {
		for o := range outputs {
			fmt.Println(o)
		}
	}()
	waitGroup.Wait()
	close(outputs)

	return nil
}

func start(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	command := fmt.Sprintf("sudo docker run -d --name %v ", container.Name())
	for _, port := range container.Ports {
		command += fmt.Sprintf("-p %v ", port)
	}
	command += container.Image

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			identifier := hostIdentifier(host)
			b, err := executer.Execute(
				host.Address,
				host.Port,
				host.User,
				command,
			)
			if err == nil {
				for _, s := range strings.Split(string(b), "\n") {
					fmt.Printf("%v%v\n", identifier, s)
				}
			} else {
				fmt.Printf("%v%v\n", identifier, err.Error())
			}
		}
	}

	return nil
}

func stop(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			identifier := hostIdentifier(host)
			command := fmt.Sprintf("sudo docker stop %v && sudo docker rm %v", container.Name(), container.Name())
			b, err := executer.Execute(
				host.Address,
				host.Port,
				host.User,
				command,
			)
			if err == nil {
				for _, s := range strings.Split(string(b), "\n") {
					fmt.Printf("%v%v\n", identifier, s)
				}
			} else {
				fmt.Printf("%v%v\n", identifier, err.Error())
			}
		}
	}

	return nil
}

func restart(configuration Configuration, image string) error {
	return nil
}

func logs(configuration Configuration, image string) error {
	return nil
}
