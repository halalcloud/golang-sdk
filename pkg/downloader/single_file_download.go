package downloader

import (
	"context"
	"errors"
	"log"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	badger "github.com/dgraph-io/badger/v4"
)

func (s *SliceDownloader) StartDownloadSingleFile(path string, id string, savePath string) (err error) {
	if s.rtcEnabled {
		log.Printf("rtc already enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// 1. get file slice

	result, err := s.fileClient.ParseFileSlice(ctx, &pubUserFile.File{

		// Parent: &pubUserFile.File{Path: currentDir},
		Path:     path,
		Identity: id,
	})
	if err != nil {
		return err
	}

	//s.sliceDB.newsc
	tx := s.sliceDB.NewTransaction(true)
	bti, err := tx.Get([]byte("_task_" + path + "_" + id))
	if err == nil {
		// TODO: check status, if task not started, start it
		s.startIfStopTask(bti)
		return errors.New("task already exists")
	}
	// create task meta...
	err = s.createTaskMeta(result)
	if err != nil {
		return err
	}
	return nil

}

func (s *SliceDownloader) startIfStopTask(*badger.Item) (err error) {
	return nil
}

func (s *SliceDownloader) StopDownloadSingleFile(path string, id string) (err error) {
	return nil
}
