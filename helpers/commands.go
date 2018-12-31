package helpers

var commands = map[string][]byte{
	"InitiateCommand": {0xFF, 0xCA, 0x00, 0x00, 0x00}, // This command is what needs to be sent to initiate any command
}
