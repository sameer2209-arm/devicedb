package main

import (
    "fmt"
    
    . "devicedb/server"
)

func init() {
    registerCommand("start", startServer, startUsage)
}

var startUsage string = 
`Usage: devicedb start -conf=[config file]
`

func startServer() {
    var sc ServerConfig
        
    err := sc.LoadFromFile(*optConfigFile)

    if err != nil {
        fmt.Printf("Unable to load config file: %s\n", err.Error())
        
        return
    }

    server, err := NewServer(sc)

    if err != nil {
        fmt.Printf("Unable to create server: %s\n", err.Error())
        
        return
    }

    sc.Hub.SyncController().Start()
    sc.Hub.StartForwardingEvents()
    sc.Hub.StartForwardingAlerts()
    server.StartGC()

    server.Start()
}