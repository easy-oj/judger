package service

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/periodic"
	"github.com/easy-oj/common/proto/common"
	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/common/settings"
	"github.com/easy-oj/judger/common/caller"
	"github.com/easy-oj/judger/common/database"
)

type judgeLeader struct {
	worker int32
	confs  map[int32]*judgeConf
}

func StartJudgeService() {
	j := newJudgeLeader()
	j.queryConf()
	periodic.NewStaticPeriodic(j.queryConf, time.Minute, periodic.MinInterval).Start()
	periodic.NewStaticPeriodic(j.tryJudge, time.Second, periodic.MinInterval).Start()
}

func newJudgeLeader() *judgeLeader {
	return &judgeLeader{}
}

func (j *judgeLeader) tryJudge() {
	if atomic.LoadInt32(&j.worker) >= int32(settings.Judger.Worker) {
		return
	}
	req := queue.NewGetMessageReq()
	resp, err := caller.QueueClient.GetMessage(context.Background(), req)
	if err != nil {
		logs.Error("[JudgeLeader] call QueueClient.GetMessage error: %s", err.Error())
	} else if resp.Message != nil {
		atomic.AddInt32(&j.worker, 1)
		go j.judge(resp.Message)
	}
}

func (j *judgeLeader) judge(message *queue.Message) {
	defer atomic.AddInt32(&j.worker, -1)
	j.setStatus(message, common.SubmissionStatus_JUDGING)
	status, ce, executions := common.SubmissionStatus_SYSTEM_ERROR, "", "[]"
	defer func() {
		j.finalize(message, status, ce, executions)
	}()
	conf, ok := j.confs[message.Lid]
	if !ok {
		logs.Warn("[JudgeLeader] sid = %d, unsupported lid: %d", message.Sid, message.Lid)
		return
	}
	cases, err := j.queryCase(message)
	if err != nil {
		logs.Error("[JudgeLeader] sid = %d, query case error: %s", message.Sid, err.Error())
		return
	}
	timeLimit, memLimit, err := j.queryLimit(message)
	if err != nil {
		logs.Error("[JudgeLeader] sid = %d, query limit error: %s", message.Sid, err.Error())
		return
	}
	status, ce, executions = newTask(message, conf, cases, timeLimit, memLimit).judge()
}

func (j *judgeLeader) setStatus(message *queue.Message, status common.SubmissionStatus) {
	_, err := database.DB.Exec("UPDATE ENTITY__SUBMISSION SET status = ? WHERE id = ?", status.String(), message.Sid)
	if err != nil {
		logs.Error("[JudgeLeader] sid = %d, set status error: %s", message.Sid, err.Error())
	}
}

func (j *judgeLeader) finalize(message *queue.Message, status common.SubmissionStatus, ce, executions string) {
	if length := len(ce); length > 4000 {
		ce = fmt.Sprintf("%s\n......\nCompile error message is too long with total %d bytes!", ce[:4000], length)
	}
	_, err := database.DB.Exec("UPDATE ENTITY__SUBMISSION SET status = ?, compile_error = ?, executions = ? WHERE id = ?", status.String(), ce, executions, message.Sid)
	if err != nil {
		logs.Error("[JudgeLeader] sid = %d, finalize error: %s", message.Sid, err.Error())
	}
	if status != common.SubmissionStatus_ACCEPTED {
		return
	}
	_, err = database.DB.Exec("UPDATE ENTITY__PROBLEM SET accepted_count = accepted_count + 1 WHERE id = ?", message.Pid)
	if err != nil {
		logs.Error("[JudgeLeader] sid = %d, update accepted count error: %s", message.Sid, err.Error())
	}
}
