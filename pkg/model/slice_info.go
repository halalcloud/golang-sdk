package model

type SliceInfo struct {
	Cid       string
	Size      int64
	Completed bool
}

func NewSliceInfo(slice string) *SliceInfo {
	return &SliceInfo{
		Cid:       slice,
		Completed: false,
	}
}
