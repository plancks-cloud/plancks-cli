package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"log"
)

const NetworkName = "plancks-net"

//CreateOverlayNetwork creates an overlay network in docker swarm
func CreateOverlayNetwork(name string) (success bool, err error) {
	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		log.Panicln(fmt.Sprintf("Error getting docker client environment: %s", err))
		return false, err
	}

	res, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{Driver: "overlay", Attachable: true})

	logrus.Infoln(fmt.Sprintf(res.ID))
	logrus.Infoln(fmt.Sprintf(res.Warning))

	if err != nil {
		logrus.Infoln(fmt.Sprintf(err.Error()))
		return false, err
	}
	success = true
	return

}

//CheckNetworkExists tells us if a network name exists
func CheckNetworkExists(name string) (exists bool, err error) {
	exists, _, err = describeNetwork(name)
	return

}

//DeleteNetwork removes a network by name
func DeleteNetwork(name string) (success bool, err error) {
	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		log.Println(fmt.Sprintf("Error getting docker client environment: %s", err))
		return false, err
	}

	exists, ID, err := describeNetwork(name)
	if err != nil {
		return false, err
	}
	if !exists {
		success = false
		return
	}
	err = cli.NetworkRemove(ctx, ID)

	if err != nil {
		return false, err
	}
	success = true
	return

}

func describeNetwork(name string) (exists bool, ID string, err error) {
	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		log.Panicln(fmt.Sprintf("Error getting docker client environment: %s", err))
		return false, "", err
	}

	list, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if len(list) == 0 {
		return false, "", err
	}

	for _, network := range list {
		if network.Name == name {
			return true, network.ID, err
		}
	}
	exists = false
	return
}
