# bsky-client
A Golang client for Bluesky that simplifies some interfaces

## Usage

### Authentication

Create a `.env` file in the root of this repo and put the following variables in it:
```
ATPROTO_HANDLE={my_handle.bsky.social}
ATPROTO_APP_PASSWORD={an_app_password_for_my_account}
```

An example CLI can be found in `cmd/cli` that allows you to create a post with tags, images, alt-text, quoting an existing post, and replying to an existing post.

```
go run cmd/cli/main.go post --text "testing a reply" --images test_cat_image.jpg --image-alt-texts "a cat licking its own nose sitting on a table" --tags cat,kitten,nose,mlem,eng --reply-to at://did:plc:o7eggiasag3efvoc3o7qo3i3/app.bsky.feed.post/3kavby7fn2u2j
```
