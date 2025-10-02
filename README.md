# Lighthouse-web3 Go SDK


### Features Implemented 

**Client & Config**
- API key auth via WithAPIKey or LIGHTHOUSE_API_KEY env var

- Configurable hosts, timeout, and user agent

**Error Handling**
- Structured error type with status and message

**Storage**

- File uploads (UploadFile, UploadReader) with optional progress callback

**Files**

- List uploaded files

- Get file info

- Delete file (by Lighthouse file ID)

- Pin file (by CID)

**Deals**
Query Filecoin deal status for a CID

**CLI (lhctl)**
```
--upload <path> : Upload file

--list : List uploaded files

--info <cid> : Get file info

--delete <id> : Delete file by ID

--deals <cid> : Check deal status
```

### Example Usage
Example Usage
###### Build CLI
```
go build -o lhctl ./cmd/lhctl
```
###### Upload file
```
LIGHTHOUSE_API_KEY=your-api-key ./lhctl --upload ./README.md
```
###### List files
```
LIGHTHOUSE_API_KEY=your-api-key ./lhctl --list
```
###### Get file info
```
LIGHTHOUSE_API_KEY=your-api-key ./lhctl --info <cid>
```
###### Delete file (by ID from --list)
```
LIGHTHOUSE_API_KEY=your-api-key ./lhctl --delete <id>
```

Quickstart (Go SDK)
```
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse"
    "github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

func main() {
    // Create client
    client := lighthouse.NewClient(nil,
        lighthouse.WithAPIKey(os.Getenv("LIGHTHOUSE_API_KEY")),
    )

    ctx := context.Background()

    // Upload file with progress
    result, err := client.Storage().UploadFile(ctx, "README.md",
        schema.WithProgress(func(p schema.Progress) {
            fmt.Printf("\rUploading: %.1f%%", p.Percent())
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nUploaded CID: %s\n", result.Data.Hash)

    // Get file info
    info, err := client.Files().Info(ctx, result.Data.Hash)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("File: %s (%s bytes)\n", info.FileName, info.FileSizeInBytes)

    // Check deals
    deals, err := client.Deals().Status(ctx, result.Data.Hash)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Deals: %d\n", len(deals))
}
```