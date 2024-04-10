package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		executable, err := os.Executable()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		log.Fatalf(":: Usage: go run %s.go <path>", filepath.Base(executable))
		return
	}

	minecraftDir := os.Args[1]

	servers, err := scanMinecraftServers(minecraftDir)
	if err != nil {
		log.Fatalf("Error scanning Minecraft servers: %v", err)
		return
	}

	fmt.Println("Scanning servers:")
	for server, worlds := range servers {
		fmt.Printf("Server %s, Worlds: %d\n", server, len(worlds))
	}

	start := time.Now()
	fmt.Printf("Run time: %v", time.Since(start))
}

func scanMinecraftServers(dir string) (map[string][]string, error) {
	servers := make(map[string][]string)

	var wg sync.WaitGroup

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error acure while readine the dir: %s :: %v\n", dir, err)
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go func(serverName string) {
				defer wg.Done()
				if isMinecraftServer(filepath.Join(dir, serverName)) {
					worlds, err := countWorlds(filepath.Join(dir, serverName))
					if err != nil {
						log.Fatalf("Error counting worlds for server %s :: %v", serverName, err)
						return
					}

					servers[serverName] = worlds
				}
			}(file.Name())
		}
	}

	wg.Wait()

	return servers, nil
}

func isMinecraftServer(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "server.properties"))
	return err == nil
}

func countWorlds(dir string) ([]string, error) {
	var worlds []string

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error acure while readine the dir: %s :: %v\n", dir, err)
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			_, err := os.Stat(filepath.Join(dir, file.Name(), "uid.dat"))
			if err == nil {
				worlds = append(worlds, file.Name())
			}
		}
	}

	return worlds, nil
}
