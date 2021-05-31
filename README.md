## Supported Go versions
----

## Feature Overview
----

- Wrapper of Notion API from Developers Beta(https://developers.notion.com/)
- Support Databases API
- Support Pages API
- Support Blocks API
- Support Users API
- Support Search API

## Guide
----

### Installation

> go get github.com/hunydev/notionapi

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
	nt := notion.New(
		api.New(os.Getenv("NOTION_AUTHORIZATION"),
			&api.Option{Timeout: time.Second * 3}))

	fmt.Println(nt.APIVersion())
}
```

## License
----
MIT