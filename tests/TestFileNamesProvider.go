package tests

import (
	"fmt"
	"narezator/filelist"
	"testing"
)

func TestFileNamesProvider(t *testing.T) {
	t.Run("Should provide file list", func(t *testing.T) {
		fl := filelist.NewDarwinFileList()
		fmt.Println(fl)
	})
}
