package service

import (
	"encoding/json"

	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/judger/common/database"
)

func (j *judgeLeader) queryLimit(message *queue.Message) (int, int, error) {
	conf := j.confs[message.Lid]
	timeLimit, memLimit := conf.TimeLimit, conf.MemLimit
	var str string
	row := database.DB.QueryRow("SELECT special_limits FROM ENTITY__PROBLEM WHERE id = ?", message.Pid)
	if err := row.Scan(&str); err != nil {
		return timeLimit, memLimit, err
	}
	var m map[int32]map[string]int
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return timeLimit, memLimit, err
	}
	if conf, ok := m[message.Lid]; ok {
		timeLimit, memLimit = conf["time_limit"], conf["mem_limit"]
	}
	return timeLimit, memLimit, nil
}
