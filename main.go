package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/user"
)

func init() {
	// load config
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Unexpected error!")
		return
	}
	if err := LoadConfigurationFrom(currentUser.HomeDir); err != nil {
		log.Fatal("Error loading the configuration file")
		return
	}
}

func main() {

	var file string
	var fileName string
	flag.StringVar(&file, "f", "", "Folder or File to be compressed")
	flag.StringVar(&fileName, "o", "", "Output static folder or file")
	flag.Parse()

	err := Zipit(file, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := sendPost(ZipFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	text, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(text))

}

func sendPost(path string) (*http.Response, error) {

	file, err := os.Open("./" + path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", Config.ServerURL+"/upload/", body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Header.Set("Authorization", "Basic YWRtaW46MTIz")
	client := &http.Client{}
	response, err := client.Do(r)
	return response, err
}
