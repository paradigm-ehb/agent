// Package that handles the parsing of dbus types
package dbushandler

// TODO: fix the return type to something explicit
// ListUnits(out a(ssssssouso) units);

// Method
// parses
// @return [][]string
func Parse(ch chan [][]string) {

	in := <-ch

	for i := range in {
		for j := range i {
			in[i][j] += "\n"
		}
	}
	ch <- in
}
