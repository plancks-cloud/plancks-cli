package pc

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/plancks-cloud/plancks-cli/docker"
	"github.com/plancks-cloud/plancks-cli/model"
	"github.com/plancks-cloud/plancks-cli/util"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Install() {

	welcomeUser()
	installNetwork()
	installHealth()

	time.Sleep(3 * time.Second)

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
	exists, err := docker.CheckNetworkExists(docker.NetworkName)
	if err != nil {
		logrus.Fatalln("Could not check if the network exists")
	}

	if !exists {
		success, err := docker.CreateOverlayNetwork(docker.NetworkName)

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

	adminPort := swarm.PortConfig{TargetPort: model.ListentPort, PublishedPort: model.AppPort}
	proxyPort := swarm.PortConfig{TargetPort: 6228, PublishedPort: 6228}
	ports := []swarm.PortConfig{adminPort, proxyPort}

	service := model.Service{
		Name:             "pc",
		Image:            "planckscloud/plancks-cloud:latest",
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

	url := fmt.Sprintf("http://127.0.0.1:%v/route", model.AppPort)

	answered := false
	attempts := 0
	for !answered {

		bytes, err := util.GetRequest(url)

		if bytes != nil && err == nil {
			logrus.Println("...Healthy!")
			answered = true
			break
		}

		if err != nil {
			logrus.Error(fmt.Sprintf("Checking the health service faild with an error"))
			logrus.Error(err.Error())
			logrus.Error(fmt.Sprintf("Check %v", url))
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
