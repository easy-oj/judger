package initial

import (
	"github.com/easy-oj/judger/common/caller"
	"github.com/easy-oj/judger/common/database"
	"github.com/easy-oj/judger/common/docker_client"
	"github.com/easy-oj/judger/service"
)

func Initialize() {
	caller.InitCaller()
	database.InitDatabase()
	docker_client.InitDockerClient()

	service.StartJudgeService()
}
