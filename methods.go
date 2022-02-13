package belt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/mitchellh/mapstructure"
)

func (f *Factory) Map(mapper Mapper) *Factory {
	if f.options.sync {
		return f.SyncMap(mapper)
	}
	for idx, item := range f.items {
		newItem, err := mapper(item, idx)
		errorHandlerWithFatal(err)
		f.items[idx] = newItem
	}
	return f
}

func (f *Factory) Filter(filter Filter) *Factory {
	if f.options.sync {
		return f.SyncFilter(filter)
	}
	items := make([]I, 0)
	for idx, item := range f.items {
		ok, err := filter(item, idx)
		errorHandlerWithFatal(err)
		if ok {
			items = append(items, item)
		}
	}
	f.items = items
	return f
}

func (f *Factory) Pipe(piper Piper) *Factory {
	if err := piper(f.items); err != nil {
		panic(err)
	}
	return f
}

func (f *Factory) Append(items ...I) *Factory {
	f.items = append(f.items, items...)
	return f
}

func (f *Factory) Build() []I {
	return f.items
}

func (f *Factory) SyncMap(mapper Mapper) *Factory {
	var wg sync.WaitGroup
	wg.Add(len(f.items))
	for idx, item := range f.items {
		go func(idx int, item I) {
			defer wg.Done()
			newItem, err := mapper(item, idx)
			errorHandlerWithPanic(err)
			f.items[idx] = newItem
		}(idx, item)
	}
	wg.Wait()
	return f
}

func (f *Factory) SyncFilter(filter Filter) *Factory {
	items := make([]I, 0)
	var wg sync.WaitGroup
	var mutex = new(sync.Mutex)
	wg.Add(len(f.items))
	for idx, item := range f.items {
		go func(idx int, item I) {
			defer wg.Done()
			defer mutex.Unlock()
			ok, err := filter(item, idx)
			errorHandlerWithPanic(err)
			mutex.Lock()
			if ok {
				items = append(items, item)
			}
		}(idx, item)
	}
	wg.Wait()
	f.items = items
	return f
}

func (f *Factory) FromQuery(dbSource string, output interface{}, query string, args ...interface{}) *Factory {
	f.items = getFromMySql(dbSource, output, query, args...)
	return f
}

func (f *Factory) FromHttp(url string, output interface{}) *Factory {
	res, err := http.Get(url)
	errorHandlerWithFatal(err)
	checkStatusCode(res)
	defer res.Body.Close()
	outputs, results := []I{}, []I{}
	json.NewDecoder(res.Body).Decode(&outputs)
	for _, item := range outputs {
		mapstructure.Decode(item.(map[string]interface{}), &output)
		results = append(results, output)
		fmt.Println(item, output)
	}
	f.items = results
	return f
}
