## Supported Go versions

## Feature Overview

- Wrapper of Notion API from Developers Beta(https://developers.notion.com/)
- Support Databases API
- Support Pages API
- Support Blocks API
- Support Users API
- Support Search API

## Guide

### Installation

> go get github.com/hunydev/notion

### Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hunydev/notion"
	api "github.com/hunydev/notion/api/v20210513"
)

func main() {
	//token := "secret_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	token := os.Getenv("NOTION_AUTHORIZATION")

	nt := notion.New(api.New(token, &api.Option{
		Timeout: time.Second * 3,
	}))

	fmt.Println(nt.APIVersion())
}
```

#### Databases
*List Databases*
```go
//List Databases
resp, err := nt.ListDatabases(&notion.PaginationRequest{}) // can use nil pointer
if err != nil {
    panic(err)
}

databases, err := resp.Databases()
if err != nil {
    panic(err)
}

for _, database := range databases {
    fmt.Println("ID:", database.ID())
}
```


## License
MIT