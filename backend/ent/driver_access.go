package ent

import "entgo.io/ent/dialect"

// Driver exposes the underlying dialect.Driver for infrastructure wiring (e.g. *sql.DB extraction).
func (c *Client) Driver() dialect.Driver {
	return c.driver
}
