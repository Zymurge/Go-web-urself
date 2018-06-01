package types

import (
//	"webstuff/types"
)

// Grid is a collection of Locs and helper functions to work with them
type Grid struct {
	locs map[string]Loc
	xmin, xmax int
	ymin, ymax int
	zmin, zmax int
}

// Build creates a grid of Loc objects spanning a size of x by y. the z axis will be calculated as a funciton 
// of x and y. The center of the grid will be 0,0,0.
func (g *Grid) Build(xSize int, ySize int) {
	g.locs = map[string]Loc{}

	g.xmax = xSize/2
	g.xmin = g.xmax * -1
	g.ymax = ySize/2
	g.ymin = g.ymax * -1

	for x := g.xmin; x<=g.xmax; x++ {
		for y := g.ymin; y<=g.ymax; y++ {
			z := (x + y) * -1
			if z < g.zmin {g.zmin = z}
			if z > g.zmax {g.zmax = z}
			if loc, err := LocFromCoords(x,y,z); err == nil {
				g.locs[loc.ID] = loc
			}
		}
	}
}

// GetLoc returns the loc with the specified ID
func (g *Grid) GetLoc(id string) Loc {
	return g.locs[id]
}

// XMin getter
func (g *Grid) XMin() int { return g.xmin }

// XMax getter
func (g *Grid) XMax() int { return g.xmax }

// YMin getter
func (g *Grid) YMin() int { return g.ymin }

// YMax getter
func (g *Grid) YMax() int { return g.ymax }

// ZMin getter
func (g *Grid) ZMin() int { return g.zmin }

// ZMax getter
func (g *Grid) ZMax() int { return g.zmax }
