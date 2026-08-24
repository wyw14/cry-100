package sector

func DefaultGrid() *Grid {
	grid := NewGrid()
	for _, sector := range []Sector{{ID: "north-ridge", Name: "North Ridge", Neighbors: []string{"pine-valley"}}, {ID: "pine-valley", Name: "Pine Valley", Neighbors: []string{"north-ridge", "river-pass"}}, {ID: "river-pass", Name: "River Pass", Neighbors: []string{"pine-valley"}}} {
		if err := grid.Add(sector); err != nil {
			panic(err)
		}
	}
	return grid
}
