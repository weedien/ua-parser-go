package uaparser

type IBrand struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type IBrowser struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Major   string `json:"major,omitempty"`
}

type ICpu struct {
	Architecture string `json:"architecture,omitempty"`
}

type IDevice struct {
	Type   string `json:"type,omitempty"` // Mobile, Desktop, Bot, Console
	Model  string `json:"model,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type IEngine struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type IOs struct {
	Platform string `json:"platform,omitempty"`
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
}

type IResult struct {
	UA      string   `json:"ua"`
	Browser IBrowser `json:"browser"`
	Cpu     ICpu     `json:"cpu"`
	Device  IDevice  `json:"device"`
	Engine  IEngine  `json:"engine"`
	Os      IOs      `json:"os"`
}
