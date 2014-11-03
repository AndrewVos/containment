package main

type Configuration struct {
	Clusters   []Cluster
	Containers []Container
}

func (c *Configuration) FindContainerByImageName(image string) (Container, bool) {
	for _, container := range c.Containers {
		if container.Image == image {
			return container, true
		}
	}
	return Container{}, false
}

func (c *Configuration) FindClustersThatShouldHaveContainer(container Container) []Cluster {
	var clusters []Cluster
	for _, clusterName := range container.Clusters {
		for _, cluster := range c.Clusters {
			if cluster.Name == clusterName {
				clusters = append(clusters, cluster)
			}
		}
	}
	return clusters
}

type Host struct {
	Address string
	Port    int
	User    string
}

type Cluster struct {
	Name  string
	Hosts []Host
}

type Container struct {
	Image    string
	Clusters []string
}
