package generator

// Book describes one book, it doesn't have any tag
type Book struct {
	ID           uint64
	Name         string
	Description  string
	privateField string
}

// Author describes one author, it doesn't have any view in tags
type Author struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Description  string
	privateField string
}

// Page describes one page, it has one view "Special"
type Page struct {
	Number  uint8  `json:"number" json(special):"specialNumber" special:"number"`
	Content string `json(special):"specialContent" json:"content" db(special):"special_content"`
}

// Paragraph describes paragraph, it has two views - "short", "extended"
type Paragraph struct {
	Title   string `json:"title"`
	Excerpt string `db:"excerpt" json(short,extended):"excerpt"`
	Content string `json(short):"-,omitempty"  json(extended):"content"`
}

// Sentence describes sentence, and it has incorrect specified views (should be copied as is)
type Sentence struct {
	ID      uint64 `json:"id" json(short):"ID:ptr"`
	Content string `json:"content" json(short:"content"`
	Visible bool   `json(long):"visible(:)" json:"visible"`
}
