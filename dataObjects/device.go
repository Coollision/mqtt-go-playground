package dataObjects

type Availability struct {
	Topic         string `json:"topic"`
	ValueTemplate string `json:"value_template"`
}

type Device struct {
	ConfigurationURL string   `json:"configuration_url,omitempty"`
	HwVersion        int      `json:"hw_version,omitempty"`
	Identifiers      []string `json:"identifiers,omitempty"`
	Manufacturer     string   `json:"manufacturer,omitempty"`
	Model            string   `json:"model,omitempty"`
	ModelID          string   `json:"model_id,omitempty"`
	Name             string   `json:"name,omitempty"`
	SwVersion        string   `json:"sw_version,omitempty"`
	ViaDevice        string   `json:"via_device,omitempty"`
}

type Origin struct {
	Name string `json:"name,omitempty"`
	Sw   string `json:"sw,omitempty"`
	URL  string `json:"url,omitempty"`
}
