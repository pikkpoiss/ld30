package main

import (
	"image/color"
	"time"

	twodee "../libs/twodee"
)

const (
	programCode int32 = iota
	settingCode
	exitCode
	musicCode
)

var (
	regColor = color.RGBA{200, 200, 200, 255}
	hiColor  = color.RGBA{255, 240, 120, 255}
	actColor = color.RGBA{200, 200, 255, 255}
)

type MenuLayer struct {
	visible  bool
	menu     *twodee.Menu
	text     *twodee.TextRenderer
	regFont  *twodee.FontFace
	hiFont   *twodee.FontFace
	actFont  *twodee.FontFace
	cache    map[int]*twodee.TextCache
	hiCache  *twodee.TextCache
	actCache *twodee.TextCache
	bounds   twodee.Rectangle
	app      *Application
}

func NewMenuLayer(app *Application) (layer *MenuLayer, err error) {
	var (
		menu    *twodee.Menu
		text    *twodee.TextRenderer
		regFont *twodee.FontFace
		hiFont  *twodee.FontFace
		actFont *twodee.FontFace
		bg      = color.Transparent
		font    = "assets/fonts/Exo-SemiBold.ttf"
	)
	if regFont, err = twodee.NewFontFace(font, 32, regColor, bg); err != nil {
		return
	}
	if hiFont, err = twodee.NewFontFace(font, 32, hiColor, bg); err != nil {
		return
	}
	if actFont, err = twodee.NewFontFace(font, 32, actColor, bg); err != nil {
		return
	}
	if text, err = twodee.NewTextRenderer(app.WinBounds); err != nil {
		return
	}
	menu, err = twodee.NewMenu([]twodee.MenuItem{
		twodee.NewKeyValueMenuItem("Music On/Off", programCode, musicCode),
		twodee.NewKeyValueMenuItem("Exit", programCode, exitCode),
	})
	if err != nil {
		return
	}
	layer = &MenuLayer{
		visible:  false,
		menu:     menu,
		text:     text,
		regFont:  regFont,
		hiFont:   hiFont,
		actFont:  actFont,
		cache:    map[int]*twodee.TextCache{},
		hiCache:  twodee.NewTextCache(hiFont),
		actCache: twodee.NewTextCache(actFont),
		bounds:   app.WinBounds,
		app:      app,
	}
	return

}

func (l *MenuLayer) HandleEvent(evt twodee.Event) bool {
	// Handle the !visible case quickly.
	if !l.visible {
		switch event := evt.(type) {
		case *twodee.KeyEvent:
			if event.Type != twodee.Press {
				break
			}
			if event.Code == twodee.KeyEscape {
				if l.visible {
					l.visible = false
					return false
				}
				l.menu.Reset()
				l.visible = true
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuOpen))
				return false
			}
		}
		return true
	}
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.visible = false
			l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuClose))
			return false
		case twodee.KeyUp:
			l.menu.Prev()
			l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuClick))
			return false
		case twodee.KeyDown:
			l.menu.Next()
			l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuClick))
			return false
		case twodee.KeyEnter:
			if data := l.menu.Select(); data != nil {
				l.handleMenuItem(data)
			}
			l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuSel))
			return false
		}
	}
	return true
}

func (l *MenuLayer) handleMenuItem(data *twodee.MenuItemData) {
	switch data.Key {
	case programCode:
		switch data.Value {
		case musicCode:
			// TODO: Write code that mutes/un-mutes music.
		case exitCode:
			l.app.InitiateCloseGame = true
		}
	}
}

func (l *MenuLayer) Update(elapsed time.Duration) {}

func (l *MenuLayer) Reset() (err error) {
	if l.text != nil {
		l.text.Delete()
	}
	if l.text, err = twodee.NewTextRenderer(l.bounds); err != nil {
		return
	}
	l.actCache.Clear()
	l.hiCache.Clear()
	for _, v := range l.cache {
		v.Clear()
	}
	return
}

func (l *MenuLayer) Delete() {
	l.text.Delete()
	l.actCache.Delete()
	l.hiCache.Delete()
	for _, v := range l.cache {
		v.Clear()
	}
}

func (l *MenuLayer) Render() {
	if !l.visible {
		return
	}
	var (
		textCache *twodee.TextCache
		texture   *twodee.Texture
		ok        bool
		y         = l.bounds.Max.Y
	)
	l.text.Bind()
	for i, item := range l.menu.Items() {
		if item.Highlighted() {
			l.hiCache.SetText(item.Label())
			texture = l.hiCache.Texture
		} else if item.Active() {
			l.actCache.SetText(item.Label())
			texture = l.actCache.Texture
		} else {
			if textCache, ok = l.cache[i]; !ok {
				textCache = twodee.NewTextCache(l.regFont)
				l.cache[i] = textCache
			}
			textCache.SetText(item.Label())
			texture = textCache.Texture
		}
		if texture != nil {
			y = y - float32(texture.Height)
			l.text.Draw(texture, 0, y)

		}
	}
	l.text.Unbind()
}
