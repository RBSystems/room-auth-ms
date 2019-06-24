package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
)

var token = "f91a62eaa4c48723412fb1925d60ef"

//recieves the UTA Number from the card and returns the byuID Number
func GetIdNumber(utanumber string) (string, *nerr.E) {
	x := utanumber
	x = x[:14]
	fmt.Println(x)
	var a [7]string
	s := 0
	t := 2
	for i := 0; i <= 6; i++ {
		a[i] = x[s:t]
		s = s + 2
		t = t + 2
	}
	test := ""
	for i := 6; i >= 0; i-- {
		test += a[i]
	}
	test = strings.ToUpper(test)

	svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"uta_id": {
				S: aws.String(test),
			},
		},
		TableName: aws.String("byu_uta_id"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", nil
	}
	if result.Item != nil {
		if result.Item["byu_id"].S != nil {
			return fmt.Sprintf("%s", *result.Item["byu_id"].S), nil
		}
		return "", nil
	}
	return "", nil

}

//takes the BYU ID number and uses the Person API to return their info
func GetNetID(idnumber string) (string, *nerr.E) {
	//this is where we get the NetID
	weburl := fmt.Sprintf("https://api.byu.edu:443/byuapi/persons/v3/%s?field_sets=basic", idnumber)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	var output structs.Person

	input := ""

	outputJSON, _, err := jsonhttp.CreateAndExecuteJSONRequest("getNetID", "GET", weburl,
		input, headers, 20, &output)
	log.L.Debug(outputJSON)
	if err != nil {
		log.L.Errorf("Failed to get NetID: %s", err.Error())
		return "", nil
	}
	NetID := output.ID
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
	}

	runner.SendEvent(Event)

}
