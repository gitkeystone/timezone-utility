package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed cities.json
var citiesJSON []byte

//go:embed postcodes.json
var postcodesJSON []byte

// Store holds the resolvable location datasets.
type Store struct {
	byAlias   map[string][]*Location
	byPostal  map[string]*Location
	Locations []*Location
}

type postalEntry struct {
	Postal   string `json:"postal"`
	Location string `json:"location"`
}

// Load builds a Store from the embedded datasets.
func Load() (*Store, error) {
	var cities []*Location
	if err := json.Unmarshal(citiesJSON, &cities); err != nil {
		return nil, fmt.Errorf("解析地点数据失败: %w", err)
	}
	var postals []postalEntry
	if err := json.Unmarshal(postcodesJSON, &postals); err != nil {
		return nil, fmt.Errorf("解析邮编数据失败: %w", err)
	}

	byName := make(map[string]*Location, len(cities))
	store := &Store{
		byAlias:   make(map[string][]*Location),
		byPostal:  make(map[string]*Location),
		Locations: cities,
	}
	for _, loc := range cities {
		byName[strings.ToLower(loc.Name)] = loc
	}
	for _, p := range postals {
		loc, ok := byName[strings.ToLower(p.Location)]
		if !ok {
			return nil, fmt.Errorf("邮编 %q 引用了未知地点 %q", p.Postal, p.Location)
		}
		if _, exists := store.byPostal[p.Postal]; exists {
			return nil, fmt.Errorf("邮编 %q 重复", p.Postal)
		}
		store.byPostal[p.Postal] = loc
		loc.PostalCodes = append(loc.PostalCodes, p.Postal)
	}
	for _, loc := range cities {
		keys := append([]string{loc.Name}, loc.Aliases...)
		for _, k := range keys {
			key := strings.ToLower(k)
			store.byAlias[key] = append(store.byAlias[key], loc)
		}
	}
	return store, nil
}

// ResolvePlace resolves a place name (case-insensitive) to a location, detecting ambiguity.
func (s *Store) ResolvePlace(name string) (*Location, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	locs := s.byAlias[key]
	switch len(locs) {
	case 0:
		return nil, fmt.Errorf("无法识别地点名称 %q", name)
	case 1:
		return locs[0], nil
	default:
		names := make([]string, len(locs))
		for i, l := range locs {
			names[i] = fmt.Sprintf("%s (%s)", l.Name, l.Timezone)
		}
		return nil, fmt.Errorf("地点名称 %q 有歧义，匹配到: %s", name, strings.Join(names, "、"))
	}
}

// ResolvePostal resolves a postal code to a location.
func (s *Store) ResolvePostal(code string) (*Location, error) {
	loc, ok := s.byPostal[strings.TrimSpace(code)]
	if !ok {
		return nil, fmt.Errorf("无法识别邮政编码 %q", code)
	}
	return loc, nil
}

// PlaceNames returns the sorted known aliases.
func (s *Store) PlaceNames() []string {
	names := make([]string, 0, len(s.byAlias))
	for k := range s.byAlias {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
