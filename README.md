# echo-binder
echo-binder 一个提供echo中数据binder和validator功能的middleware

##TODO
 - 添加gin的标签  
 - 完善注入方式

## Update
- 20161120 升级到支持echo v2 版本(可以通过tag切换版本 master上为version2.0)
- 20161018 使用[bluemonday](github.com/microcosm-cc/bluemonday)添加xss过滤,使用方式详见test TestXssBinder_Bind
- 20170922 添加form的类型解析

## Quick Start

### Installation
```
$ go get -u github.com/cnjack/echo-binder
```
### Hello, World!
```go
package main

import (
	"github.com/cnjack/echo-binder"
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" binding:"required"`
	Age   int    `json:"age" xml:"age" form:"age" binding:"gte=0,lte=130"`
	Email string `json:"email" xml:"email" form:"email" binding:"required,email"`
}

func main() {
	e := echo.New()
	e.Use(binder.BindBinder(e))
	e.POST("/", func(c echo.Context) error {
		var u User
		if err := c.Bind(&u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, "Hello, "+u.Name)
	})
	e.GET("/", func(c echo.Context) error {
		var u User
		if err := c.Bind(&u); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, "Hello, "+u.Name)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

```

## form interface EXAMPLE
```go
type Price int64

func (p Price) String() string {
	return strconv.Itoa(int(p))
}

func (p *Price) UnmarshalForm(text string) error {
	if text == "" {
		return errors.New("invalid Money text, it should not be empty")
	}
	if dotIndex := strings.IndexByte(text, '.'); dotIndex >= 0 {
		dotFront := text[:dotIndex]
		dotBehind := text[dotIndex+1:]
		switch len(dotBehind) {
		default:
			return fmt.Errorf("invalid Money text: %s", text)
		case 0:
			text = dotFront + "00"
		case 1:
			text = dotFront + dotBehind + "0"
		case 2:
			text = dotFront + dotBehind
		}
	} else {
		text += "00"
	}
	n, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Money text: %s", text)
	}
	*p = Price(n)
	return nil
}
```

## Thx.
[echo](https://github.com/labstack/echo) Fast and unfancy HTTP server framework for Go (Golang)  
[assert](https://github.com/stretchr/testify) A sacred extension to the standard go testing package  
[validator.v9](https://gopkg.in/go-playground/validator.v9) Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving  
[bluemonday](https://github.com/microcosm-cc/bluemonday) a fast golang HTML sanitizer (inspired by the OWASP Java HTML Sanitizer) to scrub user generated content of XSS  
