package handlers

import (
	"net/http"

	"github.com/byuoitav/Id-cards-microservice/helpers"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

/*
 * Function:  readTag
 * --------------------
 * reads a tag from the NFC reader
 *
 *  returns: an error if the the serial number, or ATS of the connected PICC (Proximity Integrated Circuit Card) can't be read
 *  Hopefully this will have some kind of PID (Unique ID)
 */

// func readTag() error {

// 	return http.StatusOK
// }

func GetUserInfo(context echo.Context) error {
	utanumber := context.Param("utanumber")

	idNumber, err := helpers.GetIdNumber(utanumber)
	if err != nil {
		log.L.Errorf("Failed to get ID Number: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	log.L.Debugf("The idnumber returned is %s", idNumber)

	NetID, err := helpers.GetNetID(idNumber)
	if err != nil {
		log.L.Errorf("Failed to get NetID: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, NetID)
}
