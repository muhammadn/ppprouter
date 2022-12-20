package main

import (
  "fmt"
  "net"
  "time"
  "regexp"
  "os/exec"
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


                for {
                        for i := 0; i < len(pppInt); i++ {
                                        changeMetric(conn, pppInt[i], 0)
                                        ppp := checkInternet(pppInt[i])
                                        if !ppp {
                                                fmt.Println(fmt.Sprintf("ppp device %s cannot connect to internet, deprioritizing!", pppInt[i]))
                                                changeMetric(conn, pppInt[i], 100)
                                        } else {
                                                fmt.Println(fmt.Sprintf("ppp device %s can access the internet, setting %s metric", pppInt[i], pppInt[i]))
                                                changeMetric(conn, pppInt[i], 0)
                                        }
                                        time.Sleep(20 * time.Second)
                        }
                }

        }
}

func changeMetric(conn *rtnetlink.Conn, netinterface string, metric uint32) {

        iface, _ := net.InterfaceByName(netinterface)
        attr := rtnetlink.RouteAttributes{
                OutIface: uint32(iface.Index),
                Priority: metric,
        }

        err := conn.Route.Replace(&rtnetlink.RouteMessage{
                Family:     unix.AF_INET,
                Table:      unix.RT_TABLE_MAIN,
                Protocol:   unix.RTPROT_BOOT,
                Scope:      unix.RT_SCOPE_LINK,
                Type:       unix.RTN_UNICAST,
                Attributes: attr,
        })

        if err != nil {
                fmt.Println("Error replacing route: ", err)
        }
}

func checkInternet(netinterface string) bool {
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
