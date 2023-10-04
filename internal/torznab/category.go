package torznab

import (
	_ "embed"
)

type Category struct {
	ID     int           `xml:"id,attr"`
	Name   string        `xml:"name,attr"`
	Subcat []Subcategory `xml:"subcat"`
}

type Subcategory struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Has returns true if the category has the given id or if any of its subcategories have the given id.
func (c Category) Has(id int) bool {
	if c.ID == id {
		return true
	}
	for _, sub := range c.Subcat {
		if sub.ID == id {
			return true
		}
	}
	return false
}
