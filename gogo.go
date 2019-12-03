package gogo

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

var handlers map[string]func(command string)
var lock *sync.RWMutex

// Letsgogo = Init
func Letsgogo() {
	handlers = make(map[string]func(command string))
	lock = new(sync.RWMutex)
	go func() {
		log.Print("Let's go! gogo initiated.\n")
		reader := bufio.NewReader(os.Stdin)
		input := ""
		for {
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-2]
			h := parseHeader(input)
			if _, ok := handlers[h]; ok {
				handlers[h](input)
			}
		}
	}()
}

// header: "exit" -> "exit:0" -> "0"
// header: "exit" -> "exit" -> "exit"
func RegisterCommand(header string, handler func(command string)) {
	if handlers == nil {
		handlers = make(map[string]func(command string))
	}
	lock.Lock()
	handlers[header] = handler
	lock.Unlock()
}

func parseHeader(command string) (parsed string) {
	i := strings.Index(command, ":")
	if i == -1 {
		parsed = command
	} else {
		parsed = command[i+1:]
	}
	log.Printf("[INFO] Parsed command: %s -> %s", command, parsed)
	return parsed
}
