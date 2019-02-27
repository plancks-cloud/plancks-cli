package main

import (
	"fmt"
	"github.com/plancks-cloud/plancks-docker"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func install() {

	welcomeUser()
	installNetwork()
	installHealth()
	checkHealth()
	installComplete()
	os.Exit(0)

}

//welcomeUser tells the user what is about to happen
func welcomeUser() {
	logrus.Infoln(fmt.Sprintf("Welcome to the Planck's Cloud installer."))
}

//install network tries to setup pc-net and tells the user if it succeeded
func installNetwork() {
	logrus.Infoln(fmt.Sprintf(".. Attempting to create overlay network"))
	exists, err := plancks_docker.CheckNetworkExists(plancks_docker.model.NetworkName)
	if err != nil {
		logrus.Fatalln("Could not check if the network exists")
	}

	if !exists {
		success, err := plancks_docker.CreateOverlayNetwork(plancks_docker.NetworkName)

		if err != nil {
			logrus.Fatalln("Could not check if the network exists")
		}

		if !success {
			logrus.Fatalln("Create network was not successful")
		}
	}

	logrus.Infoln(fmt.Sprintf(".. ✅ Success"))

}

//installHealth starts the Health docker image as service on the pc-net network and tells the user if it worked out
func installHealth() {
	logrus.Infoln(fmt.Sprintf(".. Attempting to install health service"))

	envVars := []string{"MODE=NORMAL"}
	mounts := []mount.Mount{{Source: "/var/run/docker.sock", Target: "/var/run/docker.sock"}}
	ports := []swarm.PortConfig{{TargetPort: model.ListentPort, PublishedPort: model.AppPort}}

	service := model.Service{
		Name:             "pc-health",
		Image:            "planckscloud/pc-health:latest",
		Network:          "plancks-net",
		RequiredMBMemory: 64,
		EnvVars:          envVars,
		Volume:           mounts,
		Ports:            ports}

	docker.CreateService(&service)

	//Install health service
	//if err != nil {
	//	logrus.Error(fmt.Sprintf("Failed to install health service. Install failed. Shutting down."))
	//	panic(0)
	//}
	logrus.Infoln(fmt.Sprintf(".. ✅ Success"))

}

//checks that the health service can be contacted in the browser by calling it's health service
func checkHealth() {

	url := fmt.Sprintf("http://localhost:%v", model.AppPort)

	answered := false
	attempts := 0
	for !answered {
		bytes, err := util.GetRequest(url)

		if bytes != nil {
			answered = true
			break
		}

		if err != nil {
			logrus.Error(fmt.Sprintf("Checking the health service faild with an error, %s", err.Error()))
			if attempts < model.InstallMaxHealthChecks {
				logrus.Infoln(fmt.Sprintf("Will try again in a few seconds"))
				time.Sleep(time.Duration(model.InstallSleepBetweenChecks) * time.Second)
				continue
			}
		}

		attempts++
		if attempts > model.InstallMaxHealthChecks {
			break
		}

	}

	if answered {
		logrus.Infoln(fmt.Sprintf(".. ✅ Success"))
	} else {
		logrus.Fatalln(fmt.Sprintf(".. Was not able to find Health service in %v attampts", model.InstallMaxHealthChecks))
	}

}

func installComplete() {
	logrus.Infoln(fmt.Sprintf("The installation completed succesfully."))
}
