package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/sina-ghaderi/nanontp/engine"
	"github.com/sina-ghaderi/nanontp/getter"
	"github.com/sina-ghaderi/nanontp/network"
)

const regIPPORT string = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):()([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$"
const regDomainIPport string = "^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|((\\*\\.)?([a-zA-Z0-9-]+\\.){0,5}[a-zA-Z0-9-][a-zA-Z0-9-]+\\.[a-zA-Z]{2,63}?)):([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$"

func main() {
	flag.Usage = usage
	ipup := flag.String("net", "0.0.0.0:123", "udp network to listen on <ipv4:port>")
	flag.Parse()
	if match, _ := regexp.MatchString(regIPPORT, *ipup); !match {
		fmt.Printf("fatal: %v is not a valid <ipv4:port> address\n", *ipup)
		flag.Usage()
		os.Exit(2)
	}
	getter.ArgNTP = flag.Args()

	if len(getter.ArgNTP) == 0 {
		fmt.Println("fatal: you must give me a upstream ntp server <ipv4-or-domain:port>")
		flag.Usage()
		os.Exit(2)

	}
	for _, x := range getter.ArgNTP {
		if match, _ := regexp.MatchString(regDomainIPport, x); !match {
			fmt.Printf("fatal: %v is not a valid ntp upstream address\n", x)
			flag.Usage()
			os.Exit(2)
		}
	}

	var handler = network.GetHandler()
	engine.Reactor.ListenUdp(*ipup, handler)
	engine.Reactor.Run()
}

func usage() {
	fmt.Printf(`usage of snix ntp server:
%v -net [ipv4:port] ntp-domain.com:port ntp-domain.org:port ...

options:
  --net string     udp network to listen on <ipv4:port> (default "0.0.0.0:123")
  --h              print this banner and exit
example: 
  %v --net 0.0.0.0:123 time.google.com:123 ntp.day.ir:123 10.10.10.10:123 

Copyright (c) 2020 slc.snix.ir, All rights reserved.
Developed BY a.esmaeilpour@irisaco.com And s.ghaderi1999@gmail.com
This work is licensed under the terms of the MIT license.
`, os.Args[0], os.Args[0])
}
