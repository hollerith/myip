package main

import (
    "fmt"
    "net"
    "os"
    "syscall"
    "unsafe"
)

func getIPAddress(ifname string) (string, error) {
    // Open a socket using AF_INET and SOCK_DGRAM
    sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
    if err != nil {
        return "", err
    }
    defer syscall.Close(sock)

    // Pack the interface name into a buffer
    var name [16]byte
    copy(name[:], ifname)

    // Send an SIOCGIFADDR request to the host's operating system
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(sock), uintptr(syscall.SIOCGIFADDR), uintptr(unsafe.Pointer(&name)))
    if errno != 0 {
        return "", fmt.Errorf("Error %d", errno)
    }

    // Convert the raw bytes of the IP address to a human-readable string
    ip := (*syscall.RawSockaddrInet4)(unsafe.Pointer(&name)).Addr[:]
    return net.IPv4(ip[0], ip[1], ip[2], ip[3]).String(), nil
}

func main() {
    // Get the interface name from the command-line arguments, or use the default if not specified
    var ifname string
    if len(os.Args) > 1 {
        ifname = os.Args[1]
    } else {
        // Find the default network interface
        ifaces, err := net.Interfaces()
        if err != nil {
            fmt.Println(err)
            return
        }
        for _, iface := range ifaces {
            if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
                ifname = iface.Name
                break
            }
        }
    }

    // Get the IP address of the specified interface
    ip, err := getIPAddress(ifname)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(ip)
    }
}

