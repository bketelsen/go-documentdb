package docdb

import (
	"os"
	"testing"
)

var (
	databaseUrl string = os.Getenv("DATABASE_URL")
	databaseKey string = os.Getenv("DATABASE_KEY")
)

type TestItem struct {
	Name string `json:"name,omitempty"`
	Resource
}

func Test_DatabaseOperations(t *testing.T) {
	client := NewClient(databaseUrl, Config{databaseKey})

	t.Log("Create database")
	db, err := client.CreateDatabase("foo")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("Create collection for db id %v\n", db.Self)
	coll, err := client.CreateCollection(db.Self, "bar", nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log("Create document")
	createItem := TestItem{"test document", Resource{Id: "id123"}}
	err = client.CreateDocument(coll.Self, &createItem)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("Read document id %v\n", createItem.Self)
	readItem := TestItem{}
	err = client.ReadDocument(createItem.Self, &readItem)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("%#v", readItem)

	var queryItems []TestItem
	query := "SELECT * FROM bar"
	t.Logf("Performing query %v", query)
	err = client.QueryDocuments(coll.Self, query, &queryItems)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v item(s) found.", len(queryItems))

	t.Logf("Delete document %v", readItem.Self)
	err = client.DeleteDocument(readItem.Self)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("Delete collection %v", coll.Self)
	err = client.DeleteCollection(coll.Self)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("Delete database %v", db.Self)
	err = client.DeleteDatabase(db.Self)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
