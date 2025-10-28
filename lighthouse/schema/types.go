package schema

import "io"

type UploadResult struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

type Progress struct {
	Uploaded int64
	Total    int64
}

type IPNSPublishResponse struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type FileEntry struct {
	Name            string `json:"fileName"`
	CID             string `json:"cid"`
	Size            int64  `json:"fileSize"`
	ID              string `json:"id,omitempty"`
	PublicKey       string `json:"publicKey,omitempty"`
	FileSizeInBytes int64  `json:"fileSizeInBytes"`
	FileSizeStr     string `json:"FileSizeStr,omitempty"`
	MimeType        string `json:"mimeType,omitempty"`
	TxHash          string `json:"txHash,omitempty"`
	Status          string `json:"status,omitempty"`
	CreatedAt       int64  `json:"createdAt,omitempty"`
	LastUpdate      int64  `json:"lastUpdate,omitempty"`
	Encryption      bool   `json:"encryption,omitempty"`
}

type FileList struct {
	Data       []FileEntry `json:"fileList"`
	LastKey    *string     `json:",omitempty"`
	TotalFiles *int        `json:"totalFiles,omitempty"`
}

type FileInfo struct {
	FileSizeInBytes int64  `json:"fileSizeInBytes"`
	CID             string `json:"cid"`
	Encryption      bool   `json:"encryption"`
	FileName        string `json:"fileName"`
	MimeType        string `json:"mimeType"`
}

type DealStatus struct {
	ChainDealID        int64  `json:"chainDealID"`
	EndEpoch           int64  `json:"endEpoch"`
	PublishCID         string `json:"publishCID"`
	StorageProvider    string `json:"storageProvider"`
	DealStatus         string `json:"dealStatus"`
	BundleID           string `json:"bundleId"`
	DealUUID           string `json:"dealUUID"`
	StartEpoch         int64  `json:"startEpoch"`
	AggregateIn        string `json:"aggregateIn"`
	ProviderCollateral string `json:"providerCollateral"`
	PieceCID           string `json:"pieceCID"`
	PayloadCID         string `json:"payloadCid"`
	PieceSize          int64  `json:"pieceSize"`
	CarFileSize        int64  `json:"carFileSize"`
	LastUpdate         int64  `json:"lastUpdate"`
	DealID             int64  `json:"dealId"`
	Miner              string `json:"miner"`
	Content            int64  `json:"content"`
}

type DealStatusResponse struct {
	Data []DealStatus `json:"data"`
}

type Usage struct {
	DataLimit int64 `json:"dataLimit"`
	DataUsed  int64 `json:"dataUsed"`
}

type IPNSKeyResponse struct {
	IPNSName string `json:"ipnsName"`
	IPNSId   string `json:"ipnsId"`
}

type IPNSRecord struct {
	IPNSName   string `json:"ipnsName"`
	IPNSId     string `json:"ipnsId"`
	PublicKey  string `json:"publicKey"`
	CID        string `json:"cid"`
	LastUpdate int64  `json:"lastUpdate"`
}

type IPNSKey struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
}

type IPNSRemoveResponse struct {
	Keys []IPNSKey `json:"Keys"`
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
	OnProgress ProgressCallback
}

func DefaultUploadOptions() *UploadOptions {
	return &UploadOptions{
		ChunkSize:  8 << 20,
		Metadata:   map[string]string{},
		MimeType:   "",
		Pin:        false,
		Public:     true,
		Progress:   nil,
		OnProgress: nil,
		EncryptKey: nil,
	}
}

func (p Progress) Percent() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(p.Uploaded) * 100.0 / float64(p.Total)
}

type ProgressCallback func(Progress)

func WithMimeType(mt string) UploadOption { return func(o *UploadOptions) { o.MimeType = mt } }
func WithPin() UploadOption               { return func(o *UploadOptions) { o.Pin = true } }
func WithPrivate() UploadOption           { return func(o *UploadOptions) { o.Public = false } }
func WithProgress(cb ProgressCallback) UploadOption {
	return func(o *UploadOptions) { o.OnProgress = cb }
}
