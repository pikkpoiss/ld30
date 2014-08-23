package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                       *Application
	backgroundMusic           *twodee.Music
	backgroundMusicObserverId int
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	a.backgroundMusic.Play(-1)
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.backgroundMusicObserverId)
	a.backgroundMusic.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		backgroundMusic *twodee.Music
	)
	if backgroundMusic, err = twodee.NewMusic("assets/music/BGM1.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:             app,
		backgroundMusic: backgroundMusic,
	}
	audioSystem.backgroundMusicObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	return
}
