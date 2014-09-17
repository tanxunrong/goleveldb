package table

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tanxunrong/goleveldb/leveldb/testutil"
)

func TestTable(t *testing.T) {
	testutil.RunDefer()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Table Suite")
}
