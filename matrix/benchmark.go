package matrix

import (
	"encoding/json"
	"fmt"
	"time"
	"math/rand"
)

type SyncResponse struct {
	NextBatch   string            `json:"next_batch"`
	Rooms       Rooms             `json:"rooms"`
	DeviceLists DeviceLists       `json:"device_lists"`
	ToDevice    ToDevice          `json:"to_device"`
	Presence    Presence          `json:"presence"`
	AccountData UserAccountData   `json:"account_data"`
}

type Rooms struct {
	Join map[string]Room `json:"join"`
}

type Room struct {
	Summary     RoomSummary     `json:"summary"`
	State       State           `json:"state"`
	Timeline    Timeline        `json:"timeline"`
	Ephemeral   Ephemeral       `json:"ephemeral"`
	AccountData RoomAccountData `json:"account_data"`
}

type RoomSummary struct {
	Heroes              []string `json:"m.heroes"`
	JoinedMemberCount   int      `json:"m.joined_member_count"`
	InvitedMemberCount  int      `json:"m.invited_member_count"`
}

type State struct {
	Events []Event `json:"events"`
}

type Timeline struct {
	Events    []Event `json:"events"`
	PrevBatch string  `json:"prev_batch"`
	Limited   bool    `json:"limited"`
}

type Ephemeral struct {
	Events []Event `json:"events"`
}

type RoomAccountData struct {
	Events []Event `json:"events"`
}

type DeviceLists struct {
	Changed []string `json:"changed"`
	Left    []string `json:"left"`
}

type ToDevice struct {
	Events []Event `json:"events"`
}

type Presence struct {
	Events []Event `json:"events"`
}

type UserAccountData struct {
	Events []Event `json:"events"`
}

type Event struct {
	Type             string                 `json:"type"`
	StateKey         string                 `json:"state_key,omitempty"`
	Content          map[string]interface{} `json:"content"`
	Sender           string                 `json:"sender"`
	OriginServerTS   int64                  `json:"origin_server_ts"`
	Unsigned         map[string]interface{} `json:"unsigned,omitempty"`
	EventID          string                 `json:"event_id"`
	RoomID           string                 `json:"room_id"`
}


var (
	names    = []string{"Alice", "Bob", "Charlie", "Dave", "Eve", "Frank"}
	messages = []string{
		"Hello, world!",
		"How are you?",
		"Good morning!",
		"What's up?",
		"See you later!",
		"Goodbye!",
		"Have a great day!",
		"Nice to meet you!",
		"Long time no see!",
		"Take care!",
	}
)

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func generateSyncResponse(numRooms, numEventsPerRoom, numUsersPerRoom int) SyncResponse {
	response := SyncResponse{
		NextBatch:   "s72595_4483_1934",
		Rooms:       Rooms{Join: make(map[string]Room)},
		DeviceLists: DeviceLists{Changed: []string{}, Left: []string{}},
		ToDevice:    ToDevice{Events: []Event{}},
		Presence:    Presence{Events: []Event{}},
		AccountData: UserAccountData{Events: []Event{}},
	}

	for i := 0; i < numRooms; i++ {
		roomID := fmt.Sprintf("!room%d:example.com", i)
		room := Room{
			Summary: RoomSummary{
				Heroes:              generateHeroes(numUsersPerRoom),
				JoinedMemberCount:   numUsersPerRoom,
				InvitedMemberCount:  0,
			},
			State:       State{Events: []Event{}},
			Timeline:    Timeline{Events: generateTimelineEvents(roomID, numEventsPerRoom), PrevBatch: "s72584_4483_1933", Limited: false},
			Ephemeral:   Ephemeral{Events: []Event{generateTypingEvent(roomID)}},
			AccountData: RoomAccountData{Events: []Event{generateFullyReadEvent(roomID)}},
		}
		response.Rooms.Join[roomID] = room
	}

	return response
}

func generateHeroes(numUsers int) []string {
	heroes := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		heroes[i] = fmt.Sprintf("@%s:example.com", randomChoice(names))
	}
	return heroes
}

func generateStateEvent(roomID, roomName string) Event {
	return Event{
		Type:   "m.room.name",
		StateKey: "",
		Content: map[string]interface{}{
			"name": roomName,
		},
		Sender:         fmt.Sprintf("@%s:example.com", randomChoice(names)),
		OriginServerTS: time.Now().UnixNano() / int64(time.Millisecond),
		EventID:        generateEventID(),
		RoomID:         roomID,
	}
}

func generateTimelineEvents(roomID string, numEvents int) []Event {
	events := make([]Event, numEvents)
	for i := 0; i < numEvents; i++ {
		events[i] = Event{
			Type:   "m.room.message",
			Content: map[string]interface{}{
				"msgtype": "m.text",
				"body":    randomChoice(messages),
			},
			Sender:         fmt.Sprintf("@%s:example.com", randomChoice(names)),
			OriginServerTS: time.Now().UnixNano() / int64(time.Millisecond),
			EventID:        generateEventID(),
			RoomID:         roomID,
		}
	}
	return events
}

func generateTypingEvent(roomID string) Event {
	return Event{
		Type: "m.typing",
		Content: map[string]interface{}{
			"user_ids": []string{fmt.Sprintf("@%s:example.com", randomChoice(names))},
		},
		RoomID: roomID,
	}
}

func generateFullyReadEvent(roomID string) Event {
	return Event{
		Type: "m.fully_read",
		Content: map[string]interface{}{
			"event_id": generateEventID(),
		},
		RoomID: roomID,
	}
}

func generateEventID() string {
	return fmt.Sprintf("$%d_%s", time.Now().UnixNano(), "example")
}

func GenerateFullSyncResp(numRooms int, numEventsPerRoom int, numUsersPerRoom int) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())

	syncResponse := generateSyncResponse(numRooms, numEventsPerRoom, numUsersPerRoom)

	return json.MarshalIndent(syncResponse, "", "  ")
}