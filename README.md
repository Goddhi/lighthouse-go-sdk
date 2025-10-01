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
    "os"

    "github.com/lighthouse-web3/lighthouse-go/lighthouse"
)

## SDK Quickstart
```
func main() {
    client, _ := lighthouse.NewClient(
        lighthouse.WithAPIKey(os.Getenv("LIGHTHOUSE_API_KEY")),
    )

    // Upload
    file, _ := os.Open("README.md")
    result, _ := client.Storage.UploadReader(context.Background(), "README.md", 0, file)
    fmt.Println("Uploaded CID:", result.Hash)

    // Info
    info, _ := client.Files.Info(context.Background(), result.Hash)
    fmt.Printf("File Info: %+v\n", info)

    // Deals
    deals, _ := client.Deals.Status(context.Background(), result.Hash)
    fmt.Printf("Deals: %+v\n", deals)
}
```