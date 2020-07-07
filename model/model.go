package model

type Process struct {
	Name      string  `json:"name"`
	MaxCPU    float64 `json:"maxCPU"`
	MaxMem    float64 `json:"maxMem"`
	RestartSh string  `json:"restartSh"`
}

type Stat struct {
	Pid string  `json:"pid"`
	CPU float64 `json:"cpu"`
	MEM float64 `json:"mem"`
}
