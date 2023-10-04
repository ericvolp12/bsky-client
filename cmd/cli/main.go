package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/ericvolp12/bsky-client/pkg/client"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:    "bsky-client",
		Usage:   "bluesky cli client",
		Version: "0.0.1",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "pds-url",
			Usage:   "http(s) url of the pds server",
			Value:   "https://bsky.social",
			EnvVars: []string{"PDS_URL"},
		},
		&cli.StringFlag{
			Name:    "plc-url",
			Usage:   "http(s) url of the plc server",
			Value:   "https://plc.directory",
			EnvVars: []string{"PLC_URL"},
		},
		&cli.StringFlag{
			Name:    "handle",
			Usage:   "handle of the user to use for authentication",
			EnvVars: []string{"ATPROTO_HANDLE"},
		},
		&cli.StringFlag{
			Name:    "app-password",
			Usage:   "app-password of the user to use for authentication",
			EnvVars: []string{"ATPROTO_APP_PASSWORD"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name: "post",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "text",
					Usage: "text of the post",
				},
				&cli.StringSliceFlag{
					Name:  "images",
					Usage: "list of paths to images to attach to the post",
				},
				&cli.StringSliceFlag{
					Name:  "image-alt-texts",
					Usage: "list of alt texts to attach to the images",
				},
				&cli.StringSliceFlag{
					Name:  "tags",
					Usage: "list of tags to attach to the post",
				},
				&cli.StringFlag{
					Name:  "reply-to",
					Usage: "uri of the post to reply to",
				},
				&cli.StringFlag{
					Name:  "quoting",
					Usage: "uri of the post to quote",
				},
				&cli.StringFlag{
					Name:  "embedded-link",
					Usage: "url of an embedded link to attach to the post",
				},
			},
			Action: Post,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Post(cctx *cli.Context) error {
	// Initialize Post Args
	postArgs := client.PostArgs{
		Text:         cctx.String("text"),
		Tags:         cctx.StringSlice("tags"),
		EmbeddedLink: cctx.String("embedded-link"),
	}

	if cctx.StringSlice("images") != nil {
		postArgs.Images = []client.Image{}
		for i, path := range cctx.StringSlice("images") {
			fi, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to stat image: %w", err)
			}
			if fi.IsDir() {
				return fmt.Errorf("image path is a directory")
			}
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open image: %w", err)
			}
			defer f.Close()
			buf := make([]byte, fi.Size())
			_, err = f.Read(buf)
			if err != nil {
				return fmt.Errorf("failed to read image: %w", err)
			}
			altText := ""
			if len(cctx.StringSlice("image-alt-texts")) > i {
				altText = cctx.StringSlice("image-alt-texts")[i]
			}
			postArgs.Images = append(postArgs.Images, client.Image{
				Reader:  bytes.NewReader(buf),
				AltText: altText,
			})
		}
	}

	if cctx.String("reply-to") != "" {
		replyingTo, err := syntax.ParseATURI(cctx.String("reply-to"))
		if err != nil {
			return fmt.Errorf("failed to parse reply-to uri: %w", err)
		}
		postArgs.ReplyTo = &replyingTo
	}

	if cctx.String("quoting") != "" {
		quoting, err := syntax.ParseATURI(cctx.String("quoting"))
		if err != nil {
			return fmt.Errorf("failed to parse quoting uri: %w", err)
		}
		postArgs.Quoting = &quoting
	}

	// Initialize Client
	c := client.New(cctx.String("pds-url"), cctx.String("plc-url"))

	// Authenticate
	err := c.Login(cctx.Context, cctx.String("handle"), cctx.String("app-password"))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// Create Post
	uri, err := c.CreatePost(cctx.Context, postArgs)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	fmt.Printf("created post: %s\n", uri.String())

	return nil
}
