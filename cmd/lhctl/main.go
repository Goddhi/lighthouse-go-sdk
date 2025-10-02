package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

func progressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	
	filled := (percent * width) / 100
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	
	bar := strings.Repeat(".", filled) + strings.Repeat(",", width-filled)
	return fmt.Sprintf("[%s]", bar)
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
		startTime := time.Now()
		var lastPercent float64

		res, err := cli.Storage().UploadFile(ctx, *upload, schema.WithProgress(func(p schema.Progress) {
			percent := p.Percent()
			if percent-lastPercent >= 1.0 || percent >= 100.0{
				elapsed := time.Since(startTime).Seconds()
				speed := float64(p.Uploaded) / elapsed / 1024 / 1024

				bar := progressBar(int(percent), 40)
				fmt.Printf("\r%s %.1f%% (%d/%d bytes) %.2f MB/s",
					bar, percent, p.Uploaded, p.Total, speed)
				lastPercent = percent
			}
			}))

			fmt.Println()

			if err != nil {
			log.Fatal(err)
			}

			fmt.Print("Upload complete!\n")
			fmt.Printf("CID %s\n", res.Hash)
			fmt.Printf("Time: %.2fs\n", time.Since(startTime).Seconds())

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