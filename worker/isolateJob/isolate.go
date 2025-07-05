package isolatejob

import (
	"OJ-Worker/schema"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	StdinFileName    = "stdin.txt"
	StdoutFileName   = "stdout.txt"
	StderrFileName   = "stderr.txt"
	MetadataFileName = "metadata.txt"
)

var boxIDCounter int64

type IsolateJob struct {
	Submission *schema.RabbitMQPayload
	Response   *schema.JudgeResponse
	BoxID      int
	WorkDir    string
	BoxDir     string
	TmpDir     string
	SourceFile string
	InputFile  string
	OutputFile string
	ErrorFile  string
	MetaFile   string
}

func ProcessSubmission(submission *schema.RabbitMQPayload, ctx context.Context) error {

	job := &IsolateJob{
		Submission: submission,
		BoxID:      int(atomic.AddInt64(&boxIDCounter, 1)) % 2147483647,
		Response:   &schema.JudgeResponse{},
	}

	return job.Execute(ctx)
}

func (j *IsolateJob) Execute(ctx context.Context) error {
	if err := j.InitializeIsolate(ctx); err != nil {
		j.Response.Result = schema.ResultSystemError
		j.CleanUp(ctx)
		return fmt.Errorf("failed to initialize isolate: %v", err)
	}
	success, err := j.Compile(ctx)
	if err != nil {
		j.Response.Result = schema.ResultSystemError
		j.CleanUp(ctx)
		return fmt.Errorf("failed to compile: %v", err)
	}
	if !success {
		j.CleanUp(ctx)
		return nil
	}

	success, err = j.Run(ctx)
	if err != nil {
		j.Response.Result = schema.ResultSystemError
		j.CleanUp(ctx)
		return fmt.Errorf("failed to run: %v", err)
	}
	if !success {
		j.CleanUp(ctx)
		return nil
	}

	j.CleanUp(ctx)
	return nil

}

func (j *IsolateJob) InitializeIsolate(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "isolate",
		"-b", strconv.Itoa(j.BoxID),
		"--init",
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to initialize isolate box: %v", err)
	}

	j.WorkDir = strings.TrimSpace(string(output))

	j.BoxDir = filepath.Join(j.WorkDir, "box")
	j.TmpDir = filepath.Join(j.WorkDir, "tmp")

	j.SourceFile = filepath.Join(j.BoxDir, j.Submission.SourceFileName)
	j.InputFile = filepath.Join(j.WorkDir, StdinFileName)
	j.OutputFile = filepath.Join(j.WorkDir, StdoutFileName)
	j.ErrorFile = filepath.Join(j.WorkDir, StderrFileName)
	j.MetaFile = filepath.Join(j.WorkDir, MetadataFileName)

	files := []string{j.SourceFile, j.InputFile, j.OutputFile, j.ErrorFile, j.MetaFile}
	for _, file := range files {
		if err := j.InitializeFiles(file, ctx); err != nil {
			return fmt.Errorf("failed to initialize file %s: %v", file, err)
		}
	}

	if err := os.WriteFile(j.SourceFile, []byte(j.Submission.SourceCode), 0644); err != nil {
		return fmt.Errorf("failed to write source code to file %s: %v", j.SourceFile, err)
	}

	if err := os.WriteFile(j.InputFile, []byte(j.Submission.StdIn), 0644); err != nil {
		return fmt.Errorf("failed to write stdin to file %s: %v", j.InputFile, err)
	}

	return nil
}

func (j *IsolateJob) InitializeFiles(filename string, ctx context.Context) error {

	user := os.Getenv("USER")

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", fmt.Sprintf("sudo touch %s && sudo chown %s: %s", filename, user, filename))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize file %s: %v", filename, err)
	}

	return nil
}

func (j *IsolateJob) Compile(ctx context.Context) (bool, error) {
	if j.Submission.CompileCmd == "" {
		return true, nil
	}
	compileScript := filepath.Join(j.BoxDir, "compile.sh")
	compileOutput := filepath.Join(j.WorkDir, "compile_output.txt")
	j.InitializeFiles(compileOutput, ctx)

	if err := os.WriteFile(compileScript, []byte(j.Submission.CompileCmd), 0755); err != nil {
		return false, fmt.Errorf("failed to write compile script to file %s: %v", compileScript, err)
	}

	cmdRun := `isolate \
	-s \
	-b %d \
	-M %s \
	--stderr-to-stdout \
	-i /dev/null \
	-t %d \
	-w %d \
	-x 0 \
	-m %d \
	-k %d \
	-p4 \
	-f %d \
	-E "HOME=/tmp" \
	-E "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
	-d "/etc:noexec" \
	--run \
	-- /bin/bash %s > %s`

	actualCompileCmd := fmt.Sprintf(cmdRun, j.BoxID, j.MetaFile, j.Submission.TimeLimit, j.Submission.WallTimeLimit, j.Submission.MemoryLimit, j.Submission.StackLimit, j.Submission.OutputLimit, filepath.Base(compileScript), compileOutput)

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", actualCompileCmd)
	err := cmd.Run()
	compileOutputText, readErr := os.ReadFile(compileOutput)
	if readErr == nil && len(compileOutputText) > 0 {
		j.Response.CompileOutput = string(compileOutputText)
	}

	metadata, _ := j.getMetadata()

	fmt.Println("----------------Compile Metadata------------")
	fmt.Println(metadata)

	filesToRemove := []string{compileScript, compileOutput}

	for _, file := range filesToRemove {
		if err := exec.CommandContext(ctx, "sudo", "rm", "-rf", file).Run(); err != nil {
			return false, fmt.Errorf("failed to remove file %s: %v", file, err)
		}
	}

	j.resetMetadata(ctx)

	if _, ok := err.(*exec.ExitError); ok {
		if status, ok := metadata["status"]; ok {
			j.Response.Message = "Compile Error"
			if status == "TO" {
				j.Response.Result = schema.ResultCompileTimeLimitExceeded
			} else {
				j.Response.Result = schema.ResultCompileError
			}
		}
		return false, nil

	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (j *IsolateJob) Run(ctx context.Context) (bool, error) {
	runScript := filepath.Join(j.BoxDir, "run.sh")

	if err := os.WriteFile(runScript, []byte(j.Submission.RunCmd), 0755); err != nil {
		return false, fmt.Errorf("failed to write run script to file %s: %v", runScript, err)
	}

	cmdRun := `isolate \
	-s \
	-b %d \
	-M %s \
	--stderr-to-stdout \
	-t %d \
	-w %d \
	-x 0 \
	-m %d \
	-k %d \
	-p4 \
	-f %d \
	-E "HOME=/tmp" \
	-E "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
	-d "/etc:noexec" \
	--run \
	-- /bin/bash %s < %s > %s 2> %s`

	actualRunCmd := fmt.Sprintf(cmdRun, j.BoxID, j.MetaFile, j.Submission.TimeLimit, j.Submission.WallTimeLimit, j.Submission.MemoryLimit, j.Submission.StackLimit, j.Submission.OutputLimit, filepath.Base(runScript), j.InputFile, j.OutputFile, j.ErrorFile)

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", actualRunCmd)
	err := cmd.Run()

	stdInText, readErr := os.ReadFile(j.InputFile)
	if readErr == nil && len(stdInText) > 0 {
		j.Response.Stdin = string(stdInText)
	}

	runOutputText, readErr := os.ReadFile(j.OutputFile)
	if readErr == nil && len(runOutputText) > 0 {
		j.Response.Stdout = string(runOutputText)
	}

	stderrOutputText, readErr := os.ReadFile(j.ErrorFile)
	if readErr == nil && len(stderrOutputText) > 0 {
		j.Response.Stderr = string(stderrOutputText)
	}

	metadata, _ := j.getMetadata()
	j.Response.ExitCode = metadata["exit-code"]
	j.Response.ExitSignal = metadata["exit-signal"]
	j.Response.Time = metadata["time"]
	j.Response.Memory = metadata["max-rss"]
	fmt.Println("----------------Run Metadata------------")
	fmt.Println(metadata)

	if err := exec.CommandContext(ctx, "sudo", "rm", "-rf", runScript).Run(); err != nil {
		return false, fmt.Errorf("failed to remove file %s: %v", runScript, err)
	}

	j.resetMetadata(ctx)

	if _, ok := err.(*exec.ExitError); ok {
		if status, ok := metadata["status"]; ok {
			switch status {
			case "TO":
				j.Response.Result = schema.ResultTimeLimitExceeded
				j.Response.Message = "Time Limit Exceeded"
			case "RE":
				j.Response.Result = schema.ResultRuntimeError
				j.Response.Message = "Runtime Error"
			}
			return false, nil
		}

	} else if err != nil {
		return false, err
	}

	if j.Response.Stdout == j.Submission.StdOut && j.Response.Stderr == "" {
		j.Response.Result = schema.ResultAccepted
	} else {
		j.Response.Result = schema.ResultWrongAnswer
	}

	return true, nil
}

func (j *IsolateJob) getMetadata() (map[string]string, error) {
	metadata := make(map[string]string)
	metadataText, err := os.ReadFile(j.MetaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %v", err)
	}

	lines := strings.SplitSeq(string(metadataText), "\n")
	for line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			metadata[parts[0]] = parts[1]
		}
	}

	return metadata, nil
}

func (j *IsolateJob) resetMetadata(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sudo", "rm", "-rf", j.MetaFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset metadata: %v", err)
	}
	j.InitializeFiles(j.MetaFile, ctx)

	return nil
}

func (j *IsolateJob) CleanUp(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "isolate", "-b", strconv.Itoa(j.BoxID), "--cleanup")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to cleanup isolate box: %v", err)
	}

	return nil
}
