package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	coreURL := flag.String("core", "http://127.0.0.1:8090", "control API base URL")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("auxitalkctl 0.1.0-dev")
		fmt.Println("commands: status, plugins, actions")
		return
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "status":
		printStatus(*coreURL)
	case "plugins":
		printPlugins(*coreURL)
	case "actions":
		printActions(*coreURL)
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func printStatus(base string) {
	data := fetch(base + "/api/status")
	fmt.Printf("%s\n", pretty(data))
}

func printPlugins(base string) {
	data := fetch(base + "/api/plugins")
	fmt.Printf("%s\n", pretty(data))
}

func printActions(base string) {
	data := fetch(base + "/api/actions")
	fmt.Printf("%s\n", pretty(data))
}

func fetch(url string) map[string]any {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]any
	_ = json.Unmarshal(body, &result)
	return result
}

func pretty(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
