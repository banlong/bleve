package index

import "github.com/blevesearch/bleve"

type IndexMap struct {
	Path 	string
	Index 	bleve.Index
}

func NewIndexMap(path string) *IndexMap  {
	return &IndexMap{
		Path: path,
		Index: NewIndex(path),
	}
}

func (im *IndexMap)AddIndex(indexValue string, value interface{})  {
	im.Index.Index(indexValue, value)
}

func (im *IndexMap)ExecQuery(queryString string)  *bleve.SearchResult{
	//Declare a search request
	stringQuery := bleve.NewQueryStringQuery(queryString)
	searchRequest := bleve.NewSearchRequest(stringQuery)

	//Execute search
	searchResult, _ := im.Index.Search(searchRequest)
	return searchResult
}


func NewIndex(path string)  bleve.Index{
	var mapping *bleve.IndexMapping
	opindex, err := bleve.Open(path)
	if err != nil{
		mapping = bleve.NewIndexMapping()
		opindex, err = bleve.New(path, mapping)
	}
	return opindex
}