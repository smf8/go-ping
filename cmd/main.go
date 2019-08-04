package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type server struct {
	addr string
	ping float64
}

func main() {
	f, err := os.Open("servers.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := 0
	scanner := bufio.NewScanner(f)
	result := make(chan server)
	for scanner.Scan() {
		text := scanner.Text()
		go func(channel chan server) {
			// ping servers
			cmd, _ := exec.Command("ping", text, "-c 5").Output()
			if strings.Contains(string(cmd), "statistics") {
				s := strings.Split(string(cmd), "\n")
				s = strings.Split(s[len(s)-2], "=")
				s = strings.Split(s[1], "/")
				// s[1] is average
				num, _ := strconv.ParseFloat(s[1], 64)
				channel <- server{text, num}
			} else {
				channel <- server{text, 99999.9999}
			}
			// fmt.Println(string(cmd))
		}(result)
		lines++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	results := make([]server, 0)
	for len(results) != lines {
		select {
		case r := <-result:
			results = append(results, r)
		}
	}
	sort.Slice(results, func(i int, j int) bool {
		return results[i].ping > results[j].ping
	})
	fmt.Println(results[:4])

}
