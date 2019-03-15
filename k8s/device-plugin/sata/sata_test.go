package main

import (
	"testing"
)

const lsblkData = `
{
   "blockdevices": [
      {"name": "sda", "size": "931.5G",
         "children": [
            {"name": "sda1", "size": "931.5G"}
         ]
      },
      {"name": "sdb", "size": "931.5G"}
   ]
}
`

func TestGetDevices(t *testing.T) {
	devices, err := GetDevicesWithData(lsblkData)
	if err != nil {
		t.Fatalf("Could not get devices, %s.", err.Error())
	}
	t.Logf("Get devices: %s.", devices)
}
