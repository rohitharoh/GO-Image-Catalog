package services

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
	"github.com/tb/image-catalog/common-packages/system"
	"time"

	"encoding/json"
	_ "errors"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/tb/image-catalog/common-packages/system"
	"io"
	"net/http"
	"os"
	"strings"

	"errors"
	"github.com/Machiel/slugify"
)

//UploadFileMedia will upload multiple files to S3, folderPath specifies the folders
func UploadFileMedia(f *os.File, filePath string, fileName string, folderPath string, referenceId string) (map[string]interface{}, error) {

	fmt.Println("filePath", filePath)
	fileContent, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {

		return nil, errors.New("fail read write error")
	}

	buff := make([]byte, 512)

	n, err := fileContent.Read(buff)
	if err != nil && err != io.EOF {

		return nil, errors.New("Reading file content error")
	}

	contentType := http.DetectContentType(buff[:n])
	err = fileContent.Close()
	if err != nil {

		return nil, err
	}
	key := "image-catalog" + "/" + referenceId + "/" + folderPath + "/" + fileName

	orgFileNameParts := strings.Split(fileName, ".")
	fileName = slugify.Slugify(strings.Join(orgFileNameParts[:len(orgFileNameParts)-1], ".")) + "." + orgFileNameParts[len(orgFileNameParts)-1]

	upFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	_, err = upFile.Read(fileBuffer)
	if err != nil {

		return nil, err
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(viper.GetString("aws.s3Region"))}))

	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:        bytes.NewReader(fileBuffer),
		Bucket:      aws.String(viper.GetString("aws.s3BucketName")),
		Key:         aws.String(key),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return nil, errors.New("upload fail error")
	}

	err = f.Close()
	if err != nil {

		return nil, err
	}
	//removed bucket name and changed the cdnURL
	s3UrlLocation := viper.GetString("aws.cdnUrl") + "image-catalog" + "/" + referenceId + "/" + folderPath + "/" + fileName

	response := make(map[string]interface{})
	response["location"] = s3UrlLocation
	response["fileName"] = fileName
	response["referenceId"] = referenceId
	return response, nil

}

//ListFilesInS3 lists all files or folders, whichever is present within the image-catalog bucket.
func ListFilesInS3(referenceID, folderPath string, skip int) ([]byte, error) {

	var response []map[string]interface{}

	svc := s3.New(session.New(&aws.Config{Region: aws.String(viper.GetString("aws.s3Region"))}))
	var subFolder string
	if folderPath == "" {
		subFolder = "image-catalog" + "/" + referenceID

	} else {

		subFolder = "image-catalog" + "/" + referenceID + "/" + folderPath

	}



	input := &s3.ListObjectsInput{
		Bucket:    aws.String(viper.GetString("aws.s3BucketName")),
		Prefix:    aws.String(subFolder + "/"),
		Delimiter: aws.String("/"),
	}

	result, err := svc.ListObjects(input)
	/*fmt.Println("common prefixes", result.CommonPrefixes)

	fmt.Println("file contents", len(result.Contents))
	fmt.Println("file contents", result.Contents)*/
	totalskip := skip
	skip = skip * system.LIMIT_RECORD_FOR_IMAGE_LIST

	var isMoreList bool
	if len(result.Contents) > skip+system.LIMIT_RECORD_FOR_IMAGE_LIST {
		isMoreList = true
	} else {
		isMoreList = false
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())

		} else {

			fmt.Println(err.Error())
		}
		return nil, err
	}

	t, tErr := time.Parse(system.DEFAULT_DATE_FORMAT, system.DATE_SPECIFIED)
	if tErr != nil {
		fmt.Printf("%v\nPlease enter the date as YYYY-MM-DD (e.g. 2016-09-25)", tErr)
		return nil, system.DateFormatMismatchErr
	}

	count := 0
	limit := system.LIMIT_RECORD_FOR_IMAGE_LIST
	s := totalskip * limit

	if len(result.Contents) > 0 {
		for k, item := range result.Contents {

			temp := make(map[string]interface{})
			fmt.Println("key", *item.Key)
			var mainFolder []string
			split := strings.Contains(*item.Key, "/")
			if split {
				mainFolder = strings.Split(*item.Key, "/")

			}

			if *item.Key == subFolder+"/" {

				fmt.Println("it is a directory")
			} else {
				if k >= s {
					if count < limit {

						specifiedTime := *item.LastModified
						if specifiedTime.After(t) {

							fileName := mainFolder[len(mainFolder)-1]
							displayFileName := strings.Split(fileName, ".")

							s3UrlLocation := viper.GetString("aws.cdnUrl") + *item.Key
							temp["isFolder"] = false
							temp["actualFileName"] = fileName
							temp["displayFileName"] = displayFileName[len(displayFileName)-2]
							temp["s3UrlLocation"] = s3UrlLocation
							temp["Size"] = *item.Size
							temp["lastModified"] = *item.LastModified
							temp["storageClass"] = *item.StorageClass
							response = append(response, temp)
						}
					}
					count++
				}

			}
		}
	}
if len(result.CommonPrefixes) > 0 {
	for j, prefix := range result.CommonPrefixes {

		if j >= s {
			if count < limit {

				temp := make(map[string]interface{})
				folderName := strings.Split(*prefix.Prefix, "/")

				temp["isFolder"] = true
				temp["folderName"] = folderName[len(folderName)-2]
				response = append(response, temp)
			}
			count++
		}
	}
}

	responseFinal := make(map[string]interface{})
	responseFinal["isMoreList"] = isMoreList
	responseFinal["fileList"] = response

	finalResponse, err := json.Marshal(responseFinal)
	if err != nil {

		return nil, err
	}
	return finalResponse, nil
}
