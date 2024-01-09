package model

import (
	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
)

type TaskMeta struct {
	slices []*SliceInfo
	info   *pubUserFile.ParseFileSliceResponse
}

func NewTaskMeta(result *pubUserFile.ParseFileSliceResponse, downloadPath string) *TaskMeta {
	slices := make([]*SliceInfo, 0)
	for _, slice := range result.RawNodes {
		sliceInfo := NewSliceInfo(slice)
		slices = append(slices, sliceInfo)
	}
	return &TaskMeta{slices: slices, info: result}
}
