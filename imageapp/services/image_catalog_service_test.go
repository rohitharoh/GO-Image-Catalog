package services

import (
	"fmt"
	"github.com/tb/image-catalog/common-packages/conf"
	"testing"
)




func Test_ListFilesInS3(t *testing.T) {

referenceId := "c3fd46be-fc49-4e8e-af91-acc8e6204cf4"
folderPath := "folder-1/folder-2"
	res, err := ListFilesInS3(referenceId,folderPath,0)
	if err != nil {
		fmt.Println("error--->",err)
	}
	fmt.Println(string(res))

}




func init() {
	conf.LoadConfigFile()
}
