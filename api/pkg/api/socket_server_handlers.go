package api

import (
	"encoding/json"
	"fmt"
	"github.com/antoniodipinto/ikisocket"
	"mizuserver/pkg/tap"
)

var browserClientSocketUUIDs = make([]string, 0)
var SocketHarOutChannel = make(chan *tap.OutputChannelItem, 1000)

func WebSocketConnect(ep *ikisocket.EventPayload) {
	if ep.Kws.GetAttribute("is_tapper") == true {
		fmt.Println(fmt.Sprintf("Websocket Connection event - Tapper connected: %s", ep.SocketUUID))
	} else {
		fmt.Println(fmt.Sprintf("Websocket Connection event - Browser socket connected: %s", ep.SocketUUID))
		browserClientSocketUUIDs = append(browserClientSocketUUIDs, ep.SocketUUID)
	}
}

func WebSocketDisconnect(ep *ikisocket.EventPayload) {
	if ep.Kws.GetAttribute("is_tapper") == true {
		fmt.Println(fmt.Sprintf("Disconnection event - Tapper connected: %s", ep.SocketUUID))
	} else {
		fmt.Println(fmt.Sprintf("Disconnection event - Browser socket connected: %s", ep.SocketUUID))
		removeSocketUUIDFromBrowserSlice(ep.SocketUUID)
	}
}

func broadcastToBrowserClients(message []byte) {
	ikisocket.EmitToList(browserClientSocketUUIDs, message)
}

func WebSocketClose(ep *ikisocket.EventPayload) {
	if ep.Kws.GetAttribute("is_tapper") == true {
		fmt.Println(fmt.Sprintf("Websocket Close event - Tapper connected: %s", ep.SocketUUID))
	} else {
		fmt.Println(fmt.Sprintf("Websocket  Close event - Browser socket connected: %s", ep.SocketUUID))
		removeSocketUUIDFromBrowserSlice(ep.SocketUUID)
	}
}

func WebSocketError(ep *ikisocket.EventPayload) {
	fmt.Println(fmt.Sprintf("Socket error - Socket uuid : %s", ep.Kws.GetStringAttribute("user_id")))
}

func WebSocketMessage(ep *ikisocket.EventPayload) {
	if ep.Kws.GetAttribute("is_tapper") == true {
		var tapOutput tap.OutputChannelItem
		err := json.Unmarshal(ep.Data, &tapOutput)
		if err != nil {
			fmt.Printf("Could not unmarshal message received from tapper websocket %v", err)
		} else {
			SocketHarOutChannel <- &tapOutput
		}
	} else {
		fmt.Println("Received Web socket message from non tapper websocket, no handler is defined for this message")
	}
}

func removeSocketUUIDFromBrowserSlice(uuidToRemove string) {
	newUUIDSlice := make([]string, 0)
	for _, uuid := range browserClientSocketUUIDs {
		if uuid != uuidToRemove {
			newUUIDSlice = append(newUUIDSlice, uuid)
		}
	}
	browserClientSocketUUIDs = newUUIDSlice
}
