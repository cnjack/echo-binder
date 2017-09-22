package binder_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/cnjack/echo-binder"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" binding:"required"`
	Age   int    `json:"age" xml:"age" form:"age" binding:"gte=0,lte=130"`
	Email string `json:"email" xml:"email" form:"email" binding:"required,email"`
	Price Price  `form:"price"`
}

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

type Xss struct {
	Data  string `json:"data" xss:"true"`
	Image string `json:"image" xss:"true"`
}

var (
	json = `{"name": "jack","age": 25,"email": "h_7357@qq.com"}`
	xml  = `<xml><name>jack</name><age>25</age><email>h_7357@qq.com</email></xml>`
	form = `name=jack&age=25&email=h_7357@qq.com&price=0.01`
	xss  = `{"data":"<a onblur='alert(secret)' href='http://www.google.com'>Google</a>", "image":"<img src='https://ssl.gstatic.com/accounts/ui/logo_2x.png'/>"}`
)

func TestFormBinder_Bind(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(echo.GET, "/?"+form, nil)
	c := e.NewContext(req, rec)
	b := binder.NewBinder(c)
	var user User
	err := b.Bind(&user, c)
	if assert.NoError(t, err) {
		assert.Equal(t, "jack", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.Equal(t, "h_7357@qq.com", user.Email)
	}
}

func TestFormPostBinder_Bind(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(echo.POST, "/", strings.NewReader(form))
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	b := binder.NewBinder(c)
	var user = User{}
	err := b.Bind(&user, c)
	if assert.NoError(t, err) {
		assert.Equal(t, "jack", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.Equal(t, "h_7357@qq.com", user.Email)
		assert.Equal(t, "1", user.Price.String())
	}
}

func TestXmlBinder_Bind(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(xml))
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, "application/xml")
	b := binder.NewBinder(c)
	var user User
	err := b.Bind(&user, c)
	if assert.NoError(t, err) {
		assert.Equal(t, "jack", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.Equal(t, "h_7357@qq.com", user.Email)
	}
}

func TestJsonBinder_Bind(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(json))
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, "application/json")
	b := binder.NewBinder(c)
	var user User
	err := b.Bind(&user, c)
	if assert.NoError(t, err) {
		assert.Equal(t, "jack", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.Equal(t, "h_7357@qq.com", user.Email)
	}
}

func TestXssBinder_Bind(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(xss))
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, "application/json")
	b := binder.NewBinder(c)
	var x Xss
	err := b.Bind(&x, c)
	fmt.Println(x.Image)
	if assert.NoError(t, err) {
		assert.Equal(t, "<a href=\"http://www.google.com\" rel=\"nofollow\">Google</a>", x.Data)
	}
}
