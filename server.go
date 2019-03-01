package main

import (
	"fmt"
	"os"

	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/room-auth-ms/helpers"
	"github.com/ebfe/scard"
)

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func waitUntilCardRelease(ctx *scard.Context, readers []string, index int) error {
	rs := make([]scard.ReaderState, 1)

	rs[0].Reader = readers[index]
	rs[0].CurrentState = scard.StatePresent

	for {

		if rs[0].EventState&scard.StateEmpty != 0 {
			return nil
		}
		rs[0].CurrentState = rs[0].EventState

		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return err
		}
	}
}

func main() {

	// Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		errorExit(err)
	}
	defer ctx.Release()

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		errorExit(err)
	}

	// connect to the hub
	messenger, er := messenger.BuildMessenger("localhost:7100", base.Messenger, 5000)
	if er != nil {
		log.L.Fatalf("failed to build messenger: %s", er)
	}

	fmt.Printf("Found %d readers:\n", len(readers))
	for i, reader := range readers {
		fmt.Printf("[%d] %s\n", i, reader)
	}

	if len(readers) > 0 {
		for {
			fmt.Println("Waiting for a Card")
			index, err := waitUntilCardPresent(ctx, readers)
			if err != nil {
				errorExit(err)
			}

			// Connect to card
			fmt.Println("Connecting to card in ", readers[index])
			card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
			if err != nil {
				errorExit(err)
			}
			defer card.Disconnect(scard.ResetCard)

			var cmd = []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}

			rsp, err := card.Transmit(cmd)
			if err != nil {
				errorExit(err)
			}
			uid := string(rsp[0:7])
			uidS := fmt.Sprintf("%x", uid)
			fmt.Printf("Tag UID is: %s\n", uidS)
			idNumber, er := helpers.GetIdNumber(uidS)
			if er != nil {
				log.L.Errorf("Failed to get ID Number: %s", er.Error())
				fmt.Printf("the error is: %s", er)
			}
			fmt.Printf("idnumber: %s", idNumber)
			if idNumber != "" {
				NetID, er := helpers.GetNetID(idNumber)
				if er != nil {
					log.L.Errorf("Failed to get NetID: %s", er.Error())
					fmt.Printf("the error is: %s", er)
				}
				fmt.Printf("Tag UID is: %s\n", NetID)
				helpers.SendEvent(NetID, *messenger)
				fmt.Printf("Writting as keyboard input...")
				fmt.Printf("Done.\n")
			}

			card.Disconnect(scard.ResetCard)

			//Wait while card will be released
			fmt.Print("Waiting for card release...")
			err = waitUntilCardRelease(ctx, readers, index)
			fmt.Println("Card released.")

		}

	}

}
