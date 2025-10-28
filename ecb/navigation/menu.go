package navigation

import (
	"context"
	"sort"
	"sync"
)

// Record is a flat representation from persistence.
type Record struct {
	ID       int64
	ParentID *int64
	Title    string
	Icon     string
	URL      string
	Order    int
}

// Item is the hierarchical menu node used by the GUI.
type Item struct {
	ID       int64
	Title    string
	Icon     string
	URL      string
	Order    int
	Children []*Item
}

// BreadcrumbEntry captures the path context for the active screen.
type BreadcrumbEntry struct {
	Title string
	Icon  string
	URL   string
}

// Repository hides storage details (DB, file, mock).
type Repository interface {
	FetchTree(ctx context.Context) ([]Record, error)
}

// Service builds menus and breadcrumbs for the Fyne navigation.
type Service struct {
	repo Repository

	mu         sync.RWMutex
	root       []*Item
	indexByURL map[string]*Item
	indexByID  map[int64]*Item
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Refresh reloads the navigation records and rebuilds the tree.
func (s *Service) Refresh(ctx context.Context) error {
	records, err := s.repo.FetchTree(ctx)
	if err != nil {
		return err
	}
	tree, byURL, byID := buildTree(records)
	s.mu.Lock()
	s.root = tree
	s.indexByURL = byURL
	s.indexByID = byID
	s.mu.Unlock()
	return nil
}

// MainMenu returns the cached root menu items.
func (s *Service) MainMenu() []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneTree(s.root)
}

// Breadcrumb builds the breadcrumb chain for the given URL.
func (s *Service) Breadcrumb(url string) []BreadcrumbEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	node := s.indexByURL[url]
	if node == nil {
		return nil
	}
	var chain []BreadcrumbEntry
	for current := node; current != nil; current = parentOf(current, s.indexByID) {
		chain = append([]BreadcrumbEntry{{Title: current.Title, Icon: current.Icon, URL: current.URL}}, chain...)
	}
	return chain
}

// CurrentMenu finds a menu entry by URL.
func (s *Service) CurrentMenu(url string) *Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.indexByURL == nil {
		return nil
	}
	if node, ok := s.indexByURL[url]; ok {
		return cloneNode(node)
	}
	return nil
}

// CurrentMenuByID finds a menu entry by ID.
func (s *Service) CurrentMenuByID(id int64) *Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.indexByID == nil {
		return nil
	}
	if node, ok := s.indexByID[id]; ok {
		return cloneNode(node)
	}
	return nil
}

func buildTree(records []Record) ([]*Item, map[string]*Item, map[int64]*Item) {
	nodes := make(map[int64]*Item, len(records))
	root := make([]*Item, 0, len(records))
	for _, rec := range records {
		nodes[rec.ID] = &Item{
			ID:    rec.ID,
			Title: rec.Title,
			Icon:  rec.Icon,
			URL:   rec.URL,
			Order: rec.Order,
		}
	}
	for _, rec := range records {
		node := nodes[rec.ID]
		if rec.ParentID == nil {
			root = append(root, node)
			continue
		}
		parent := nodes[*rec.ParentID]
		if parent == nil {
			root = append(root, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	sortChildren(root)
	indexByURL := make(map[string]*Item, len(nodes))
	indexByID := make(map[int64]*Item, len(nodes))
	for _, node := range nodes {
		if node.URL != "" {
			indexByURL[node.URL] = node
		}
		indexByID[node.ID] = node
	}
	return root, indexByURL, indexByID
}

func sortChildren(nodes []*Item) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Order == nodes[j].Order {
			return nodes[i].Title < nodes[j].Title
		}
		return nodes[i].Order < nodes[j].Order
	})
	for _, node := range nodes {
		if len(node.Children) == 0 {
			continue
		}
		sortChildren(node.Children)
	}
}

func cloneTree(nodes []*Item) []*Item {
	if len(nodes) == 0 {
		return nil
	}
	out := make([]*Item, len(nodes))
	for i, node := range nodes {
		out[i] = cloneNode(node)
	}
	return out
}

func cloneNode(node *Item) *Item {
	if node == nil {
		return nil
	}
	copy := *node
	copy.Children = cloneTree(node.Children)
	return &copy
}

func parentOf(node *Item, index map[int64]*Item) *Item {
	if node == nil {
		return nil
	}
	for _, candidate := range index {
		for _, child := range candidate.Children {
			if child.ID == node.ID {
				return candidate
			}
		}
	}
	return nil
}
