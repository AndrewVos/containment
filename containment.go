package main

import (
	"errors"
	"fmt"
)

var executer Executer

func main() {
	executer = &SSHExecuter{}
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
				fmt.Sprintf("docker pull %v", container.Image),
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
