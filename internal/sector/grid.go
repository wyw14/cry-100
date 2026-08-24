package sector

import (
	"fmt"
	"github.com/wyw14/cry-100/internal/model"
	"sort"
	"sync"
)

type Sector struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Center    model.Point `json:"center"`
	Neighbors []string    `json:"neighbors"`
}
type Grid struct {
	mu      sync.RWMutex
	sectors map[string]Sector
}

func NewGrid() *Grid { return &Grid{sectors: make(map[string]Sector)} }
func (g *Grid) Add(sector Sector) error {
	if sector.ID == "" || sector.Name == "" {
		return fmt.Errorf("sector id and name are required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.sectors[sector.ID]; ok {
		return fmt.Errorf("sector %s already exists", sector.ID)
	}
	sector.Neighbors = uniqueSorted(sector.Neighbors)
	g.sectors[sector.ID] = sector
	return nil
}
func (g *Grid) List() []Sector {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]Sector, 0, len(g.sectors))
	for _, sector := range g.sectors {
		result = append(result, sector)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value != "" {
			seen[value] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
