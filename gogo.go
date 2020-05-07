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
var arrayhandlers map[string]func(commands []string)
var lock *sync.RWMutex
var Logger func(v ...interface{})

func init() {
	Logger = log.Print
}

// Letsgogo = Init
func Letsgogo() {
	handlers = make(map[string]func(command string))
	arrayhandlers = make(map[string]func(commands []string))
	lock = new(sync.RWMutex)
	go func() {
		Logger("[GOGO] Let's go! gogo initiated.\n")
		defer Logger("[ERROR > GOGO] GOGO is down!\n")
		reader := bufio.NewReader(os.Stdin)
	loop:
		for {
			var input string
			var err error
			input, err = reader.ReadString('\n')
			if err != nil {
				Logger(fmt.Sprintf("[GOGO] An error occured when reading input: %s\n", err))
				continue
			}
			strings.TrimRight(input, "\r\n")
			if len(input) == 0 {
				continue
			}
			Logger(fmt.Sprintf("[GOGO] Parsed command: %s", input))
			h, p := parseHeader(input)
			for i := 0; i < len(p); i++ {
				if p[i] == ',' {
					if _, ok := arrayhandlers[h]; ok {
						arrayhandlers[h](strings.Split(p, ","))
					} else {
						Logger(fmt.Sprintf("[GOGO] Unknown array command: %s\n", input))
					}
					continue loop
				}
			}
			if _, ok := handlers[h]; ok {
				handlers[h](p)
			} else {
				Logger(fmt.Sprintf("[GOGO] Unknown command: %s\n", input))

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
	defer lock.Unlock()
	handlers[header] = handler
	Logger(fmt.Sprintf("[GOGO] Registered command: %s\n", header))
}

func RegisterArrayCommand(header string, handler func(command []string)) {
	if arrayhandlers == nil {
		arrayhandlers = make(map[string]func(command []string))
	}
	lock.Lock()
	defer lock.Unlock()
	arrayhandlers[header] = handler
	Logger(fmt.Sprintf("[GOGO] Registered array command: %s\n", header))
}

func ClearRegisteredCommand(header string) {
	if handlers == nil {
		handlers = make(map[string]func(command string))
	}
	lock.Lock()
	defer lock.Unlock()
	delete(handlers, header)
}

func ClearRegisteredArrayCommand(header string) {
	if arrayhandlers == nil {
		arrayhandlers = make(map[string]func(command []string))
	}
	lock.Lock()
	defer lock.Unlock()
	delete(arrayhandlers, header)
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
