package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const Name = "test-network"

func TestNetwork(t *testing.T) {

	cleanupForTest()

	//Create a network
	success, err := CreateOverlayNetwork(Name)
	assert.True(t, success, "The network create method should return success is true")
	assert.Nil(t, err, "The error should be nil")

	//Check if it exists
	exists, err := CheckNetworkExists(Name)
	assert.True(t, exists, "The network exists method should be able to find the network")
	assert.Nil(t, err, "The error should be nil")

	//Delete the network
	deleted, err := DeleteNetwork(Name)
	assert.True(t, deleted, "The network delete method should say that it was deleted")
	assert.Nil(t, err, "The error should be nil")

	//Network exists
	nowExists, err := CheckNetworkExists(Name)
	assert.False(t, nowExists, "The network exists method should NOT be able to find the network")
	assert.Nil(t, err, "The error should be nil")

}

func cleanupForTest() {
	DeleteNetwork(Name)
}
