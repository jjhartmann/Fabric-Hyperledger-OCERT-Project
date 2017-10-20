/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hyperledger/fabric-ca/api"
	. "github.com/hyperledger/fabric-ca/lib"
	"github.com/hyperledger/fabric-ca/lib/dbutil"
	"github.com/hyperledger/fabric-ca/lib/spi"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "/tmp/dbtesting"

	sqliteTruncateTables = `
DELETE FROM Users;
DELETE FROM affiliations;
`
)

type TestAccessor struct {
	Accessor *Accessor
	DB       *sqlx.DB
}

func (ta *TestAccessor) Truncate() {
	Truncate(ta.DB)
}

func TestSQLite(t *testing.T) {
	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dbPath, 0755)
		}
	} else {
		os.RemoveAll(dbPath)
		os.MkdirAll(dbPath, 0755)
	}
	dataSource := dbPath + "/fabric-ca.db"
	db, _, err := dbutil.NewUserRegistrySQLLite3(dataSource)
	if err != nil {
		t.Error("Failed to open connection to DB")
	}
	accessor := NewDBAccessor()
	accessor.SetDB(db)

	ta := TestAccessor{
		Accessor: accessor,
		DB:       db,
	}
	testEverything(ta, t)
	removeDatabase()
}

// Truncate truncates the DB
func Truncate(db *sqlx.DB) {
	var sql []string
	sql = []string{sqliteTruncateTables}

	for _, expr := range sql {
		if len(strings.TrimSpace(expr)) == 0 {
			continue
		}
		if _, err := db.Exec(expr); err != nil {
			panic(err)
		}
	}
}

func removeDatabase() {
	os.RemoveAll(dbPath)
}

func testEverything(ta TestAccessor, t *testing.T) {
	testInsertAndGetUser(ta, t)
	testDeleteUser(ta, t)
	testUpdateUser(ta, t)
	testInsertAndGetAffiliation(ta, t)
	testDeleteAffiliation(ta, t)
}

func testInsertAndGetUser(ta TestAccessor, t *testing.T) {
	t.Log("TestInsertAndGetUser")
	ta.Truncate()

	insert := spi.UserInfo{
		Name:       "testId",
		Pass:       "123456",
		Type:       "client",
		Attributes: []api.Attribute{},
	}

	err := ta.Accessor.InsertUser(insert)
	if err != nil {
		t.Errorf("Error occured during insert query of ID: %s, error: %s", insert.Name, err)
	}

	user, err := ta.Accessor.GetUser(insert.Name, nil)
	if err != nil {
		t.Errorf("Error occured during querying of id: %s, error: %s", insert.Name, err)
	}

	if user.GetName() != insert.Name {
		t.Error("Incorrect ID retrieved")
	}
}

func testDeleteUser(ta TestAccessor, t *testing.T) {
	t.Log("TestDeleteUser")
	ta.Truncate()

	insert := spi.UserInfo{
		Name:       "testId",
		Pass:       "123456",
		Type:       "client",
		Attributes: []api.Attribute{},
	}

	err := ta.Accessor.InsertUser(insert)
	if err != nil {
		t.Errorf("Error occured during insert query of id: %s, error: %s", insert.Name, err)
	}

	err = ta.Accessor.DeleteUser(insert.Name)
	if err != nil {
		t.Errorf("Error occured during deletion of ID: %s, error: %s", insert.Name, err)
	}

	_, err = ta.Accessor.GetUser(insert.Name, nil)
	if err == nil {
		t.Error("Should have errored, and not returned any results")
	}
}

func testUpdateUser(ta TestAccessor, t *testing.T) {
	t.Log("TestUpdateUser")
	ta.Truncate()

	insert := spi.UserInfo{
		Name:           "testId",
		Pass:           "123456",
		Type:           "client",
		Attributes:     []api.Attribute{},
		MaxEnrollments: 1,
	}

	err := ta.Accessor.InsertUser(insert)
	if err != nil {
		t.Errorf("Error occured during insert query of ID: %s, error: %s", insert.Name, err)
	}

	insert.Pass = "654321"

	ta.Accessor.UpdateUser(insert)
	if err != nil {
		t.Errorf("Error occured during update query of ID: %s, error: %s", insert.Name, err)
	}

	user, err := ta.Accessor.GetUser(insert.Name, nil)
	if err != nil {
		t.Errorf("Error occured during querying of ID: %s, error: %s", insert.Name, err)
	}

	err = user.Login(insert.Pass, -1)
	if err != nil {
		t.Error("Failed to login in user: ", err)
	}

}

func testInsertAndGetAffiliation(ta TestAccessor, t *testing.T) {
	ta.Truncate()

	err := ta.Accessor.InsertAffiliation("Bank1", "Banks")
	if err != nil {
		t.Errorf("Error occured during insert query of group: %s, error: %s", "Bank1", err)
	}

	group, err := ta.Accessor.GetAffiliation("Bank1")
	if err != nil {
		t.Errorf("Error occured during querying of name: %s, error: %s", "Bank1", err)
	}

	if group.GetName() != "Bank1" {
		t.Error("Failed to query")
	}

}

func testDeleteAffiliation(ta TestAccessor, t *testing.T) {
	ta.Truncate()

	err := ta.Accessor.InsertAffiliation("Banks.Bank2", "Banks")
	if err != nil {
		t.Errorf("Error occured during insert query of group: %s, error: %s", "Bank2", err)
	}

	err = ta.Accessor.DeleteAffiliation("Banks.Bank2")
	if err != nil {
		t.Errorf("Error occured during deletion of group: %s, error: %s", "Bank2", err)
	}

	_, err = ta.Accessor.GetAffiliation("Banks.Bank2")
	if err == nil {
		t.Error("Should have errored, and not returned any results")
	}
}
