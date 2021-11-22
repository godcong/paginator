# paginator

## Paginator is a go-based page-turning extension package

## Usage

```go
p:= paginator.New() //create the paginator module for use
//you can create a paginator with a custom options
p:= paginator.New(PerPageOption(30)) //paging 30 items per page
p:= paginator.New(PerPageKeyOption("perPage")) //set per page key with option
p:= paginator.New(PageKeyOption("page")) //set page key with option
p:= paginator.New(PerPageOption(30)，PerPageKeyOption("perPage"),PageKeyOption("page")) //create the paginator module with all custom options

//use the paginator
page,err:=p.Parse(Turnable) //parse will return the current page data and error

//Turnable at least 3 interfaces that need to be implemented
Counter   //count the total data
Finder    //find the data by page
Requester //get the http request from your set

// If you need to customize the query you need to implement the following 2 interfaces
Hooker //auto call the hook function when the page have setted value
CustomHooker //parse customize the query from request all you want todo
InitHooker //do something before the hook function do

type pageExample struct {
ctx     context.Context
model   *Model
request *http.Request
query   *Query
}

func (p pageExample) Count(paginator.Iterator) (int64, error) {
count, err := p.Query().Count(p.ctx)
return int64(count), err
}

func (p pageExample) Find(page paginator.PrePager) (interface{}, error) {
return p.Query().Limit(page.Limit()).Offset(page.Offset()).All(p.ctx)
}

func (p pageExample) Request() *http.Request {
return p.request
}

func (p *pageExample) hookUpdate(v []string) bool {
id:=v[0]
if err == nil {
p.query = p.query.Where(page.IDEq(id))
}
return true
}

func (p *pageExample) Hooks() map[string]func([]string) bool {
return map[string]func([]string) bool{
"id": p.hookID,
}
}

func (p *pageExample) Initialize() {
p.query = p.model.Query()
}

func (p *pageExample) Query() *Query {
return p.query.Clone
}

page,err:=p.Parse(&pageExample{
    //set init data for the paginator
})



```

