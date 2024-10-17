package dataObjects

// Generated by https://quicktype.io
type DeviceConfigCover struct {
	Availability        []Availability `json:"availability,omitempty"`
	AvailabilityMode    string         `json:"availability_mode,omitempty"`
	CommandTopic        string         `json:"command_topic,omitempty"`
	Device              Device         `json:"device,omitempty"`
	DeviceClass         string         `json:"device_class,omitempty"`
	Name                interface{}    `json:"name,omitempty"`
	ObjectID            string         `json:"object_id,omitempty"`
	Origin              Origin         `json:"origin,omitempty"`
	PositionTemplate    string         `json:"position_template,omitempty"`
	PositionTopic       string         `json:"position_topic,omitempty"`
	SetPositionTemplate string         `json:"set_position_template,omitempty"`
	SetPositionTopic    string         `json:"set_position_topic,omitempty"`
	StateClosed         string         `json:"state_closed,omitempty"`
	StateClosing        string         `json:"state_closing,omitempty"`
	StateOpen           string         `json:"state_open,omitempty"`
	StateOpening        string         `json:"state_opening,omitempty"`
	StateStopped        string         `json:"state_stopped,omitempty"`
	PayloadOpen         string         `json:"payload_open,omitempty"`
	PayloadClose        string         `json:"payload_close,omitempty"`
	PayloadStop         string         `json:"payload_stop,omitempty"`
	StateTopic          string         `json:"state_topic,omitempty"`
	UniqueID            string         `json:"unique_id,omitempty"`
	ValueTemplate       string         `json:"value_template,omitempty"`
}
