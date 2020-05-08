package gogo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var handlers map[string]func(params []string)
var lock *sync.RWMutex
var Logger func(v ...interface{})

func init() {
	Logger = log.Print
}

// Letsgogo = Init
func Letsgogo() {
	handlers = make(map[string]func(params []string))
	lock = new(sync.RWMutex)
	go func() {
		Logger("[GOGO] Let's go! gogo initiated.\n")
		defer Logger("[ERROR > GOGO] GOGO is down!\n")
		reader := bufio.NewReader(os.Stdin)
		for {
			var input string
			var err error
			input, err = reader.ReadString('\n')
			if err != nil {
				Logger(fmt.Sprintf("[GOGO] An error occured when reading input: %s\n", err))
				continue
			}
			input = strings.TrimRight(input, "\r\n")
			if len(input) == 0 {
				continue
			}
			Trigger(input)
		}
	}()
}

func Trigger(command string) {
	commands := strings.Split(command, " ")
	lock.RLock()
	defer lock.RUnlock()
	if _, ok := handlers[commands[0]]; ok {
		handlers[commands[0]](commands[1:])
	} else {
		Logger(fmt.Sprintf("[GOGO] Unknown command header: %s\n", command))
	}
}

// header: "exit" -> "exit:0" -> "0"
// header: "exit" -> "exit" -> "exit"
func RegisterCommand(header string, handler func(params []string)) {
	if handlers == nil {
		handlers = make(map[string]func(params []string))
	}
	lock.Lock()
	defer lock.Unlock()
	handlers[header] = handler
	Logger(fmt.Sprintf("[GOGO] Registered command: %s\n", header))
}

func UnregisterCommand(command string) {
	if handlers == nil {
		handlers = make(map[string]func(params []string))
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := handlers[command]; ok {
		delete(handlers, command)
	} else {
		Logger(fmt.Sprintf("[GOGO] Not registered: %s", command))
	}
}

func ShowRegisteredCommands() {
	fmt.Print("[GOGO] GOGO registered commands:\n")
	for k, v := range handlers {
		if v == nil {
			continue
		}
		fmt.Printf(" > %s\n", k)
	}
}
