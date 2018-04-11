package cachego

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type MySQLTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
	db     *sql.DB
}

var (
	cacheTable string = "cache"
)

func (s *MySQLTestSuite) SetupTest() {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cachedb")

	if err != nil {
		s.T().Skip()
	}

	s.cache, err = NewMySQL(db, cacheTable)

	if err != nil {
		s.T().Skip()
	}

	s.assert = assert.New(s.T())
	s.db = db
}

func (s *MySQLTestSuite) TearDownTest() {
}

func (s *MySQLTestSuite) TestCreateInstanceThrowAnError() {
	s.db.Close()

	_, err := NewMySQL(s.db, cacheTable)

	s.assert.Error(err)
}

func (s *MySQLTestSuite) TestSaveThrowAnError() {
	s.db.Close()

	s.assert.Error(s.cache.Save("foo", "bar", 0))
}

func (s *MySQLTestSuite) TestSaveThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(s.cache.Save("foo", "bar", 0))
}

func (s *MySQLTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *MySQLTestSuite) TestFetchThrowAnError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 1)

	result, err := s.cache.Fetch(key)

	s.assert.Error(err)
	s.assert.Empty(result)
}

func (s *MySQLTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *MySQLTestSuite) TestFetchWithLongLifetime() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 10*time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *MySQLTestSuite) TestContainsThrowAnError() {
	s.assert.False(s.cache.Contains("bar"))
}

func (s *MySQLTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *MySQLTestSuite) TestDeleteThrowAnError() {
	s.db.Close()

	s.assert.Error(
		s.cache.Delete("cccc"),
	)
}

func (s *MySQLTestSuite) TestDeleteThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(
		s.cache.Delete("cccc"),
	)
}

func (s *MySQLTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.Nil(s.cache.Delete("foo"))
}

func (s *MySQLTestSuite) TestFlushThrowAnError() {
	s.db.Close()

	s.assert.Error(s.cache.Flush())
}

func (s *MySQLTestSuite) TestFlushThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(s.cache.Flush())
}

func (s *MySQLTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *MySQLTestSuite) TestFetchMultiReturnNoItemsWhenThrowAnError() {
	s.db.Close()

	result := s.cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *MySQLTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *MySQLTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestMySQLRunSuite(t *testing.T) {
	suite.Run(t, new(MySQLTestSuite))
}
