package test

import (
	"testing"
	"fmt"
	"net"
	"log"
)

func TestMac(t *testing.T)  {
	inter, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(inter.HardwareAddr.String())
}
