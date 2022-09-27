# paginator

## Paginator is a go-based page-turning extension package

## Usage

- Create the paginator object

```go
package example

import "github.com/godcong/paginator/v3"

func main() {
	p := paginator.New() //create the paginator module for use
	//you can create a paginator with a custom options
	p := paginator.New(paginator.PerPageOption(30))           //paging 30 items per page
	p := paginator.New(paginator.PerPageKeyOption("perPage")) //set per page key with option
	p := paginator.New(paginator.PageKeyOption("page"))       //set page key with option
	p := paginator.New(paginator.PerPageOption(30),
		paginator.PerPageKeyOption("perPage"), paginator.PageKeyOption("page")) //create the paginator module with all custom options

	//use the paginator
	page, err := p.Parse(Queryable) //parse will return the current page data and error 
}

```

- Implement the Queryable interface

```
//Turnable at least 3 interfaces that need to be implemented
Queryable //return the Finder for paginator query

Counter //count the total data
Getter  //find the data by page
//optional
Cloner //all the data 
```

- A Queryable example

```go
package example

import (
	"context"

	"github.com/godcong/paginator/v3"
)

type pageExample struct {
	query *Query
}

func (p pageExample) Count(ctx context.Context) (int64, error) {
	count, err := p.query.Count(ctx)
	return int64(count), err
}

func (p pageExample) Clone() paginator.Finder {
	return p.query.Clone()
}

func (p pageExample) Finder(parser paginator.Parser) paginator.Finder {
	v := parser.FindValue("catch", "")
	if v != "" {
		p.query = p.query.Where(Cacth(v))
	}
	id := parser.FindValue("id", "")
	if id != "" {
		p.query = p.query.Where(page.IDEq(id))
	}

	return p
}

func main() {
	//then use
	p.SetDefaultQuery(&pageExample{})
	page, err := p.Parse(paginator.NewHTTPParser(req))
	//or
	page, err := p.ParseWithQuery(paginator.NewHTTPParser(req), &pageExample)
}
```

- Request from web

```
   http://127.0.0.1/api/v0/example?per_page=20&page=xx&id=xx,
```

result will like this:

```json
{
  "current_page": 1,
  "last_page": 1,
  "per_page": 20,
  "data": [
    {
      "id": "1",
      "name": "test1"
    },
    {
      "id": "2",
      "name": "test2"
    }
  ],
  "total": 2,
  "first_page_url": "127.0.0.1/api/v0/example?page=1&per_page=20",
  "last_page_url": "127.0.0.1/api/v0/example?page=1&per_page=20",
  "next_page_url": "",
  "prev_page_url": "",
  "path": "127.0.0.1/api/v0/example"
}
```