package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type cidr struct {
	address [4]int
	suffix  int
}

func main() {
	var cidrValue cidr
	flag.Var(&cidrValue, "cidr", "CIDR range input e.g. 192.168.0.1/24")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("please provide a CIDR using the -cidr flag")
		return
	}

	resolution := SolveCidr(cidrValue)
	fmt.Println("Network Address:", resolution.networkAddress)
	fmt.Println("Subnet Mask:", resolution.subnetMask)
	fmt.Println("First Host:", resolution.firstHost)
	fmt.Println("Last Host:", resolution.lastHost)
	fmt.Println("Broadcast Address:", resolution.broadcastAddress)
}

func (c *cidr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d/%d", c.address[0], c.address[1], c.address[2], c.address[3], c.suffix)
}

func SolveCidr(c cidr) cidrResolution {
	subnetMask := generateSubnetMask(c.suffix)

	networkAddress := [4]int{}
	for i := range c.address {
		networkAddress[i] = c.address[i] & subnetMask[i]
	}

	broadcastAddress := [4]int{}
	invertedMask := [4]int{}
	for i := range subnetMask {
		invertedMask[i] = 255 - subnetMask[i]
		broadcastAddress[i] = networkAddress[i] | invertedMask[i]
	}

	firstHost := [4]int{
		networkAddress[0],
		networkAddress[1],
		networkAddress[2],
		networkAddress[3],
	}

	lastHost := [4]int{
		broadcastAddress[0],
		broadcastAddress[1],
		broadcastAddress[2],
		broadcastAddress[3],
	}

	if c.suffix < 31 {
		firstHost[3]++
		lastHost[3]--
	}

	return cidrResolution{
		networkAddress:   formatIP(networkAddress),
		subnetMask:       formatIP(subnetMask),
		firstHost:        formatIP(firstHost),
		lastHost:         formatIP(lastHost),
		broadcastAddress: formatIP(broadcastAddress),
	}
}

func generateSubnetMask(bits int) [4]int {
	result := [4]int{0, 0, 0, 0}

	for i := 0; i < bits; i++ {
		octetIndex := i / 8
		bitPosition := 7 - (i % 8)
		result[octetIndex] |= 1 << bitPosition
	}

	return result
}

func formatIP(ip [4]int) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

type cidrResolution struct {
	networkAddress   string
	firstHost        string
	lastHost         string
	broadcastAddress string
	numberOfHosts    int
	subnetMask       string
}

func (c *cidr) Set(value string) error {
	if value == "" {
		return errors.New("nil CIDR is invalid")
	}
	addrParts := strings.Split(value, "/")
	if len(addrParts) != 2 {
		return errors.New("invalid CIDR block")
	}

	address := strings.Split(addrParts[0], ".")
	if len(address) != 4 {
		return errors.New("invalid CIDR block")
	}

	suffixPart := addrParts[1]
	bits, err := strconv.Atoi(suffixPart)
	if err != nil {
		return err
	}

	if bits < 0 || bits > 32 {
		return errors.New("invalid subnet mask size (must be between 0 and 32)")
	}

	for i := range address {
		addrPart, err := strconv.Atoi(address[i])
		if err != nil {
			return err
		}
		if addrPart < 0 || addrPart > 255 {
			return errors.New("invalid IP address octet (must be between 0 and 255)")
		}
		c.address[i] = addrPart
	}
	c.suffix = bits
	return nil
}
