package service

import (
	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/judger/common/database"
)

type _case struct {
	id     int
	input  string
	output string
}

func (j *judgeLeader) queryCase(message *queue.Message) ([]*_case, error) {
	rows, err := database.DB.Query("SELECT id, input, output FROM tb_case WHERE pid = ?", message.Pid)
	if err != nil {
		return nil, err
	}
	cases := make([]*_case, 0)
	var id int
	var input, output string
	for rows.Next() {
		if err := rows.Scan(&id, &input, &output); err != nil {
			return nil, err
		}
		cases = append(cases, &_case{
			id:     id,
			input:  input,
			output: output,
		})
	}
	return cases, nil
}
