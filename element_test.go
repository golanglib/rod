package rod_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/ysmood/kit"
)

func (s *S) TestClick() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("button")
	el.MustClick()

	s.True(p.MustHas("[a=ok]"))

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustClick()
	})
	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustClick()
	})
	s.Panics(func() {
		defer s.errorAt(3, nil)()
		el.MustClick()
	})
	s.Panics(func() {
		defer s.errorAt(4, nil)()
		el.MustClick()
	})
	s.Panics(func() {
		defer s.errorAt(5, nil)()
		el.MustClick()
	})
}

func (s *S) TestClickable() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	s.True(p.MustElement("button").MustClickable())
}

func (s *S) TestNotClickable() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("button")

	// cover the button with a green div
	p.MustWaitLoad().MustEval(`() => {
		let div = document.createElement('div')
		div.style = 'position: absolute; left: 0; top: 0; width: 500px; height: 500px;'
		document.body.append(div)
	}`)
	s.Panics(func() {
		el.MustClick()
	})

	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustClickable()
	})
	s.Panics(func() {
		defer s.errorAt(4, nil)()
		el.MustClickable()
	})
	s.Panics(func() {
		defer s.errorAt(8, nil)()
		el.MustClickable()
	})
}

func (s *S) TestHover() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("button")
	el.MustEval(`this.onmouseenter = () => this.dataset['a'] = 1`)
	el.MustHover()
	s.Equal("1", el.MustEval(`this.dataset['a']`).String())
}

func (s *S) TestElementContext() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("button").Timeout(time.Minute).CancelTimeout()
	s.Error(el.Click(proto.InputMouseButtonLeft))
}

func (s *S) TestIframes() {
	p := s.page.MustNavigate(srcFile("fixtures/click-iframes.html"))
	frame := p.MustElement("iframe").Frame().MustElement("iframe").Frame()
	el := frame.MustElement("button")
	el.MustClick()
	s.True(frame.MustHas("[a=ok]"))

	id := el.MustNodeID()
	s.Panics(func() {
		defer s.errorAt(2, nil)()
		p.MustElementFromNode(id)
	})
	s.Panics(func() {
		defer s.at(4, func(d []byte, err error) ([]byte, error) {
			return sjson.SetBytes(d, "result", rod.Array{})
		})()
		p.MustElementFromNode(id).MustText()
	})
	s.Panics(func() {
		defer s.errorAt(7, nil)()
		p.MustElementFromNode(id)
	})
	s.Panics(func() {
		defer s.errorAt(12, nil)()
		p.MustElementFromNode(id)
	})
	s.Panics(func() {
		defer s.errorAt(16, nil)()
		p.MustElementFromNode(id)
	})
}

func (s *S) TestContains() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	a := p.MustElement("button")

	b := p.MustElementFromNode(a.MustNodeID())
	s.True(a.MustContainsElement(b))

	box := a.MustBox()
	c := p.MustElementFromPoint(int(box.X)+3, int(box.Y)+3)
	s.True(a.MustContainsElement(c))

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		a.MustContainsElement(b)
	})
}

func (s *S) TestShadowDOM() {
	p := s.page.MustNavigate(srcFile("fixtures/shadow-dom.html")).MustWaitLoad()
	el := p.MustElement("#container")
	s.Equal("inside", el.MustShadowRoot().MustElement("p").MustText())

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustShadowRoot()
	})
	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustShadowRoot()
	})
}

func (s *S) TestPress() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("[type=text]")
	el.MustPress('A')
	el.MustPress(' ')
	el.MustPress('b')

	s.Equal("A b", el.MustText())

	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustPress(' ')
	})
	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustSelectAllText()
	})
}

func (s *S) TestKeyDown() {
	p := s.page.MustNavigate(srcFile("fixtures/keys.html"))
	p.MustElement("body")
	p.Keyboard.MustDown('j')

	s.True(p.MustHas("body[event=key-down-j]"))
}

func (s *S) TestKeyUp() {
	p := s.page.MustNavigate(srcFile("fixtures/keys.html"))
	p.MustElement("body")
	p.Keyboard.MustUp('x')

	s.True(p.MustHas("body[event=key-up-x]"))
}

func (s *S) TestText() {
	text := "雲の上は\nいつも晴れ"

	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("textarea")
	el.MustInput(text)

	s.Equal(text, el.MustText())
	s.True(p.MustHas("[event=textarea-change]"))

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustText()
	})
}

func (s *S) TestCheckbox() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("[type=checkbox]")
	s.True(el.MustClick().MustProperty("checked").Bool())
}

func (s *S) TestSelectText() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("textarea")
	el.MustInput("test")
	el.MustSelectAllText()
	el.MustInput("test")
	s.Equal("test", el.MustText())

	el.MustSelectText(`es`)
	el.MustInput("__")

	s.Equal("t__t", el.MustText())

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustSelectText("")
	})
	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustSelectAllText()
	})

	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustInput("")
	})
	s.Panics(func() {
		defer s.errorAt(4, nil)()
		el.MustInput("")
	})
}

func (s *S) TestBlur() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("#blur").MustInput("test").MustBlur()

	s.Equal("ok", *el.MustAttribute("a"))
}

func (s *S) TestSelectOptions() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("select")
	el.MustSelect("B", "C")

	s.Equal("B,C", el.MustText())
	s.EqualValues(1, el.MustProperty("selectedIndex").Int())
}

func (s *S) TestMatches() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("textarea")
	s.True(el.MustMatches(`[cols="30"]`))

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustMatches("")
	})
}

func (s *S) TestAttribute() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("textarea")
	cols := el.MustAttribute("cols")
	rows := el.MustAttribute("rows")

	s.Equal("30", *cols)
	s.Equal("10", *rows)

	p = s.page.MustNavigate(srcFile("fixtures/click.html"))
	el = p.MustElement("button").MustClick()

	s.Equal("ok", *el.MustAttribute("a"))
	s.Nil(el.MustAttribute("b"))

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustAttribute("")
	})
}

func (s *S) TestProperty() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("textarea")
	cols := el.MustProperty("cols")
	rows := el.MustProperty("rows")

	s.Equal(float64(30), cols.Num)
	s.Equal(float64(10), rows.Num)

	p = s.page.MustNavigate(srcFile("fixtures/open-page.html"))
	el = p.MustElement("a")

	s.Equal("link", el.MustProperty("id").Str)
	s.Equal("_blank", el.MustProperty("target").Str)
	s.Equal(gjson.Null, el.MustProperty("test").Type)

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustProperty("")
	})
}

func (s *S) TestSetFiles() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement(`[type=file]`)
	el.MustSetFiles(
		slash("fixtures/click.html"),
		slash("fixtures/alert.html"),
	)

	list := el.MustEval("Array.from(this.files).map(f => f.name)").Array()
	s.Len(list, 2)
	s.Equal("alert.html", list[1].String())
}

func (s *S) TestSelectQuery() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("select")
	el.MustSelect("[value=c]")

	s.EqualValues(2, el.MustEval("this.selectedIndex").Int())
}

func (s *S) TestSelectQueryNum() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("select")
	el.MustSelect("123")

	s.EqualValues(-1, el.MustEval("this.selectedIndex").Int())
}

func (s *S) TestEnter() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("[type=submit]")
	el.MustPress(input.Enter)

	s.True(p.MustHas("[event=submit]"))
}

func (s *S) TestWaitInvisible() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	h4 := p.MustElement("h4")
	btn := p.MustElement("button")
	timeout := 3 * time.Second

	s.True(h4.MustVisible())

	h4t := h4.Timeout(timeout)
	h4t.MustWaitVisible()
	h4t.CancelTimeout()

	go func() {
		kit.Sleep(0.03)
		h4.MustEval(`this.remove()`)
		kit.Sleep(0.03)
		btn.MustEval(`this.style.visibility = 'hidden'`)
	}()

	h4.Timeout(timeout).MustWaitInvisible()
	btn.Timeout(timeout).MustWaitInvisible()

	s.False(p.MustHas("h4"))
}

func (s *S) TestWaitStable() {
	p := s.page.MustNavigate(srcFile("fixtures/wait-stable.html"))
	el := p.MustElement("button")
	el.MustWaitStable()
	el.MustClick()
	p.MustHas("[event=click]")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		kit.Sleep(0.2)
		cancel()
	}()
	s.Error(el.Context(ctx, cancel).WaitStable(time.Minute))
}

func (s *S) TestCanvasToImage() {
	p := s.page.MustNavigate(srcFile("fixtures/canvas.html"))
	src, err := png.Decode(bytes.NewBuffer(p.MustElement("#canvas").MustCanvasToImage("", 1.0)))
	utils.E(err)
	s.Equal(src.At(50, 50), color.NRGBA{0xFF, 0x00, 0x00, 0xFF})
}

func (s *S) TestResource() {
	p := s.page.MustNavigate(srcFile("fixtures/resource.html"))
	el := p.MustElement("img").MustWaitLoad()
	s.Equal(15456, len(el.MustResource()))

	func() {
		defer s.at(3, func(res []byte, err error) ([]byte, error) {
			return kit.MustToJSONBytes(proto.PageGetResourceContentResult{
				Content:       "ok",
				Base64Encoded: false,
			}), nil
		})()
		s.Equal([]byte("ok"), el.MustResource())
	}()

	s.Panics(func() {
		defer s.errorAt(2, nil)()
		el.MustResource()
	})
	s.Panics(func() {
		defer s.errorAt(3, nil)()
		el.MustResource()
	})
}

func (s *S) TestElementScreenshot() {
	f := filepath.Join("tmp", kit.RandString(8)+".png")
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("h4")

	data := el.MustScreenshot(f)
	img, err := png.Decode(bytes.NewBuffer(data))
	utils.E(err)
	s.EqualValues(200, img.Bounds().Dx())
	s.EqualValues(30, img.Bounds().Dy())
	s.FileExists(f)

	s.Panics(func() {
		defer s.errorAt(1, nil)()
		el.MustScreenshot()
	})
	s.Panics(func() {
		s.countCall()
		defer s.errorAt(2, nil)()
		el.MustScreenshot()
	})
	s.Panics(func() {
		defer s.errorAt(3, nil)()
		el.MustScreenshot()
	})
}

func (s *S) TestUseReleasedElement() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	btn := p.MustElement("button")
	btn.MustRelease()
	s.EqualError(btn.Click("left"), "context canceled")

	btn = p.MustElement("button")
	utils.E(proto.RuntimeReleaseObject{ObjectID: btn.ObjectID}.Call(p))
	s.EqualError(btn.Click("left"), "{\"code\":-32000,\"message\":\"Could not find object with given id\",\"data\":\"\"}")
}

func (s *S) TestElementMultipleTimes() {
	// To see whether chrome will reuse the remote object ID or not.
	// Seems like it will not.

	page := s.page.MustNavigate(srcFile("fixtures/click.html"))

	btn01 := page.MustElement("button")
	btn02 := page.MustElement("button")

	s.Equal(btn01.MustText(), btn02.MustText())
	s.NotEqual(btn01.ObjectID, btn02.ObjectID)
}

func (s *S) TestFnErr() {
	p := s.page.MustNavigate(srcFile("fixtures/click.html"))
	el := p.MustElement("button")

	_, err := el.Eval(true, "foo()", nil)
	s.Error(err)
	s.Contains(err.Error(), "ReferenceError: foo is not defined")
	s.True(errors.Is(err, rod.ErrEval))

	_, err = el.ElementByJS("foo()", nil)
	s.Error(err)
	s.Contains(err.Error(), "ReferenceError: foo is not defined")
	s.True(errors.Is(err, rod.ErrEval))
}

func (s *S) TestElementEWithDepth() {
	checkStr := `green tea`
	p := s.page.MustNavigate(srcFile("fixtures/describe.html"))

	ulDOMNode, err := p.MustElement(`ul`).Describe(-1, true)
	s.Nil(errors.Unwrap(err))

	data, err := json.Marshal(ulDOMNode)
	s.Nil(errors.Unwrap(err))
	// The depth is -1, should contain checkStr
	s.Contains(string(data), checkStr)
}

func (s *S) TestElementOthers() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("form")
	s.IsType(p.GetContext(), el.GetContext())
	el.MustFocus()
	el.MustScrollIntoView()
	s.EqualValues(784, el.MustBox().Width)
	s.Equal("submit", el.MustElement("[type=submit]").MustText())
	s.Equal("<input type=\"submit\" value=\"submit\">", el.MustElement("[type=submit]").MustHTML())
	el.MustWait(`true`)
	s.Equal("form", el.MustElementByJS(`this`).MustDescribe().LocalName)
	s.Len(el.MustElementsByJS(`[]`), 0)
}

func (s *S) TestElementErrors() {
	p := s.page.MustNavigate(srcFile("fixtures/input.html"))
	el := p.MustElement("form")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := el.Context(ctx, cancel).Describe(-1, true)
	s.Error(err)

	err = el.Context(ctx, cancel).Focus()
	s.Error(err)

	err = el.Context(ctx, cancel).Press('a')
	s.Error(err)

	err = el.Context(ctx, cancel).Input("a")
	s.Error(err)

	err = el.Context(ctx, cancel).Select([]string{"a"})
	s.Error(err)

	err = el.Context(ctx, cancel).WaitStable(0)
	s.Error(err)

	_, err = el.Context(ctx, cancel).Box()
	s.Error(err)

	_, err = el.Context(ctx, cancel).Resource()
	s.Error(err)

	err = el.Context(ctx, cancel).Input("a")
	s.Error(err)

	err = el.Context(ctx, cancel).Input("a")
	s.Error(err)

	_, err = el.Context(ctx, cancel).HTML()
	s.Error(err)

	_, err = el.Context(ctx, cancel).Visible()
	s.Error(err)

	_, err = el.Context(ctx, cancel).CanvasToImage("", 0)
	s.Error(err)

	err = el.Context(ctx, cancel).Release()
	s.Error(err)

	s.Panics(func() {
		s.countCall()
		defer s.errorAt(2, nil)()
		el.MustNodeID()
	})
}
