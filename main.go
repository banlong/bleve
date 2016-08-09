package main

import (
	"log"
	"github.com/blevesearch/bleve"
	"bleve/index"
)


type Message  struct{
	Id 	string
	From 	string
	Body 	string
	Value   int
}

var bleveIdx bleve.Index

func main() {
	//log.Print("Hello")
	////bolt.Open("bolt.db", config.READWRITE, &bolt.Options{Timeout: 1 * time.Second})
	message :=  Message {
			Id:   "example",
			From: "marty.schoch@gmail.com",
			Body: "bleve indexing is easy",
			Value: 10,
		}

	im := index.NewIndexMap("examples")

	im.AddIndex(message.Id, message)

	result := im.ExecQuery("Value:>5")

	log.Println(result)


	//var opindex bleve.Index
	//var mapping *bleve.IndexMapping
	//indexPath := "example"
	//opindex, err := bleve.Open(indexPath)
	//if err != nil{
	//	mapping = bleve.NewIndexMapping()
	//	opindex, err = bleve.New("example", mapping)
	//}
	//
	////Add new index, field to be used as Index, we can define multiple index at the same time for one object
	//opindex.Index(message.Id, message)
	////This message.Id is the return value from the search result.
	////Assume I save message object into key-value bolt (key= message.Id, value = message).
	////If I want to search the messages in bolt that satisfy my conditions, if match found,
	////the search result will contain the key that I will then use it to retrive the object
	////from BoltDB.
	//
	//
	//// Case 1: search for the "easy". Plain terms without any other syntax are
	//// interpreted as a match query for the term in the default field.
	//// The default field is "_all" unless overridden in the index mapping.
	//searchPlanValue := bleve.NewQueryStringQuery("easy")
	//
	////Declare a search request
	//searchRequest := bleve.NewSearchRequest(searchPlanValue)
	//
	////Execute search
	//searchResult, _ := opindex.Search(searchRequest)
	//
	////Display result
	//if searchResult.Hits.Len() > 0{
	//	log.Println(searchResult.Hits[0].ID);
	//	log.Println(searchResult.Hits[0].Index);
	//	log.Println(searchResult.Hits[0].Fields);
	//	log.Println(searchResult.Hits[0].Fragments);
	//	log.Println(searchResult.Hits[0].Locations);
	//	log.Println(searchResult.Hits[0].Score);
	//}
	//log.Println("----------------------------------");
	//log.Println(searchResult);
	//
	//
	//// Case 2: Field Scoping, search for the "marty.schoch@gmail.com" in field "From".
	//searchValueInField := bleve.NewQueryStringQuery("From:marty.schoch@gmail.com")
	////Declare a search request
	//searchRequest = bleve.NewSearchRequest(searchValueInField)
	//
	////Execute search
	//searchResult, _ = opindex.Search(searchRequest)
	//
	////Display result
	//log.Println(searchResult);
	//
	//// Case 3: Required, Optional, and Exclusion.
	//// Example: +description:water -light beer will perform a Boolean Query that MUST satisfy
	//// the Match Query for the term water in the description field, MUST NOT satisfy the Match
	//// Query for the term light in the default field, and SHOULD satisfy the Match Query for
	//// the term beer in the default field.
	//
	//searchBooleanValue := bleve.NewQueryStringQuery("+From:marty.schoch@gmail.com -easy bleve")
	////Declare a search request
	//searchRequest = bleve.NewSearchRequest(searchBooleanValue)
	//
	////Execute search
	//searchResult, _ = opindex.Search(searchRequest)
	//
	////Display result
	//log.Println(searchResult);
	//
	////Case 4: Numeric Ranges
	////You can perform numeric ranges by using the >, >=, <, and <= operators, followed by a numeric value.
	////Example: abv:>10 will perform an Numeric Range Query on the abv field for values greater than ten.
	//
	//searchRannge := bleve.NewQueryStringQuery("Value:>5 Value:<11 ")
	////Declare a search request
	//searchRequest = bleve.NewSearchRequest(searchRannge)
	//
	////Execute search
	//searchResult, _ = opindex.Search(searchRequest)
	//
	////Display result
	//log.Println(searchResult);
}