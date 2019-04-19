package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/ahmetb/go-dexec"
	"github.com/fsouza/go-dockerclient"

	"github.com/easy-oj/common/proto/common"
	"github.com/easy-oj/common/settings"
	"github.com/easy-oj/judger/tools"
)

type execution struct {
	t *task  `json:"-"`
	c *_case `json:"-"`

	s common.ExecutionStatus `json:"-"`

	Cid       int    `json:"cid"`
	Status    string `json:"status"`
	TimeUsed  int    `json:"time_used"`
	MemUsed   int    `json:"mem_used"`
	ExtraData int    `json:"extra_data"`
}

func newExecution(t *task, c *_case) *execution {
	return &execution{
		t:   t,
		c:   c,
		Cid: c.id,
	}
}

func (e *execution) execute() {
	defer func() {
		e.Status = e.s.String()
	}()
	exec, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Name: e.t.uuid,
		Config: &docker.Config{
			Image:           e.t.conf.ExecuteImage,
			WorkingDir:      "/ws",
			NetworkDisabled: true,
		},
		HostConfig: &docker.HostConfig{
			CPUPeriod: 100000,
			CPUQuota:  100000,
			Memory:    512 * 1024 * 1024,
			ShmSize:   32 * 1024 * 1024,
			Tmpfs: map[string]string{
				"/tmp": "rw,noexec,nosuid,size=32768k",
			},
			ReadonlyRootfs: true,
			VolumeDriver:   "bind",
			Binds: []string{
				fmt.Sprintf("%s:/ws:ro", e.t.dirPath),
				fmt.Sprintf("%s:%s:ro", settings.Path.Executor, settings.Path.Executor),
			},
			CapAdd: []string{"SYS_PTRACE"},
		},
	})
	args := append([]string{strconv.Itoa(e.t.timeLimit), strconv.Itoa(e.t.memLimit)}, e.t.conf.ExecuteTarget...)
	cmd := e.t.docker.Command(exec, settings.Path.Executor, args...)
	cmd.Stdin = bytes.NewReader([]byte(e.c.input))
	bs, err := cmd.Output()
	output := strings.TrimSpace(string(bs))
	if idx := strings.LastIndexByte(output, '#'); idx >= 0 {
		ss := strings.Split(output[idx+1:], ":")
		e.TimeUsed, _ = strconv.Atoi(ss[0])
		e.MemUsed, _ = strconv.Atoi(ss[1])
		e.ExtraData, _ = strconv.Atoi(ss[2])
		if strings.HasSuffix(output, fmt.Sprintf("#%d:%d:%d", e.TimeUsed, e.MemUsed, e.ExtraData)) {
			output = strings.TrimSpace(output[:idx])
		}
	}
	if err != nil {
		if ee, ok := err.(*dexec.ExitError); ok {
			switch ee.ExitCode {
			case 2:
				e.s = common.ExecutionStatus_RUNTIME_ERROR
			case 3:
				e.s = common.ExecutionStatus_TIME_LIMIT_EXCEED
			case 4:
				e.s = common.ExecutionStatus_MEMORY_LIMIT_EXCEED
			case 5:
				e.s = common.ExecutionStatus_ILLEGAL_SYSCALL
			}
		} else {
			e.s = common.ExecutionStatus_EXECUTION_ERROR
		}
	} else {
		if strings.TrimSpace(e.c.output) == strings.TrimSpace(output) {
			e.s = common.ExecutionStatus_EXECUTION_ACCEPTED
		} else if tools.StripWriteSpace(e.c.output) == tools.StripWriteSpace(output) {
			e.s = common.ExecutionStatus_PRESENTATION_ERROR
		} else {
			e.s = common.ExecutionStatus_WRONG_ANSWER
		}
	}
}
