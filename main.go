package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	// Fork used because original version is not providing forwarded ports info
	"github.com/bsideup/go-virtualbox"
	"github.com/fsouza/go-dockerclient"
)

const forward2dockerPrefix = "forward2docker-"

var runOnce = flag.Bool("run-once", false, "stop execution after the first port assignment")

func main() {
	flag.Parse()
	client, err := docker.NewClientFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	reload(client)

	if !*runOnce {
		listener := make(chan *docker.APIEvents)
		client.AddEventListener(listener)
		for {
			select {
			case event := <-listener:
				switch event.Status {
				case "start", "die":
					reload(client)
				}
			}
		}
	}
}

func reload(client *docker.Client) {
	machineName := os.Getenv("DOCKER_MACHINE_NAME")

	if len(machineName) == 0 {
		machineName = "boot2docker-vm"
	}

	machine, err := virtualbox.GetMachine(machineName)

	if err != nil {
		log.Fatal(err)
	}
	forwards := make(map[int64]bool)

	for _, forwarding := range machine.Forwardings {
		if !strings.HasPrefix(forwarding.Name, forward2dockerPrefix) {
			continue
		}

		forwards[int64(forwarding.GuestPort)] = true
	}

	log.Printf("Previously mapped: %v\n", forwards)

	opts := docker.ListContainersOptions{}
	containers, err := client.ListContainers(opts)

	if err != nil {
		log.Println("Failed to list containers", err)
		return
	}

	for _, container := range containers {
		for _, binding := range container.Ports {
			_, exists := forwards[binding.PublicPort]
			if !exists {
				if binding.PublicPort == 0 {
					log.Println("Will not add mapping for port 0")
					continue
				}
				log.Printf("Adding port mapping for %d\n", binding.PublicPort)
				if binding.PublicPort < 1024 {
					log.Printf("Exposed port %d is lower than 1024 and will not work if VirtualBox is started without root privileges\n", binding.PublicPort)
				}
				rule := virtualbox.PFRule{
					Proto:     virtualbox.PFTCP,
					HostPort:  uint16(binding.PublicPort),
					GuestPort: uint16(binding.PublicPort),
				}
				machine.AddNATPF(1, forward2dockerPrefix+strconv.FormatInt(binding.PublicPort, 10), rule)
			}

			delete(forwards, binding.PublicPort)
		}
	}

	if len(forwards) > 0 {
		log.Printf("To be deleted: %v\n", forwards)
		for port := range forwards {
			machine.DelNATPF(1, forward2dockerPrefix+strconv.FormatInt(port, 10))
		}
	}

}
