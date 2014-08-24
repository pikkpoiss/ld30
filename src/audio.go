package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                        *Application
	backgroundMusic            *twodee.Music
	planetDropEffect           *twodee.SoundEffect
	backgroundMusicObserverId  int
	planetDropEffectObserverId int
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	a.backgroundMusic.Play(-1)
}

func (a *AudioSystem) PlayPlanetDropEffect(e twodee.GETyper) {
	a.planetDropEffect.PlayChannel(2, 1)
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.backgroundMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(DropPlanet, a.planetDropEffectObserverId)
	a.backgroundMusic.Delete()
	a.planetDropEffect.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		backgroundMusic  *twodee.Music
		planetDropEffect *twodee.SoundEffect
	)
	if backgroundMusic, err = twodee.NewMusic("assets/music/BGM1.ogg"); err != nil {
		return
	}
	if planetDropEffect, err = twodee.NewSoundEffect("assets/sound_effects/PlanetDrop.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:              app,
		backgroundMusic:  backgroundMusic,
		planetDropEffect: planetDropEffect,
	}
	planetDropEffect.SetVolume(60)
	audioSystem.backgroundMusicObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.planetDropEffectObserverId = app.GameEventHandler.AddObserver(DropPlanet, audioSystem.PlayPlanetDropEffect)
	return
}
