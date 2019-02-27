package docker

import (
	"context"
	"fmt"
	"github.com/plancks-cloud/plancks-cli/model"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func CreateService(service *model.Service) {
	log.Debugln(fmt.Sprintf("createService method!"))

	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		log.Panicln(fmt.Sprintf("Error getting docker client environment: %s", err))
	}

	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: service.Name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:  service.Image,
				Env:    service.EnvVars,
				Mounts: service.Volume,
			},
			Resources: &swarm.ResourceRequirements{
				Limits: &swarm.Resources{
					MemoryBytes: int64(service.RequiredMBMemory * 1024 * 1024),
				},
			},
			ForceUpdate: 1,
		},
		EndpointSpec: &swarm.EndpointSpec{Ports: service.Ports},
	}

	_, err = cli.ServiceCreate(
		ctx,
		spec,
		types.ServiceCreateOptions{},
	)

	if err != nil {
		log.Errorln(fmt.Sprintf("Error creating docker service: %s", err))
	}
}
