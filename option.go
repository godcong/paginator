package paginator

// Option is an Option for paginator
type Option struct {
	customize       bool
	usePathPaginate bool
	ignoreKeys      []Key
	noSuccess       bool
	keySet          map[Key]string
	valSet          map[Value]int
}

// SetPathKey set the path key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetPathKey(val string) *Option {
	o.SetKey(KeyPath, val)
	return o
}

// SetCurrentPageKey set the current page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetCurrentPageKey(val string) *Option {
	o.SetKey(KeyCurrentPage, val)
	return o
}

// SetTotalKey set the total key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetTotalKey(val string) *Option {
	o.SetKey(KeyTotal, val)
	return o
}

// SetLastPageKey set the last page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetLastPageKey(val string) *Option {
	o.SetKey(KeyLastPage, val)
	return o
}

// SetNextPageKey set the next page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetNextPageKey(val string) *Option {
	o.SetKey(KeyNextPageUrl, val)
	return o
}

// SetPrevPageKey set the prev page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetPrevPageKey(val string) *Option {
	o.SetKey(KeyPrevPageUrl, val)
	return o
}

// SetDataKey set the data key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetDataKey(val string) *Option {
	o.SetKey(KeyData, val)
	return o
}

// SetFirstPageKey set the first page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetFirstPageKey(val string) *Option {
	o.SetKey(KeyFirstPageUrl, val)
	return o
}

// SetPerPageKey set the per page key to paginator
// @receiver *Option
// @param string
func (o *Option) SetPerPageKey(val string) *Option {
	o.SetKey(KeyPerPage, val)
	return o
}

// SetPageKey set the page key to paginator
// @receiver *Option
// @param string
// @return *Option
func (o *Option) SetPageKey(val string) *Option {
	o.SetKey(KeyPage, val)
	return o
}

// SetStartIndex set the start index to paginator
// @receiver *Option
// @param int
// @return *Option
func (o *Option) SetStartIndex(si int) *Option {
	o.SetValue(ValueStartIndex, si)
	return o
}

// SetStartPage ...
// @receiver *Option
// @param int
// @return *Option
func (o *Option) SetStartPage(sp int) *Option {
	o.SetValue(ValueStartPage, sp)
	return o
}

// SetPerPage ...
// @receiver *Option
// @param int
// @return *Option
func (o *Option) SetPerPage(perPage int) *Option {
	o.SetValue(ValuePerPage, perPage)
	return o
}

// AddIgnore add ignore key
// @receiver *Option
// @param Key
func (o *Option) AddIgnore(key Key) *Option {
	o.ignoreKeys = append(o.ignoreKeys, key)
	return o
}

// SetIgnores set ignore keys
// @receiver *Option
// @param ...Key
func (o *Option) SetIgnores(keys ...Key) *Option {
	o.ignoreKeys = keys
	return o
}

func (o *Option) NoSuccess() *Option {
	o.noSuccess = true
	return o
}

func (o *Option) UsePathPaginate() *Option {
	o.usePathPaginate = true
	return o
}

// SettingPerPage set the per page to paginator
// @Description: per page init on paginator
// @param int
// @return Option
func SettingPerPage(perPage int) *Option {
	opts := SettingOption()
	opts.SetValue(ValuePerPage, perPage)
	return opts
}

// SettingStartIndex set the start index to paginator
// @Description: start index init on paginator
// @param int
// @return Option
func SettingStartIndex(si int) *Option {
	opts := SettingOption()
	opts.SetValue(ValueStartIndex, si)
	return opts
}

// SettingStartPage set the start page to paginator
// @Description: start page init on paginator
// @param int
// @return *Option
func SettingStartPage(sp int) *Option {
	opts := SettingOption()
	opts.SetValue(ValueStartPage, sp)
	return opts
}

func defaultOption() *Option {
	return &Option{
		customize: false,
		keySet:    make(map[Key]string),
		valSet:    make(map[Value]int),
	}
}

func SettingOption() *Option {
	return defaultOption()
}

func (o *Option) SetKey(key Key, value string) {
	o.customize = true
	if o.keySet == nil {
		o.keySet = make(map[Key]string)
	}
	o.keySet[key] = value
}

func (o *Option) SetValue(key Value, num int) {
	if o.valSet == nil {
		o.valSet = make(map[Value]int)
	}
	o.valSet[key] = num
}

func (o *Option) ClearKeys() {
	o.customize = false
	o.keySet = make(map[Key]string)
}

func (o *Option) ClearValues() {
	o.valSet = make(map[Value]int)
}

func (o *Option) Clear() {
	o.ClearKeys()
	o.ClearValues()
}
