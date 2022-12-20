package main

import (
  "fmt"
  "net"
  "time"
  "regexp"
  "os/exec"
//  "os"
  "golang.org/x/sys/unix"
  "github.com/jsimonetti/rtnetlink"
)

func main() {
        interfaces, err := net.Interfaces()
        if err != nil {
                fmt.Println(err)
        }

        var pppInt []string
        for i := 0; i < len(interfaces); i++ {
                match, _ := regexp.MatchString("ppp([0-9]+)", interfaces[i].Name)

                if match {
                        pppInt = append(pppInt, interfaces[i].Name)
                }
        }

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

        if len(pppInt) > 0 {
                fmt.Println("There are ppp interfaces(s): ", pppInt)


                /* Still a work-in-progress
                for {
                        for i := 0; i < len(pppInt); i++ {
                                        ppp := testNetwork2(pppInt[i])
                                        if !ppp {
                                                	iface, _ := net.InterfaceByName(pppInt[i])
	                                                attr := rtnetlink.RouteAttributes{
		                                                OutIface: uint32(iface.Index),
                                                                Gateway: net.ParseIP("127.0.0.1"),
                                                                Priority: 100,
	                                                }

                                                        fmt.Println(fmt.Sprintf("ppp device %s cannot connect to internet, deprioritizing!", pppInt[i]))
    	                                                err = conn.Route.Replace(&rtnetlink.RouteMessage{
		                                                Family:     unix.AF_INET,
		                                                Table:      unix.RT_TABLE_MAIN,
		                                                Protocol:   unix.RTPROT_BOOT,
		                                                Scope:      unix.RT_SCOPE_LINK,
		                                                Type:       unix.RTN_UNICAST,
		                                                Attributes: attr,
	                                               })

                                        } else {
                                                        iface, _ := net.InterfaceByName(pppInt[i])
                                                        attr := rtnetlink.RouteAttributes{
                                                                OutIface: uint32(iface.Index),
                                                                Gateway: net.ParseIP("127.0.0.1"),
                                                                Priority: 0,
                                                        }


                                                        //routes, _ := conn.Route.List()
                                                        //fmt.Println("Routes: ", routes[0])
                                                        //address, _ := conn.Address.List()
                                                        //fmt.Println("Address: ", address[1].Attributes.Label)


                                                        fmt.Println(fmt.Sprintf("ppp device %s can access the internet, setting %s metric", pppInt[i], pppInt[i]))
                                                        err = conn.Route.Replace(&rtnetlink.RouteMessage{
                                                                Family:     unix.AF_INET,
                                                                Table:      unix.RT_TABLE_MAIN,
                                                                Protocol:   unix.RTPROT_BOOT,
                                                                Scope:      unix.RT_SCOPE_LINK,
                                                                Type:       unix.RTN_UNICAST,
                                                                Attributes: attr,
                                                       })
                                      }
                                      time.Sleep(5 * time.Second)
                        }
                }

        } */

        if len(pppInt) > 0 {
                fmt.Println("There are ppp interfaces(s): ", pppInt)

                for {
                        for i := 0; i < len(pppInt); i++ {
                                        changeMetric(pppInt[i], "10")
                                        ppp := testNetwork2(pppInt[i])
                                        if ppp {
                                                /* if pppInt[i] == "ppp0" {
                                                        fmt.Println("Setting ppp0 as priority")
                                                        changeMetric(pppInt[i], "20") 
                                                } else { 
                                                        changeMetric(pppInt[i], "30")
                                                        fmt.Println(fmt.Sprintf("ppp device %s can access the internet, setting device as backup", pppInt[i]))
                                                } */

                                                changeMetric(pppInt[i], "30")
                                                fmt.Println(fmt.Sprintf("ppp device %s can access the internet, setting %s metric", pppInt[i], pppInt[i]))
                                        } else {
                                                fmt.Println(fmt.Sprintf("ppp device %s cannot connect to internet, deprioritizing!", pppInt[i]))
                                                changeMetric(pppInt[i], "100")
                                        }
                                        time.Sleep(5 * time.Second)
                       }
                }
        } else {
                fmt.Println("There is no ppp interface")
                os.Exit(1)
        }

        //ppp0, _ := netlink.LinkByName("ppp0")
        //ppp1, _ := netlink.LinkByName("ppp1")
}

func changeMetric(netinterface string, priority string) {
        cmd := exec.Command("ifmetric", netinterface, priority)

        if err := cmd.Start(); err != nil {
                fmt.Println("cmd.Start: %v", err)
        }

        if err := cmd.Wait(); err != nil {
                //if exiterr, ok := err.(*exec.ExitError); ok {
                if _, ok := err.(*exec.ExitError); ok {
                    //fmt.Println("ifmetric exit error: ", exiterr)
                } else {
                    fmt.Println("cmd.Wait: %v", err)
                }
        }
}

func testNetwork2(netinterface string) bool {
        cmd := exec.Command("curl", "--interface", netinterface, "--connect-timeout", "10", "https://8.8.8.8")

        if err := cmd.Start(); err != nil {
                fmt.Println("cmd.Start: %v", err)
                return false
        }

        if err := cmd.Wait(); err != nil {
                //if exiterr, ok := err.(*exec.ExitError); ok {
                if _, ok := err.(*exec.ExitError); ok {
                    //fmt.Println("curl exit error: ", exiterr)
                    return false
                } else {
                    fmt.Println("cmd.Wait: %v", err)
                    return false
                }
        }

        return true
}

func testNetwork(netinterface string) bool {
        nic, err := net.InterfaceByName(netinterface)
        if err != nil {
            fmt.Println(err)
            return false
        }

        adds, err := nic.Addrs()
        if err != nil {
           fmt.Println(err)
           return false
        }

        tcpAddr := &net.TCPAddr{
            IP: adds[0].(*net.IPNet).IP,
        }

        d := net.Dialer{LocalAddr: tcpAddr, Timeout: time.Second * 10}

        _, err = d.Dial("tcp", "8.8.8.8:443")
        if err != nil {
            fmt.Println(err)
            return false
        }

        return true
}
