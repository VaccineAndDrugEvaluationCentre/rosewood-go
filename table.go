package rosewood

//Table holds all the info needed to render a table
type table struct {
	identifier string
	contents   *tableContents
	caption    *section
	header     *section
	footnotes  *section
	cmdList    []*Command
}

func newTable() *table {
	return &table{}
}

func (t *table) Merge(ra Range) error {
	//	fmt.Println("range in Merge:", ra)
	return t.contents.merge(ra)
}
