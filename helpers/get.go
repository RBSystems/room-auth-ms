package helpers

import (
	"os"
	"strings"
	"time"

	"github.com/byuoitav/wso2services/wso2requests"

	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/room-auth-ms/structs"
)

//takes the Card Serial Number and uses the Person API to return their info
func GetNetID(cardNumber string) (string, *nerr.E) {
	//this is where we get the NetID

	var output structs.WSO2CredentialResponse

	err := wso2requests.MakeWSO2Request("GET", "https://api.byu.edu:443/byuapi/persons/v3/?credentials.credential_type=CARD_SERIAL_NUMBER&credentials.credential_id="+cardNumber, "", &output)
	if err != nil {
		log.L.Debugf("Error when making WSO2 request %v", err)
	}
	//log.L.Debugf("this is the output %v", output)
	NetID := output.Values[0].Basic.NetID.Value
	return NetID, nil
}

//fix send event thing
func SendEvent(netid string, runner messenger.Messenger) {

	room := os.Getenv("SYSTEM_ID")
	a := strings.Split(room, "-")
	roominfo := events.BasicRoomInfo{}
	if len(a) == 3 {
		roominfo = events.BasicRoomInfo{
			BuildingID: a[0],
			RoomID:     a[0] + "-" + a[1],
		}
	} else {
		roominfo = events.BasicRoomInfo{
			BuildingID: room,
			RoomID:     room,
		}
	}

	basicdevice := events.BasicDeviceInfo{
		BasicRoomInfo: roominfo,
		DeviceID:      os.Getenv("SYSTEM_ID"),
	}

	Event := events.Event{
		GeneratingSystem: os.Getenv("SYSTEM_ID"),
		Timestamp:        time.Now(),
		Key:              "Login",
		Value:            "True",
		User:             netid,
		TargetDevice:     basicdevice,
		AffectedRoom:     roominfo,
		EventTags: []string{
			events.Heartbeat,
		},
	}

	runner.SendEvent(Event)

}
