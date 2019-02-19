package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	command = "lsblk -n -o name,size -J"
	key     = "blockdevices"
	devPath = "/dev/"

	pattern = `sd[a-z]`
)

type SATADevice struct {
	Name     string                   `json:"name"`
	Children []map[string]interface{} `json:"children"`
	Path     string                   `json:"path"`
	Size     string                   `json:"size"`
}

type SATADeviceSlice []SATADevice

func GetDevices() (SATADeviceSlice, error) {
	data, err := runCombinedOutput(command)
	if err != nil {
		return nil, err
	}
	return parseSATADevices(data)
}

func GetDevicesWithData(data string) (SATADeviceSlice, error) {
	return parseSATADevices(data)
}

func parseSATADevices(data string) (SATADeviceSlice, error) {
	var deviceMarshal map[string]interface{}
	if err := json.Unmarshal([]byte(data), &deviceMarshal); err != nil {
		return nil, err
	}

	bds, ok := deviceMarshal[key]
	if !ok {
		return nil, fmt.Errorf("Parse block device JSON error")
	}

	raw, err := json.Marshal(bds)
	if err != nil {
		return nil, err
	}

	var sds SATADeviceSlice
	if err := json.Unmarshal(raw, &sds); err != nil {
		return nil, err
	}

	var newSATADevices SATADeviceSlice
	for _, device := range sds {
		m, _ := regexp.MatchString(pattern, device.Name)
		if len(device.Children) == 0 && m {
			device.Path = devPath + device.Name
			newSATADevices = append(newSATADevices, device)
		}
	}
	glog.V(4).Infof("Get devices: %s", newSATADevices)
	return newSATADevices, nil
}

func runCombinedOutput(cmd string) (string, error) {
	glog.V(4).Infof("Run with output:", cmd)
	c := exec.Command("/bin/bash", "-c", cmd)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "running command: %s\n output: %s", cmd, out)
	}
	return string(out), nil
}
