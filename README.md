LD30
====

Ludum Dare 30 - theme 'Connected Worlds'.

You are trying to populate a solar system with life.  You can place Earth-like planets in orbit around a star.  Every planet pulls on each other so you need to balance the gravity wells of every object in the system.  When a planet falls into a tolerable "life zone", then its population grows according to a standard population growth curve.  If a planet goes too close to the sun it burns up, if it falls too far away, it freezes.  If a planet collides with another planet or the sun it is destroyed.  All of these have disastrous consequences for any life on the planet.  The goal is to have the highest possible aggregate population in the solar system before the sun goes supernova.

## TODO

  * [x] Come up with idea.
  * [ ] Integrate twodee library
  * [ ] Open window, blank screen
  * [ ] Main game loop
  * [ ] Render sun
  * [ ] Drop planet when click
  * [ ] Have planet orbit sun
  * [ ] Multiple planets
  * [ ] Placeholder art
  * [ ] Music loading
  * [ ] Heads up display of temp, population for planets
    
## Setup

Complete the setup steps for the twodee lib.

Also:

	go get -u github.com/kurrik/tmxgo

Run:

	git submodule init
	git submodule update
	
## Brainstorming
Ideas

### Jumping back and forth between worlds
 
  * Close to what we had with Heavy Drizzle
  * Probably what most people will pick
 
### "Get your ass to Mars"
 
  * Need to terraform / colonize a planet
  * Shuttle makes constant repeating trips back and forth
  * Need to launch payloads to be picked up
  * Limited payload size
  * Need to balance replenishing resources / building new terraforming equipment, and ferrying people back and forth.
  * Goal is to move entire population
  * What is difficult?  What will make people play differently?
   
### Linked characters

  * Two characters in split screen, wandering.
  * Controls move both at once
  * Need to run one into walls, etc, to make them line up correctly
  
### Orbiting planets

  * Place planets aside a gravity well
  * They start with a velocity and begin to orbit
  * Try to place more planets, all gravity calculated to adjust orbits
  * Planets can fly off into space, collide with each other, fly into sun
  * What's the point of the game?  Try to get X planets in a stable orbit?