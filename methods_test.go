package belt

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	loop = 100
)

type Parent struct {
	id   int
	name string
}
type Child struct {
	Parent
	extra string
}

func map1(item I, idx int) (I, error) {
	p := item.(Parent)
	time.Sleep(1 * time.Millisecond)
	return Parent{
		id:   p.id + idx,
		name: p.name,
	}, nil
}

func map2(item I, idx int) (I, error) {
	p := item.(Parent)
	return Child{
		p,
		"Hello",
	}, nil
}

func filter1(item I, _ int) (bool, error) {
	i := item.(Child)
	if i.id > 15 {
		return true, nil
	}
	return false, nil
}

func printLen(items []I) error {
	fmt.Println(len(items))
	return nil
}

func TestMethodsWithNoConcurrency(t *testing.T) {
	items := make([]I, 0)
	for i := 0; i < loop; i++ {
		items = append(items, Parent{id: i, name: "Parent " + strconv.Itoa(i)})
	}
	p := Parent{loop, "New Parent"}
	newChild := Child{p, "New Child"}
	option := FactoryOptions{false}
	items = NewFactory(items, option).Map(map1).Map(map2).Filter(filter1).Append(newChild).Pipe(printLen).Build()
	// test for Append
	assert.Equal(t, items[len(items)-1], newChild)
	p = Parent{16, "Parent 8"}
	// test for Map and Filter
	expect := Child{p, "Hello"}
	assert.Equal(t, items[0], expect)
	assert.Equal(t, len(items), loop-7)
}

func TestMethodsWithSync(t *testing.T) {
	items := make([]I, 0)
	for i := 0; i < loop; i++ {
		items = append(items, Parent{id: i, name: "Parent " + strconv.Itoa(i)})
	}
	p := Parent{loop, "New Parent"}
	newChild := Child{p, "New Child"}
	option := FactoryOptions{sync: true}
	items = NewFactory(items, option).Map(map1).Map(map2).Filter(filter1).Append(newChild).Pipe(printLen).Build()
	// test for Append
	assert.Equal(t, items[len(items)-1], newChild)
	p = Parent{16, "Parent 8"}
	// test for Map and Filter
	assert.Equal(t, len(items), loop-7)
}

type transport struct {
	UUID        string
	Temperature string
}

type newTransport struct {
	UUID         string
	Temperatures []float32
}

func map3(item I, idx int) (I, error) {
	t := item.(transport)
	_temperatures := strings.Split(t.Temperature, ",")
	temperatures := make([]float32, len(_temperatures))
	for i, v := range _temperatures {
		conv, err := strconv.ParseFloat(v, 32)
		if err != nil {
			log.Print(err)
			continue
		}
		temperatures[i] = float32(conv)
	}
	return newTransport{
		UUID:         t.UUID,
		Temperatures: temperatures,
	}, nil
}

func TestFromQueryAndMap(t *testing.T) {
	dbSource := ""
	query := "select td.uuid, td.temperature from transport_done td where td.companyId = ?"
	output := transport{}
	items := NewFactory([]I{}, FactoryOptions{false}).FromQuery(dbSource, output, query, 1).Map(map3).Build()
	ListItems(items)
}

type Commit struct {
	Commits  int    `json:"commits"`
	Date     string `json:"date"`
	Day      int    `json:"day"`
	Month    int    `json:"month"`
	Week     int    `json:"week"`
	LastWeek bool   `json:"lastWeek"`
	LastDay  bool   `json:"lastDay"`
}

func TestFromHttp(t *testing.T) {
	url := "https://gw.alipayobjects.com/os/antvdemo/assets/data/github-commit.json"
	output := Commit{}
	items := NewFactory([]I{}, FactoryOptions{false}).FromHttp(url, output).Build()
	ListItems(items)
}
