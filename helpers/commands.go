package helpers

// Commands will be the byte arrays holding the hex commands for the readers
var commands = map[string][]byte{
	"InitiateCommand": {0xFF, 0xCA, 0x00, 0x00, 0x00}, // this is what needs to be sent to initiate any command
	"GET_UID":         {0xFF, 0xCA, 0x00, 0x00, 0x04}, // send this to get the serial nnumber of the connected PICC
}
