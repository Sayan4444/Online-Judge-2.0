package schema

type JudgeResponse struct {
	Stderr           string   `json:"stderr"`
	Time             string   `json:"time"`
	Memory           string   `json:"memory"`
	ExitSignal       string   `json:"exit_signal"`
	ExitCode         string   `json:"exit_code"`
	Message          string   `json:"message"`
	Result           string   `json:"result"`
	CompileOutput    string   `json:"compile_output"`
	WrongAnswers     []WrongAnswer `json:"wrong_answers"`
}

type WrongAnswer struct {
	TestCaseID string `json:"test_case_id"`
	Stdout     string `json:"stdout"`
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
