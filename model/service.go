package model

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

type Service struct {
	Name             string        `json:"name"`
	Image            string        `json:"image"`
	Network          string        `json:"network"`
	HealthyManaged   bool          `json:"healthyManaged"`
	RequiredMBMemory int           `json:"requiredMBMemory"`
	EnvVars          []string      `json:"envVars"`
	Volume           []mount.Mount `json:"volumes"`
	Ports            []swarm.PortConfig
}
