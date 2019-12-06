package gogo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var handlers map[string]func(command string)
var arrayhandlers map[string]func(command []string)
var lock *sync.RWMutex

// Letsgogo = Init
func Letsgogo() {
	handlers = make(map[string]func(command string))
	lock = new(sync.RWMutex)
	go func() {
		log.Print("[GOGO] Let's go! gogo initiated.\n")
		reader := bufio.NewReader(os.Stdin)
		input := ""
		for {
			input, _ = reader.ReadString('\n')
			input = input[:len(input)-2]
			h, p := parseHeader(input)
			if strings.ContainsAny(p, ",") {
				if _, ok := arrayhandlers[h]; ok {
					arrayhandlers[h](strings.Split(p, ","))
				} else {
					log.Printf("[GOGO] Unknown command: %s.\n", input)
				}
			} else {
				if _, ok := handlers[h]; ok {
					handlers[h](p)
				} else {
					log.Printf("[GOGO] Unknown command: %s.\n", input)

				}
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

func RegisterArrayCommand(header string, handler func(command []string)) {
	if arrayhandlers == nil {
		arrayhandlers = make(map[string]func(command []string))
	}
	lock.Lock()
	arrayhandlers[header] = handler
	lock.Unlock()
}

func parseHeader(command string) (header string, parsed string) {
	i := strings.Index(command, ":")
	if i == -1 {
		header = command
		parsed = command
	} else {
		header = command[:i]
		parsed = command[i+1:]
	}
	log.Printf("[INFO] Parsed command: %s -> %s", command, parsed)
	return header, parsed
}

func ShowRegisteredCommands() {
	fmt.Print("[GOGO] GOGO registered commands:\n")
	for k, v := range handlers {
		fmt.Printf("> %s: %v\n", k, v == nil)
	}
}