package lighthouse

import (
	"context"
	"io"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type StorageService interface {
	UploadFile(ctx context.Context, path string, opts ...schema.UploadOption) (*schema.UploadResult, error)
	UploadReader(ctx context.Context, name string, size int64, r io.Reader, opts ...schema.UploadOption) (*schema.UploadResult, error)
}

type FilesService interface {
	List(ctx context.Context, lastKey *string) (*schema.FileList, error)
	Info(ctx context.Context, cid string) (*schema.FileInfo, error)
	Pin(ctx context.Context, cid, name string) error	
	Delete(ctx context.Context, id string) error
}

type DealsService interface {
	Status(ctx context.Context, cid string) ([]schema.DealStatus, error)
}