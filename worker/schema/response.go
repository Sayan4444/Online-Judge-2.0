package schema

type JudgeResponse struct {
	Stdin         string `json:"stdin"`
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	Time          string `json:"time"`
	Memory        string `json:"memory"`
	ExitSignal    string `json:"exit_signal"`
	ExitCode      string `json:"exit_code"`
	Message       string `json:"message"`
	Result        string `json:"result"`
	CompileOutput string `json:"compile_output"`
}

const (
	ResultAccepted                 = "AC"
	ResultWrongAnswer              = "WA"
	ResultTimeLimitExceeded        = "TLE"
	ResultMemoryLimitExceeded      = "MLE"
	ResultRuntimeError             = "RE"
	ResultCompileError             = "CE"
	ResultCompileTimeLimitExceeded = "CTLE"
	ResultOutputLimitExceeded      = "OLE"
	ResultSystemError              = "SE"
	ResultUnknownError             = "UE"
)
