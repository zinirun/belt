<div align="center">
<h1>BELT</h1>
<p>Super-fast and easy data transforming with any interfaces</p>
<p>Just connect Belt to your Factory</p>
</div>

## Methods
- Transform data
  - Map
  - Filter
- Assistance
  - Pipe
  - Append
- Load data
  - FromHttp
  - FromQuery
- Initialize and build
  - NewFactory
  - Build
- Debug
  - ListItems
  - Type

## Getting Started
- Initialize Factory with sync option
    ```go
    // true: with go sync, false: without go sync
	option := FactoryOptions{true}
	items = NewFactory(items, option)
    ```
    - or use `FromHttp` to load data from http request
        ```go
        type Commit struct {
            Commits  int    `json:"commits"`
            Date     string `json:"date"`
            Day      int    `json:"day"`
            Month    int    `json:"month"`
            Week     int    `json:"week"`
            LastWeek bool   `json:"lastWeek"`
            LastDay  bool   `json:"lastDay"`
        }
        func example() {
            url := "https://gw.alipayobjects.com/os/antvdemo/assets/data/github-commit.json"
            output := Commit{}
            items := NewFactory([]I{}, FactoryOptions{false}).FromHttp(url, output)
        }
        ```
    - or use `FromQuery` to load data from MySQL database
        ```go
        type User struct {
            id          int
            email       string
        }
        func example(){
            dbSource := "USER/PASSWORD@tcp(HOST:PORT)/DATABASE"
            query := "select id, email from user where isDeleted = ?"
            output := User{}
            items := NewFactory([]I{}, FactoryOptions{false}).FromQuery(dbSource, output, query, 0)
        }
        ```
- Transform data and build result
  - create functions for transforming data
    ```go
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
    ```
  - connect Belts to Factory and build
    ```go
    func example(){
        items = NewFactory(items, option).Map(map1).Map(map2).Filter(filter1).Append(newChild).Pipe(printLen).Build()
    }
    ```

## License
© Zini, Released under the MIT License