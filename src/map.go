package main

import (
	twodee "../libs/twodee"
	"fmt"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"path/filepath"
)

func getTexturePath(m *tmxgo.Map, path string) (out string, err error) {
	var prefix = filepath.Dir(path)
	for i := 0; i < len(m.Tilesets); i++ {
		if m.Tilesets[i].Image == nil {
			continue
		}
		out = filepath.Join(prefix, m.Tilesets[i].Image.Source)
		return
	}
	err = fmt.Errorf("Could not find suitable tileset")
	return
}

func LoadMap(path string) (batch *twodee.Batch, err error) {
	var (
		tilemeta twodee.TileMetadata
		maptiles []*tmxgo.Tile
		textiles []twodee.TexturedTile
		maptile  *tmxgo.Tile
		m        *tmxgo.Map
		i        int
		data     []byte
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	if path, err = getTexturePath(m, path); err != nil {
		return
	}
	tilemeta = twodee.TileMetadata{
		Path:      path,
		PxPerUnit: int(PxPerUnit),
		Interpolation: twodee.LinearInterpolation,
	}
	if maptiles, err = m.TilesFromLayerName("tiles"); err != nil {
		return
	}
	textiles = make([]twodee.TexturedTile, len(maptiles))
	for i, maptile = range maptiles {
		if maptile != nil {
			textiles[i] = maptile
		}
	}
	if batch, err = twodee.LoadBatch(textiles, tilemeta); err != nil {
		return
	}
	return
}
