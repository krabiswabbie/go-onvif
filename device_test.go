package onvif

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGetInformation(t *testing.T) {
	log.Println("Test GetInformation")

	res, err := testDevice.GetInformation()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCapabilities(t *testing.T) {
	log.Println("Test GetCapabilities")

	res, err := testDevice.GetCapabilities()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetDiscoveryMode(t *testing.T) {
	log.Println("Test GetDiscoveryMode")

	res, err := testDevice.GetDiscoveryMode()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestGetScopes(t *testing.T) {
	log.Println("Test GetScopes")

	res, err := testDevice.GetScopes()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetHostname(t *testing.T) {
	log.Println("Test GetHostname")

	res, err := testDevice.GetHostname()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestContinuousMove(t *testing.T) {
	AppPTZMove("up")
	time.Sleep(3 * time.Second)
	AppPTZMove("stop")
}
