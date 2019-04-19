package service

import (
	"strings"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/judger/common/database"
)

type judgeConf struct {
	NeedCompile   bool
	CompileImage  string
	CompileCmd    string
	CompileArgs   []string
	CompileEnvs   []string
	ExecuteImage  string
	ExecuteTarget []string
	TimeLimit     int
	MemLimit      int
}

func (j *judgeLeader) queryConf() {
	rows, err := database.DB.Query("SELECT id, need_compile, compile_image, compile_cmd, compile_env, execute_image, execute_cmd, time_limit, mem_limit FROM ENTITY__LANGUAGE")
	if err != nil {
		logs.Error("[JudgerLeader] query conf error: %s", err.Error())
		return
	}
	confs := make(map[int32]*judgeConf)
	var (
		id           int32
		needCompile  string
		compileImage string
		compileCmd   string
		compileEnv   string
		executeImage string
		executeCmd   string
		timeLimit    int
		memLimit     int
	)
	for rows.Next() {
		err := rows.Scan(&id, &needCompile, &compileImage, &compileCmd, &compileEnv, &executeImage, &executeCmd, &timeLimit, &memLimit)
		if err != nil {
			logs.Error("[JudgerLeader] query conf scan row error: %s", err.Error())
			continue
		}
		conf := &judgeConf{
			ExecuteImage:  executeImage,
			ExecuteTarget: strings.Split(executeCmd, " "),
			TimeLimit:     timeLimit,
			MemLimit:      memLimit,
		}
		if needCompile == "\x01" {
			cmd := strings.Split(compileCmd, " ")
			conf.NeedCompile = true
			conf.CompileImage = compileImage
			conf.CompileCmd = cmd[0]
			conf.CompileArgs = cmd[1:]
			if compileEnv != "" {
				conf.CompileEnvs = strings.Split(compileEnv, " ")
			}
		}
		confs[id] = conf
	}
	j.confs = confs
}
