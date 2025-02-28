package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

var (
	ErrUploadBlockstorePutFailed = errors.New("upload blockstore put failed")
)

type hasResponse struct {
	Has bool `json:"has"`
}

type uploadBlockstore struct {
	client  *http.Client
	baseURL string
	taskID  string
	rootCID string
}

func NewUploadBlockstore(httpClient *http.Client, baseURL string, taskID string, rootCID string) blockstore.Blockstore {
	return &uploadBlockstore{
		client:  httpClient,
		baseURL: baseURL,
		taskID:  taskID,
		rootCID: rootCID,
	}
}

func (s *uploadBlockstore) DeleteBlock(ctx context.Context, req cid.Cid) error {
	return nil
}

// Has(context.Context, cid.Cid) (bool, error)
func (s *uploadBlockstore) Has(ctx context.Context, chk cid.Cid) (bool, error) {
	requestURL := s.baseURL + "/" + s.taskID + "/" + chk.String()
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)

		log.Printf("Has failed: %s => %s", resp.Status, string(responseBody))
		return false, ErrUploadBlockstorePutFailed
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var hasResp hasResponse
	if err := json.Unmarshal(data, &hasResp); err != nil {
		return false, err
	}
	nHasResp := hasResp.Has
	return nHasResp, nil
}

// Get(context.Context, cid.Cid) (blocks.Block, error)
func (s *uploadBlockstore) Get(ctx context.Context, req cid.Cid) (blocks.Block, error) {

	// data, err := storagePool.Read(oid, )
	return nil, nil
}

// GetSize(context.Context, cid.Cid) (int, error)
func (s *uploadBlockstore) GetSize(ctx context.Context, req cid.Cid) (int, error) {

	return 0, nil
}

// Put(context.Context, blocks.Block) error
func (s *uploadBlockstore) Put(ctx context.Context, block blocks.Block) error {
	requestURL := s.baseURL + "/" + s.taskID + "/" + block.Cid().String()
	postData := block.RawData()
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(postData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)

		log.Printf("Put failed: %s => %s", resp.Status, string(responseBody))
		return ErrUploadBlockstorePutFailed
	}
	return nil
}

// PutMany(context.Context, []blocks.Block) error
func (s *uploadBlockstore) PutMany(ctx context.Context, blocks []blocks.Block) error {
	for _, block := range blocks {
		if err := s.Put(ctx, block); err != nil {
			return err
		}
	}
	return nil
}

// AllKeysChan(ctx context.Context) (<-chan cid.Cid, error)
func (s *uploadBlockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, nil
}

// HashOnRead(enabled bool)
func (s *uploadBlockstore) HashOnRead(enabled bool) {
}
