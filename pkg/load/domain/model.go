package domain

type AgentInfo struct {
	AgentId string `json:"agentId" form:"agentId,omitempty"`

	Hostname string `json:"hostname" form:"hostname,omitempty"`
	Username string `json:"username" form:"username,omitempty"`
	TcpPort  string `json:"tcpPort" form:"tcpPort,omitempty"`
	Shutdown bool   `json:"shutdown" form:"shutdown,omitempty"`
}

func NewAgentInfo() AgentInfo {
	return AgentInfo{}
}

type LoadTest struct {
	TestId string `json:"testId" form:"TestId,omitempty"`

	Protocol string `json:"protocol" form:"protocol,omitempty"`
	Hostname string `json:"hostname" form:"hostname,omitempty"`
	Port     string `json:"port" form:"port,omitempty"`
	Path     string `json:"path" form:"port,omitempty"`
	BodyData string `json:"bodyData" form:"bodyData,omitempty"`

	Threads   string `json:"threads" form:"threads,omitempty"`
	RampTime  string `json:"rampTime" form:"rampTime,omitempty"`
	LoopCount string `json:"loopCount" form:"loopCount,omitempty"`
}

type LoadTestStatistics struct {
	Label         string  `json:"label"`
	RequestCount  int     `json:"requestCount"`
	Average       float64 `json:"average"`
	Median        float64 `json:"median"`
	NinetyPercent float64 `json:"ninetyPercent"`
	NinetyFive    float64 `json:"ninetyFive"`
	NinetyNine    float64 `json:"ninetyNine"`
	MinTime       float64 `json:"minTime"`
	MaxTime       float64 `json:"maxTime"`
	ErrorPercent  float64 `json:"errorPercent"`
	Throughput    float64 `json:"throughput"`
	ReceivedKB    float64 `json:"receivedKB"`
	SentKB        float64 `json:"sentKB"`
}