package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type cidr struct {
	address string
	suffix  int
}

func main() {
	var cidrValue cidr
	flag.Var(&cidrValue, "cidr", "CIDR range input e.g. 192.168.0.1/24")
	flag.Parse()
	SolveCidr(cidrValue)
}

func (c *cidr) String() string {
	return fmt.Sprint(*c)
}

func SolveCidr(cidr cidr) cidrResolution {
	bitMaskSize := 32 - cidr.suffix
	subnetMask := generateSubnetMask(bitMaskSize)
	fmt.Println(subnetMask)
	return cidrResolution{}
}

func generateSubnetMask(bitMaskSize int) [4]int {
	octets := [4]int{0, 0, 0, 0}

	for i := 0; i < bitMaskSize; i++ {
		octetIndex := i / 8
		bitPosition := 7 - (i % 8)
		octets[octetIndex] |= 1 << bitPosition
	}

	return octets
}

type cidrResolution struct {
	networkPart      string
	firstHost        string
	broadcastAddress string
	numberOfHosts    int
	subnetMask       string
}

// Yes i could use regex
func (c *cidr) Set(value string) error {
	if value == "" {
		return errors.New("nil CIDR is invalid")
	}
	addrParts := strings.Split(value, "/")
	if len(addrParts) != 2 {
		return errors.New("invalid CIDR block")
	}

	if len(strings.Split(addrParts[0], ".")) != 4 {
		return errors.New("invalid CIDR block")
	}

	suffixPart := addrParts[1][1:]
	bits, err := strconv.Atoi(suffixPart)
	if err != nil {
		return err
	}
	c.address = addrParts[0]
	c.suffix = bits
	return nil
}
