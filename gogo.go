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
		defer log.Print("[ERROR > GOGO] GOGO is down!\n")
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("[GOGO] An error occured when reading input: %s\n", err)
				continue
			}
			input = input[:len(input)-2]
			h, p := parseHeader(input)
			if strings.ContainsAny(p, ",") {
				if _, ok := arrayhandlers[h]; ok {
					arrayhandlers[h](strings.Split(p, ","))
				} else {
					log.Printf("[GOGO] Unknown command: %s\n", input)
				}
			} else {
				if _, ok := handlers[h]; ok {
					handlers[h](p)
				} else {
					log.Printf("[GOGO] Unknown command: %s\n", input)

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
	log.Printf("[GOGO] Parsed command:\n%s -> %s\n", command, parsed)
	return header, parsed
}

func ShowRegisteredCommands() {
	fmt.Print("[GOGO] GOGO registered commands:\n")
	fmt.Print("These handlers expect single arguments:\n")
	for k, v := range handlers {
		if v == nil {
			continue
		}
		fmt.Printf(" > %s\n", k)
	}
	fmt.Print("These handlers expect array arguments:\n")
	for k, v := range arrayhandlers {
		if v == nil {
			continue
		}
		fmt.Printf(" > %s\n", k)
	}
}
