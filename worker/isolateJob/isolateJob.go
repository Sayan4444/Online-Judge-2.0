package isolatejob

import (
	"OJ-Worker/config"
	model "OJ-Worker/models"
	"OJ-Worker/schema"
	"context"
	"fmt"
	"log"
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

func isRootUser() bool {
	return os.Getuid() == 0
}

type IsolateJob struct {
	Submission *schema.RabbitMQPayload
	Response   *schema.JudgeResponse
	Language   *model.Language
	TestCases  []model.TestCase
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

func ProcessSubmission(submission *schema.RabbitMQPayload, response *schema.JudgeResponse, ctx context.Context) error {
	log.Printf("Processing submission %s for problem %s by user %s", submission.SubmissionID, submission.ProblemID, submission.UserID)
	db := config.DB
	var language model.Language
	if err := db.Where("name = ?", submission.Language).First(&language).Error; err != nil {
		return fmt.Errorf("failed to find language: %v", err)
	}

	var testCases []model.TestCase
	if err := db.Where("problem_id = ?", submission.ProblemID).Find(&testCases).Error; err != nil {
		return err
	}

	job := &IsolateJob{
		Submission: submission,
		BoxID:      int(atomic.AddInt64(&boxIDCounter, 1)) % 2147483647,
		Response:   response,
		Language:   &language,
		TestCases:  testCases,
	}

	err := job.execute(ctx)
	log.Printf("worker-1: JudgeResponse fields for submission %s:", submission.SubmissionID)
	log.Printf("  Stderr: %s", response.Stderr)
	log.Printf("  Time: %s", response.Time)
	log.Printf("  Memory: %s", response.Memory)
	log.Printf("  ExitCode: %s", response.ExitCode)
	log.Printf("  Result: %s", response.Result)
	log.Printf("  CompileOutput: %s", response.CompileOutput)
	return err
}

func (j *IsolateJob) execute(ctx context.Context) error {
	defer j.cleanUp(ctx)
	if err := j.initializeIsolate(ctx); err != nil {
		j.Response.Result = schema.ResultSystemError
		return fmt.Errorf("failed to initialize isolate: %v", err)
	}
	log.Println("Isolate initialized successfully")
	success, err := j.compile(ctx)
	if err != nil {
		j.Response.Result = schema.ResultSystemError
		return fmt.Errorf("failed to compile: %v", err)
	}
	if !success {
		return nil
	}
	log.Println("Code Compiled successfully")

	success, err = j.run(ctx)
	if err != nil {
		j.Response.Result = schema.ResultSystemError
		return fmt.Errorf("failed to run: %v", err)
	}
	if !success {
		return nil
	}

	j.cleanUp(ctx)
	log.Println("Code Ran successfully")
	return nil

}

func (j *IsolateJob) initializeIsolate(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "isolate",
		"-b", strconv.Itoa(j.BoxID),
		"--init",
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to initialize isolate box: %v", err)
	}

	j.WorkDir = strings.TrimSpace(string(output))
	log.Printf("Isolate work directory: %s", j.WorkDir)
	j.BoxDir = filepath.Join(j.WorkDir, "box")
	j.TmpDir = filepath.Join(j.WorkDir, "tmp")
	j.SourceFile = filepath.Join(j.BoxDir, j.Language.SrcFile)
	j.InputFile = filepath.Join(j.WorkDir, StdinFileName)
	j.OutputFile = filepath.Join(j.WorkDir, StdoutFileName)
	j.ErrorFile = filepath.Join(j.WorkDir, StderrFileName)
	j.MetaFile = filepath.Join(j.WorkDir, MetadataFileName)

	files := []string{j.SourceFile, j.InputFile, j.OutputFile, j.ErrorFile, j.MetaFile}
	for _, file := range files {
		if err := j.initializeFiles(file, ctx); err != nil {
			return fmt.Errorf("failed to initialize file %s: %v", file, err)
		}
	}

	if err := os.WriteFile(j.SourceFile, []byte(j.Submission.SourceCode), 0644); err != nil {
		return fmt.Errorf("failed to write source code to file %s: %v", j.SourceFile, err)
	}

	return nil
}

func (j *IsolateJob) initializeFiles(filename string, ctx context.Context) error {
	user := os.Getenv("USER")
	var cmdStr string
	if isRootUser() {
		cmdStr = fmt.Sprintf("touch %s && chown %s: %s", filename, user, filename)
	} else {
		cmdStr = fmt.Sprintf("sudo touch %s && sudo chown %s: %s", filename, user, filename)
	}
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", cmdStr)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize file %s: %v", filename, err)
	}

	return nil
}

func (j *IsolateJob) compile(ctx context.Context) (bool, error) {
	log.Println("Starting compilation process...")
	// making compile script
	compileScript := filepath.Join(j.BoxDir, "compile.sh")
	compileOutput := filepath.Join(j.WorkDir, "compile_output.txt")
	if err := j.initializeFiles(compileOutput, ctx); err != nil {
		return false, fmt.Errorf("failed to initialize compile output file: %v", err)
	}

	log.Printf("Compile command from language config: %s\n", j.Language.CompileCommand)
	if err := os.WriteFile(compileScript, []byte(j.Language.CompileCommand), 0755); err != nil {
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
	-x %d \
	-m %d \
	-k %d \
	-p120 \
	-f %d \
	-E "HOME=/tmp" \
	-E "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
	-d "/etc:noexec" \
	--run \
	-- /bin/bash %s > %s`

	// Time limits in seconds
	compilationTimeLimit := 5
	compilationWallTimeLimit := 20
	compilationExtraTimeLimit := 2 // Extra time for startup

	// Memory limits in Kilobytes (KB)
	compilationMemoryLimit := 512000 // 512 MB
	compilationStackLimit := 128000  // 128 MB

	actualCompileCmd := fmt.Sprintf(cmdRun, j.BoxID, j.MetaFile, compilationTimeLimit, compilationWallTimeLimit, compilationExtraTimeLimit, compilationMemoryLimit, compilationStackLimit, j.Language.OutputLimit, filepath.Base(compileScript), compileOutput)
	log.Printf("Executing actual compile command:\n%s\n", actualCompileCmd)

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", actualCompileCmd)
	err := cmd.Run()

	compileOutputText, readErr := os.ReadFile(compileOutput)
	if readErr != nil {
		log.Printf("Error reading compile output file %s: %v\n", compileOutput, readErr)
	}
	if len(compileOutputText) > 0 {
		j.Response.CompileOutput = string(compileOutputText)
		log.Printf("Compile output:\n%s\n", j.Response.CompileOutput)
	} else {
		log.Println("Compile output is empty.")
	}

	// Checking the correctness
	metadata, metaErr := j.getMetadata()
	if metaErr != nil {
		log.Printf("Error getting metadata after compile: %v\n", metaErr)
	} else {
		log.Printf("Compilation metadata: %v\n", metadata)
	}

	filesToRemove := []string{compileScript, compileOutput}

	for _, file := range filesToRemove {
		var rmCmd *exec.Cmd
		if isRootUser() {
			rmCmd = exec.CommandContext(ctx, "rm", "-rf", file)
		} else {
			rmCmd = exec.CommandContext(ctx, "sudo", "rm", "-rf", file)
		}
		if err := rmCmd.Run(); err != nil {
			// Log this error but don't fail the entire job because of a cleanup issue
			log.Printf("Warning: failed to remove file %s: %v\n", file, err)
		}
	}

	if err := j.resetMetadata(ctx); err != nil {
		log.Printf("Warning: failed to reset metadata after compile: %v\n", err)
	}

	if err != nil {
		log.Printf("Compile command finished with error: %v\n", err)
		if _, ok := err.(*exec.ExitError); ok {
			if metaErr == nil {
				j.Response.ExitCode = metadata["exitcode"]
				if status, ok := metadata["status"]; ok {
					log.Printf("Compilation status from metadata: %s\n", status)
					if status == "TO" {
						j.Response.Result = schema.ResultCompileTimeLimitExceeded
						log.Println("Result set to: Compile Time Limit Exceeded")
					} else {
						j.Response.Result = schema.ResultCompileError
						log.Println("Result set to: Compile Error due to status " + status)
					}
				} else {
					j.Response.Result = schema.ResultCompileError
					log.Println("Result set to: Compile Error (status not in metadata)")
				}
			} else {
				j.Response.Result = schema.ResultSystemError
				log.Println("Result set to: System Error (could not read metadata after compile error)")
			}
			return false, nil
		}
		// For other errors (not ExitError), it's likely a system issue
		return false, fmt.Errorf("compile command failed with non-exit error: %v", err)
	}

	log.Println("Compilation successful.")
	return true, nil
}

func (j *IsolateJob) run(ctx context.Context) (bool, error) {
	runScript := filepath.Join(j.BoxDir, "run.sh")

	if err := os.WriteFile(runScript, []byte(j.Language.RunCommand), 0755); err != nil {
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
		-p120 \
		-f %d \
		-E "HOME=/tmp" \
		-E "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
		-d "/etc:noexec" \
		--run \
		-- /bin/bash %s < %s > %s 2> %s`

	actualRunCmd := fmt.Sprintf(cmdRun, j.BoxID, j.MetaFile, j.Language.TimeLimit, j.Language.WallLimit, j.Language.MemoryLimit, j.Language.StackLimit, j.Language.OutputLimit, filepath.Base(runScript), j.InputFile, j.OutputFile, j.ErrorFile)
	log.Printf("Actual run command template:\n%s", actualRunCmd)

	
	var combinedInput strings.Builder
	numTestCasesStr := strconv.Itoa(len(j.TestCases))
	combinedInput.WriteString(numTestCasesStr + " ")
	for _, testCase := range j.TestCases {
		combinedInput.WriteString(testCase.Input + " ")
	}

	if err := os.WriteFile(j.InputFile, []byte(combinedInput.String()), 0755); err != nil {
		return false, fmt.Errorf("failed to write combined stdin to file %s: %v", j.InputFile, err)
	}

	
	success, err := j.executeTestCase(ctx, actualRunCmd)
	if err != nil {
		return false, err
	}
	
	var rmCmd *exec.Cmd
	if isRootUser() {
		rmCmd = exec.CommandContext(ctx, "rm", "-rf", runScript)
	} else {
		rmCmd = exec.CommandContext(ctx, "sudo", "rm", "-rf", runScript)
	}
	if err := rmCmd.Run(); err != nil {
		return false, fmt.Errorf("failed to remove file %s: %v", runScript, err)
	}

	if !success {
		j.Response.Result = schema.ResultWrongAnswer
		return false, nil
	} else {
		j.Response.Result = schema.ResultAccepted
		return true, nil
	}
}

func (j *IsolateJob) executeTestCase(ctx context.Context, actualRunCmd string) (bool, error) {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", actualRunCmd)
	err := cmd.Run()

	var stdout string
	runOutputText, readErr := os.ReadFile(j.OutputFile)
	if readErr == nil && len(runOutputText) > 0 {
		stdout = string(runOutputText)
	}

	stderrOutputText, readErr := os.ReadFile(j.ErrorFile)
	if readErr == nil && len(stderrOutputText) > 0 {
		j.Response.Stderr = string(stderrOutputText)
	}

	metadata, metaErr := j.getMetadata()
	if metaErr != nil {
		log.Printf("Error getting metadata after run: %v", metaErr)
	} else {
		log.Printf("Run metadata: %v", metadata)
	}
	j.Response.ExitCode = metadata["exit-code"]
	if currentTime := metadata["time"]; currentTime != "" {
		if j.Response.Time == "" || currentTime > j.Response.Time {
			j.Response.Time = currentTime
		}
	}
	if currentMemory := metadata["max-rss"]; currentMemory != "" {
		if j.Response.Memory == "" || currentMemory > j.Response.Memory {
			j.Response.Memory = currentMemory
		}
	}
	j.resetMetadata(ctx)

	if _, ok := err.(*exec.ExitError); ok {
		if status, ok := metadata["status"]; ok {
			switch status {
			case "TO":
				j.Response.Result = schema.ResultTimeLimitExceeded
			case "RE":
				j.Response.Result = schema.ResultRuntimeError
			}
			return false, nil
		}

	} else if err != nil {
		return false, err
	}

	actualOutputs := strings.Split(strings.TrimSpace(stdout), "\n")
	for i, output := range actualOutputs {
		testCase := j.TestCases[i]
		expectedOutput := strings.TrimSpace(testCase.Output)
		actualOutput := strings.TrimSpace(output)

		if actualOutput != expectedOutput {
			log.Printf("Wrong answer for test case index %d, ID %s. Expected: '%s', Got: '%s'", i, testCase.ID, expectedOutput, actualOutput)
			j.Response.Result = schema.ResultWrongAnswer
			j.Response.WrongAnswers = append(j.Response.WrongAnswers, schema.WrongAnswer{
				TestCaseID:     testCase.ID.String(),
				Stdout:   actualOutput,
			})

			return false, nil
		}
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
	var cmd *exec.Cmd
	if isRootUser() {
		cmd = exec.CommandContext(ctx, "rm", "-rf", j.MetaFile)
	} else {
		cmd = exec.CommandContext(ctx, "sudo", "rm", "-rf", j.MetaFile)
	}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset metadata: %v", err)
	}
	j.initializeFiles(j.MetaFile, ctx)

	return nil
}

func (j *IsolateJob) cleanUp(ctx context.Context) error {
	log.Printf("Cleaning up isolate box %d", j.BoxID)
	cmd := exec.CommandContext(ctx, "isolate", "-b", strconv.Itoa(j.BoxID), "--cleanup")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to cleanup isolate box: %v", err)
	}

	return nil
}
