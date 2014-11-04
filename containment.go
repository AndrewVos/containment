package main

import (
	"bufio"
	"bytes"
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
		case "restart":
			if len(os.Args) == 3 {
				err := restart(configuration, os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		case "status":
			if len(os.Args) == 3 {
				err := status(configuration, os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}

	fmt.Println(`Usage:
  containment status IMAGE
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

func executeCommandAndWriteOutput(container Container, host Host, command string) error {
	b, err := executer.Execute(host, command)
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		fmt.Printf("%v %v\n", host.Identifier(), scanner.Text())
	}
	return err
}

func status(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			waitGroup.Add(1)
			go func(container Container, host Host) {
				command := fmt.Sprintf("sudo docker inspect -f '{{.State.Running}}' %v", container.Name())
				b, err := executer.Execute(host, command)

				if err == nil {
					status := "running"
					if strings.TrimSpace(string(b)) != "true" {
						status = "stopped"
					}
					fmt.Printf("%v %v %v\n", host.Identifier(), container.Image, status)
				} else {
					scanner := bufio.NewScanner(bytes.NewReader(b))
					for scanner.Scan() {
						fmt.Printf("%v %v\n", host.Identifier(), scanner.Text())
					}

				}
				waitGroup.Done()
			}(container, host)
		}
	}
	waitGroup.Wait()

	return nil
}

func update(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			waitGroup.Add(1)
			go func(container Container, host Host) {
				command := fmt.Sprintf("sudo docker pull %v", container.Image)
				output, err := executer.Execute(host, command)
				if err == nil {
					fmt.Printf("%v Updated %v\n", host.Identifier(), container.Image)
				} else {
					fmt.Printf("%v Failed to update %v\n%v\n", host.Identifier(), container.Image, string(output))
				}
				waitGroup.Done()
			}(container, host)
		}
	}
	waitGroup.Wait()

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
			output, err := executer.Execute(host, command)
			if err == nil {
				fmt.Printf("%v Started %v\n", host.Identifier(), container.Image)
			} else {
				fmt.Printf("%v Failed to start %v\n%v\n", host.Identifier(), container.Image, string(output))
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
			command := fmt.Sprintf("sudo docker stop %v && sudo docker rm %v", container.Name(), container.Name())
			output, err := executer.Execute(host, command)
			if err == nil {
				fmt.Printf("%v Stopped %v\n", host.Identifier(), container.Image)
			} else {
				fmt.Printf("%v Failed to stop %v\n%v\n", host.Identifier(), container.Image, string(output))
			}
		}
	}
	return nil
}

func restart(configuration Configuration, image string) error {
	container, clusters, err := findContainerAndClusters(configuration, image)
	if err != nil {
		return err
	}

	runCommand := fmt.Sprintf("sudo docker run -d --name %v ", container.Name())
	for _, port := range container.Ports {
		runCommand += fmt.Sprintf("-p %v ", port)
	}
	commands := []string{
		fmt.Sprintf("sudo docker stop %v", container.Name()),
		fmt.Sprintf("sudo docker rm %v", container.Name()),
		runCommand,
	}
	command := strings.Join(commands, " && ")
	command += container.Image

	for _, cluster := range clusters {
		for _, host := range cluster.Hosts {
			executeCommandAndWriteOutput(container, host, command)
		}
	}

	return nil
}

func logs(configuration Configuration, image string) error {
	return nil
}
