package queue

import "github.com/easy-oj/common/proto/base"

func NewPutMessageReq() *PutMessageReq {
	return &PutMessageReq{
		BaseReq: base.NewBaseReq(),
	}
}

func NewPutMessageResp() *PutMessageResp {
	return &PutMessageResp{
		BaseResp: base.NewBaseResp(),
	}
}

func NewGetMessageReq() *GetMessageReq {
	return &GetMessageReq{
		BaseReq: base.NewBaseReq(),
	}
}

func NewGetMessageResp() *GetMessageResp {
	return &GetMessageResp{
		BaseResp: base.NewBaseResp(),
	}
}
