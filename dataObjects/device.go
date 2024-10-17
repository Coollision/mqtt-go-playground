package dataObjects

type Availability struct {
	Topic string `json:"topic"`
}

type Device struct {
	ConfigurationURL string   `json:"configuration_url,omitempty"`
	Identifiers      []string `json:"identifiers,omitempty"`
	Manufacturer     string   `json:"manufacturer,omitempty"`
	Model            string   `json:"model,omitempty"`
	Name             string   `json:"name,omitempty"`
	SwVersion        string   `json:"sw_version,omitempty"`
	ViaDevice        string   `json:"via_device,omitempty"`
}

type Origin struct {
	Name string `json:"name,omitempty"`
	Sw   string `json:"sw,omitempty"`
	URL  string `json:"url,omitempty"`
}
