package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"vessel/api"
	"vessel/client"
	"vessel/common"
	"vessel/manager"
	"vessel/runtime"
	"vessel/server"
)

const DEBUG = false

func invalid() {
	fmt.Println("Invalid Usage")
	os.Exit(1)
}

func createVessel() {
	vesselClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	appDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	vesselId, err := vesselClient.NewVessel("test", appDir, "alpine", "buildfile", "/bin/sh", "/app/entrypoint")
	if err != nil {
		panic(err)
	}
	fmt.Println("Vessel Created:", vesselId)
	vesselClient.AddHandler(common.UPDATE_RELAY_STDOUT, client.NewRelayedOutputHandler(stdOutRelayer))
	vesselClient.AddHandler(common.UPDATE_RELAY_STDERR, client.NewRelayedOutputHandler(stdErrRelayer))
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println("Error while reading:", err)
				return
			}
			log.Println("You typed: ", string(buf[:n]))
			err = vesselClient.RelayStdin(vesselId, buf[:n])
			if err != nil {
				log.Println("Failed to relay stdin:", err)
			}
			log.Println("Are you deadlocked?")
		}
	}()
	err = vesselClient.Listen()
	if err != nil {
		fmt.Println("Failed to start listener:", err.Error())
	}
}

func stopVessel(id string) {
	vesselClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	stopped, err := vesselClient.StopVessel(id)
	if stopped {
		fmt.Println("Stopped Successfully....")
		return
	}
	fmt.Println("Failed to stop:", err.Error())
}

func stdOutRelayer(vesselId string, data []byte) error {
	_, err := os.Stdout.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func stdErrRelayer(vesselId string, data []byte) error {
	_, err := os.Stderr.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func startVessel(id string) {
	vesselClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	_, err = vesselClient.StartVessel(id)
	if err != nil {
		fmt.Println("Failed to start:", err.Error())
		return
	}
	fmt.Println("Started Successfully....")
	vesselClient.AddHandler(common.UPDATE_RELAY_STDOUT, client.NewRelayedOutputHandler(stdOutRelayer))
	vesselClient.AddHandler(common.UPDATE_RELAY_STDERR, client.NewRelayedOutputHandler(stdErrRelayer))
	go func() {
		buf := make([]byte, 1024)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("Error while reading:", err)
			return
		}
		err = vesselClient.RelayStdin(id, buf[:n])
		if err != nil {
			log.Println("Failed to relay stdin:", err)
		}
	}()
	err = vesselClient.Listen()
	if err != nil {
		fmt.Println("Failed to start listener:", err.Error())
	}
}

func startPreIsolatedVessel(args []string) {
	if len(args) > 2 {
		vesselId := args[2]
		var build = false
		if len(args) > 3 && args[3] == "true" {
			build = true
		}
		err := runtime.StartVessel(vesselId, build, DEBUG)
		if err != nil {
			log.Println("Failed to run vessel:", err)
		}
		return
	}
	invalid()
}

func startDaemon() {
	l := log.Default()
	vm := manager.NewVesselsManager()
	s := api.NewApi(l, vm)
	serv := server.NewServer(l)
	s.RegisterHandlers(serv)
	err := serv.Start()
	if err != nil {
		panic(err)
	}
}

func echo() {
	client, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	msg, err := client.Echo()
	if err != nil {
		panic(err)
	}
	fmt.Println("Echo:", msg)
}

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "daemon":
			startDaemon()
			return
		case "echo":
			echo()
			return
		case "stop":
			if len(args) < 2 {
				invalid()
				return
			}
			id := args[2]
			stopVessel(id)
			return
		case "start":
			if len(args) < 2 {
				invalid()
				return
			}
			id := args[2]
			startVessel(id)
			return
		case "create":
			createVessel()
			return
		case "start_pre_iso_vessel":
			startPreIsolatedVessel(args)
			return
		}
	}
	fmt.Println("bruh")
	// vessel, err := runtime.Create("test", ".", "alpine", "buildfile", "/bin/sh", "/app/entrypoint")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(vessel)
	// fmt.Println("Start:", vessel.Start())
}
