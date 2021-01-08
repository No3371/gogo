# gogo
Console command interpreter for go

# Usage
`gogo.Letsgogo()` to initialize the interpreter, it will start taking input from stdin.

`gogo.RegisterCommand(header string, handler func(params []string))` to add command handler.
For example, `gogo histroy 5`, `gogo` will be used as header (key to handler), all following string splitted by whitespace will be pass to handlers ([`histroy`, `5`]). 

You can also un-register a command with `gogo.UnregisterCommand(header string)`.
