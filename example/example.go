//go:generate ../bin/go-generated-views -b example

package example

// Book describes one book
type Book struct {
	ID          uint64 `binding(get):"required,gt=0" json:"id" json(add):"-,omitempty"`
	Name        string `json:"name" binding(add):"required"`
	Description string `json:"description"`
}
