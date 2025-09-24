package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("LIGHTHOUSE_API_KEY")
	if apiKey == "" {
		log.Fatal("set LIGHTHOUSE_API_KEY env var")
	}

	upload := flag.String("upload", "", "file path to upload")
	info := flag.String("info", "", "CID to fetch info for")
	list := flag.Bool("list", false, "list uploaded files")
	flag.Parse()

	cli := lighthouse.NewClient(nil, lighthouse.WithAPIKey(apiKey))

	if *upload != "" {
		res, err := cli.Storage().UploadFile(
			ctx,
			*upload,
			schema.WithProgress(func(p schema.Progress) {
				fmt.Printf("uploaded %d/%d (%.2f%%)\n", p.Uploaded, p.Total, p.Percent())
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("CID:", res.Data.Hash)
		return
	}

	if *info != "" {
		i, err := cli.Files().Info(ctx, *info)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name=%s Size=%s\n", i.Data.Name, i.Data.Size)
		return
	}

	if *list {
		ls, err := cli.Files().List(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range ls.Data {
			fmt.Printf("%s\t%s\t%d\n", f.CID, f.Name, f.Size)
		}
	}
}
