package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                             *Application
	backgroundMusic                 *twodee.Music
	planetDropEffect                *twodee.SoundEffect
	planetFireDeathEffect           *twodee.SoundEffect
	planetCollisionEffect           *twodee.SoundEffect
	backgroundMusicObserverId       int
	planetDropEffectObserverId      int
	planetFireDeathEffectObserverId int
	planetCollisionEffectObserverId int
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	a.backgroundMusic.Play(-1)
}

func (a *AudioSystem) PlayPlanetDropEffect(e twodee.GETyper) {
	a.planetDropEffect.PlayChannel(2, 1)
}

func (a *AudioSystem) PlayPlanetFireDeathEffect(e twodee.GETyper) {
	a.planetFireDeathEffect.PlayChannel(3, 1)
}

func (a *AudioSystem) PlayPlanetCollisionEffect(e twodee.GETyper) {
	a.planetCollisionEffect.PlayChannel(4, 1)
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.backgroundMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(DropPlanet, a.planetDropEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(DropPlanet, a.planetFireDeathEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(DropPlanet, a.planetCollisionEffectObserverId)
	a.backgroundMusic.Delete()
	a.planetDropEffect.Delete()
	a.planetFireDeathEffect.Delete()
	a.planetCollisionEffect.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		backgroundMusic       *twodee.Music
		planetDropEffect      *twodee.SoundEffect
		planetFireDeathEffect *twodee.SoundEffect
		planetCollisionEffect *twodee.SoundEffect
	)
	if backgroundMusic, err = twodee.NewMusic("assets/music/Birth_of_a_Phantom_Planet.ogg"); err != nil {
		return
	}
	if planetDropEffect, err = twodee.NewSoundEffect("assets/sound_effects/PlanetDrop.ogg"); err != nil {
		return
	}
	if planetFireDeathEffect, err = twodee.NewSoundEffect("assets/sound_effects/PlanetFireDeath.ogg"); err != nil {
		return
	}
	if planetCollisionEffect, err = twodee.NewSoundEffect("assets/sound_effects/PlanetCollision.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:                   app,
		backgroundMusic:       backgroundMusic,
		planetDropEffect:      planetDropEffect,
		planetFireDeathEffect: planetFireDeathEffect,
		planetCollisionEffect: planetCollisionEffect,
	}
	planetDropEffect.SetVolume(60)
	audioSystem.backgroundMusicObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.planetDropEffectObserverId = app.GameEventHandler.AddObserver(DropPlanet, audioSystem.PlayPlanetDropEffect)
	audioSystem.planetFireDeathEffectObserverId = app.GameEventHandler.AddObserver(DropPlanet, audioSystem.PlayPlanetFireDeathEffect)
	audioSystem.planetCollisionEffectObserverId = app.GameEventHandler.AddObserver(DropPlanet, audioSystem.PlayPlanetCollisionEffect)
	return
}
