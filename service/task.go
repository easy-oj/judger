package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ahmetb/go-dexec"
	"github.com/fsouza/go-dockerclient"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/proto/common"
	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/common/settings"
	"github.com/easy-oj/common/tools"
	"github.com/easy-oj/judger/common/docker_client"
)

type task struct {
	message   *queue.Message
	conf      *judgeConf
	cases     []*_case
	timeLimit int
	memLimit  int
	uuid      string
	dirPath   string
	docker    *dexec.Docker
}

func newTask(message *queue.Message, conf *judgeConf, cases []*_case, timeLimit, memLimit int) *task {
	uuid := tools.GenUUID()
	return &task{
		message:   message,
		conf:      conf,
		cases:     cases,
		timeLimit: timeLimit,
		memLimit:  memLimit,
		uuid:      uuid,
		dirPath:   path.Join(settings.Judger.Path, uuid),
		docker: &dexec.Docker{
			Client: docker_client.Client,
		},
	}
}

func (t *task) judge() (common.SubmissionStatus, string, string) {
	defer t.clean()
	if err := t.release(); err != nil {
		return common.SubmissionStatus_SYSTEM_ERROR, "", "[]"
	}
	if t.conf.NeedCompile {
		if bs, err := t.compile(); err != nil {
			return common.SubmissionStatus_COMPILE_ERROR, strings.TrimSpace(string(bs)), "[]"
		}
	}
	status := common.SubmissionStatus_ACCEPTED
	es := make([]*execution, len(t.cases))
	for idx, c := range t.cases {
		time.Sleep(time.Second)
		e := newExecution(t, c)
		e.execute()
		if e.s != common.ExecutionStatus_EXECUTION_ACCEPTED {
			status = common.SubmissionStatus_FAILURE
		}
		es[idx] = e
	}
	if bs, err := json.Marshal(es); err == nil {
		return status, "", string(bs)
	}
	return status, "", "[]"
}

func (t *task) release() error {
	if err := os.MkdirAll(t.dirPath, 0777); err != nil {
		logs.Error("[JudgeTask] sid = %d, mkdir error: %s", t.message.Sid, err.Error())
		return err
	}
	for p, c := range t.message.Content {
		if err := ioutil.WriteFile(path.Join(t.dirPath, p), []byte(c), 0666); err != nil {
			logs.Error("[JudgeTask] sid = %d, write file error: %s", t.message.Sid, err.Error())
			return err
		}
	}
	return nil
}

func (t *task) compile() ([]byte, error) {
	e, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Name: t.uuid,
		Config: &docker.Config{
			Image:           t.conf.CompileImage,
			WorkingDir:      "/ws",
			NetworkDisabled: true,
			Env:             t.conf.CompileEnvs,
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
				fmt.Sprintf("%s:/ws:rw", t.dirPath),
			},
		},
	})
	cmd := t.docker.Command(e, t.conf.CompileCmd, t.conf.CompileArgs...)
	return cmd.CombinedOutput()
}

func (t *task) clean() {
	if err := os.RemoveAll(t.dirPath); err != nil {
		logs.Error("[JudgeTask] sid = %d, remove error: %s", t.message.Sid, err.Error())
	}
}
