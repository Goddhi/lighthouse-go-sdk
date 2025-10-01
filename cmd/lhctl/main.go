package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

func valueOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func valueOrNil(p *string) string {
	if p == nil || *p == "" {
		return "<nil>"
	}
	return *p
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("LIGHTHOUSE_API_KEY")
	if apiKey == "" {
		log.Fatal("set LIGHTHOUSE_API_KEY env var")
	}

	upload := flag.String("upload", "", "file path to upload")
	info := flag.String("info", "", "CID to fetch info for")
	list := flag.Bool("list", false, "list uploaded files")
	lastKey := flag.String("last-key", "", "pagination cursor for --list")
	deals := flag.String("deals", "", "CID to fetch Filecoin deal status")
	del := flag.String("delete", "", "file ID to delete (use --list to find IDs)")
	flag.Parse()

	cli := lighthouse.NewClient(nil, lighthouse.WithAPIKey(apiKey))


	switch {
	case *upload != "":
		res, err := cli.Storage().UploadFile(
			ctx,
			*upload,
			schema.WithProgress(func(p schema.Progress) {
				fmt.Printf("\ruploaded %d/%d (%.2f%%)", p.Uploaded, p.Total, p.Percent())
			}),
		)
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Uploaded %s (%s bytes)\nCID: %s\n", res.Name,  res.Size, res.Hash)

	case *info != "":
		i, err := cli.Files().Info(ctx, *info)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name=%s  Size=%v CID=%s MimeType=%s Encryption=%v\n", i.FileName, i.FileSizeInBytes, i.CID, i.MimeType, i.Encryption)

	case *list:
		var cursor *string
		if *lastKey != "" {
			cursor = lastKey
		}
		ls, err := cli.Files().List(ctx, cursor)
		if err != nil {
			log.Fatal(err)
		}
		// Show a few helpful columns (including ID so users can --delete)
		fmt.Printf("TOTAL=%d  NEXT=%v\n", valueOrZero(ls.TotalFiles), valueOrNil(ls.LastKey))
		fmt.Println("CID\tID\tSIZE(bytes)\tNAME")
		for _, f := range ls.Data {
			fmt.Printf("%s\t%s\t%d\t%s\n", f.CID, f.ID, f.Size, f.Name)
		}

	case *deals != "":
		ds, err := cli.Deals().Status(ctx, *deals)
		if err != nil {
			log.Fatal(err)
		}
		
		if len(ds) == 0 {
			fmt.Println("No deals found.")
			return
		}
		
		fmt.Printf("Found %d deal(s):\n\n", len(ds))
		fmt.Printf("%-30s %-20s %-50s %s\n", "Provider", "Status", "PieceCID", "ChainDealID")
		fmt.Println(strings.Repeat("-", 120))
		
		for _, d := range ds {
			provider := d.StorageProvider
			if provider == "" {
				provider = "(empty)"
			}
			status := d.DealStatus
			if status == "" {
				status = "(empty)"
			}
			pieceCID := d.PieceCID
			if pieceCID == "" {
				pieceCID = "(empty)"
			}
			
			fmt.Printf("%-30s %-20s %-50s %d\n", provider, status, pieceCID, d.ChainDealID)
		}
	case *del != "":
		if err := cli.Files().Delete(ctx, *del); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Delete request completed.")

	default:
		usage()
	}
}

func usage() {
	fmt.Println(`Usage:
  lhctl --upload <path>                 Upload a file (shows progress)
  lhctl --info <cid>                    Fetch file info by CID
  lhctl --list [--last-key <cursor>]    List uploaded files (shows IDs)
  lhctl --deals <cid>                   Show Filecoin deal status for a CID
  lhctl --delete <id>                   Delete a file by ID (from --list)

Environment:
  LIGHTHOUSE_API_KEY  API key for authenticated endpoints`)
}