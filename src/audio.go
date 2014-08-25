package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                             *Application
	backgroundMusic                 *twodee.Music
	planetDropEffect                *twodee.SoundEffect
	planetFireDeathEffect           *twodee.SoundEffect
	planetCollisionEffect           *twodee.SoundEffect
	backgroundMusicObserverId       int
	pauseMusicObserverId            int
	resumeMusicObserverId           int
	planetDropEffectObserverId      int
	planetFireDeathEffectObserverId int
	planetCollisionEffectObserverId int
	gameOverObserverId              int
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	a.backgroundMusic.Play(-1)
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if twodee.MusicIsPaused() {
		twodee.ResumeMusic()
	}
}

func (a *AudioSystem) PlayPlanetDropEffect(e twodee.GETyper) {
	if a.planetDropEffect.IsPlaying(2) == 0 {
		a.planetDropEffect.PlayChannel(2, 1)
	}
}

func (a *AudioSystem) PlayPlanetFireDeathEffect(e twodee.GETyper) {
	if a.planetFireDeathEffect.IsPlaying(3) == 0 {
		a.planetFireDeathEffect.PlayChannel(3, 1)
	}
}

func (a *AudioSystem) PlayPlanetCollisionEffect(e twodee.GETyper) {
	if a.planetCollisionEffect.IsPlaying(4) == 0 {
		a.planetCollisionEffect.PlayChannel(4, 1)
	}
}

func (a *AudioSystem) OnGameOver(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	// TODO(kalev): Play a game over effect.
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.backgroundMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ReleasePlanet, a.planetDropEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlanetFireDeath, a.planetFireDeathEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlanetCollision, a.planetCollisionEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(GameOver, a.gameOverObserverId)
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
	planetDropEffect.SetVolume(100)
	planetFireDeathEffect.SetVolume(60)
	planetCollisionEffect.SetVolume(60)
	audioSystem.backgroundMusicObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.planetDropEffectObserverId = app.GameEventHandler.AddObserver(ReleasePlanet, audioSystem.PlayPlanetDropEffect)
	audioSystem.planetFireDeathEffectObserverId = app.GameEventHandler.AddObserver(PlanetFireDeath, audioSystem.PlayPlanetFireDeathEffect)
	audioSystem.planetCollisionEffectObserverId = app.GameEventHandler.AddObserver(PlanetCollision, audioSystem.PlayPlanetCollisionEffect)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	audioSystem.gameOverObserverId = app.GameEventHandler.AddObserver(GameOver, audioSystem.OnGameOver)
	return
}
