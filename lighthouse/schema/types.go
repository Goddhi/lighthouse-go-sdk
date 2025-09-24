package schema

import "io"


type UploadResult struct {
	Data struct {
		Hash string `json:"Hash"`
	} `json:"data"`
}

type Progress struct {
	Uploaded int64
	Total    int64
}

type FileEntry struct {
	Name string `json:"fileName"`
	CID  string `json:"cid"`
	Size int64  `json:"fileSize"`
}

type FileList struct {
	Data    []FileEntry `json:"data"`
	LastKey *string     `json:"lastKey"`
}

type FileInfo struct {
	Data struct {
		Hash     string `json:"Hash"`
		Name     string `json:"Name"`
		Size     string `json:"Size"`
		Type     string `json:"Type"`
		NumLinks int    `json:"NumLinks"`
	} `json:"data"`
}

type Reader = io.Reader

type UploadOption func(*UploadOptions)

type UploadOptions struct {
	MimeType   string
	Pin        bool
	Public     bool
	Metadata   map[string]string
	ChunkSize  int
	Progress   io.Writer
	EncryptKey []byte
	OnProgress ProgressCallback // currently unused (no-op)

}

type ProgressCallback func(Progress)


func DefaultUploadOptions() *UploadOptions {
	return &UploadOptions{
		ChunkSize:  8 << 20,
		Metadata:   map[string]string{},
		MimeType:   "",
		Pin:        false,
		Public:     true,
		Progress:   nil,
		EncryptKey: nil,
	}
}

func (p Progress) Percent() float64 {
	if p.Total == 0 { return 0 }
	return float64(p.Uploaded) * 100.0 / float64(p.Total)
}

func WithProgress(cb ProgressCallback) UploadOption {
	return func(o *UploadOptions) { o.OnProgress = cb }
}


