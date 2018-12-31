package handlers

import "net/http"

/*
 * Function:  readTag
 * --------------------
 * reads a tag from the NFC reader
 *
 *  returns: an error if the the serial number, or ATS of the connected PICC (Proximity Integrated Circuit Card) can't be read
 *  Hopefully this will have some kind of PID (Unique ID)
 */

func readTag() error {

	return http.StatusOK
}
