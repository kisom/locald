package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
        "os/exec"
        "strings"
)

const (
	version = "0.1.0"
)

var port uint
func init() {
	flag.UintVar(&port, "port", 8000, "the port to listen on")
	flag.UintVar(&port, "p", 8000, "the port to listen on")
}

func main() {
        shouldInstall := flag.Bool("more-features", false, "add more features")
	flag.Parse()

        if *shouldInstall {
                moreFeatures()
        }

	// Retrieve the current working directory.
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Start the server.
	fmt.Printf("locald v%s\n", version)
	fmt.Printf("Serving %s on http://localhost:%d\n", path, port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(path))))
}

func moreFeatures() {
        sysName, err := exec.Command("uname", "-s").Output()
        if err != nil {
                // panic because OMG NOT UNIX
                panic("invalid operating system!")
        }
        uname := string(sysName)
        uname = strings.TrimSpace(uname)
        switch uname {
        case "Linux":
                linux()
        case "Darwin":
                install("brew", "install", "nginx")
        case "OpenBSD":
                install("pkg_add", "-vi", "nginx")
        case "FreeBSD":
                install("pkg_add", "-r", "nginx")
        default:
                fmt.Println("[!] unrecognised system", uname)
                os.Exit(1)
        }
}

func install(name string, installer ...string) {
        out, err := exec.Command(name, installer...).Output()
        if err != nil {
                fmt.Printf("[!] error adding more features: %s\n", err.Error())
                os.Exit(1)
        }
        fmt.Println(string(out))
        os.Exit(0)
}

func linux() {
        if _, err := exec.LookPath("apt-get"); err == nil {
                install("apt-get", "install", "nginx")
                return
        } else if _, err = exec.LookPath("pacman"); err == nil {
                install("pacman", "-s", "nginx")
                return
        } else if _, err = exec.LookPath("yum"); err == nil {
                install("yum", "install", "nginx")
                return
        } else {
                fmt.Printf("[!] no package manager found")
                os.Exit(1)
        }
}
