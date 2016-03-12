package main

import (
	"bytes"
	"encoding/base64"
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
	var err error
	CurrentUser, err = user.Current()
	if err != nil {
		log.Fatal("Unexpected error!")
		return
	}
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {

	var (
		setServer      bool
		setUser        bool
		addUser        bool
		changePassword bool
		deleteUser     bool
		upload         string
		withName       string
	)

	flag.BoolVar(&setServer, "setserver", false, "Server URL")
	flag.BoolVar(&setUser, "setuser", false, "User and Password")
	flag.BoolVar(&addUser, "adduser", false, "Add new user")
	flag.BoolVar(&changePassword, "changepassword", false, "Change user password")
	flag.BoolVar(&deleteUser, "deleteuser", false, "Delete user")

	flag.StringVar(&upload, "upload", "", "Folder or File to be uploaded")
	flag.StringVar(&withName, "withname", "", "Custom name to store")
	flag.Parse()

	if err := LoadConfiguration(); err != nil {
		log.Fatal("Error to load or write config file!")
		return
	}

	if setServer {
		newServer := flag.Args()[0]
		if newServer == "" {
			log.Fatal("Please, provide a valid server url")
			return
		}
		Config.ServerURL = newServer
		WriteConfiguration()
		log.Println("Server OK")
		return
	}
	// if !*configMode || *server == "" || *userName == "" || *userPassword == "" {
	// 		log.Fatal("Error loading the configuration file.")
	// 		return
	// 	}
	// 	Config = &Configuration{}
	// 	Config.ServerURL = *server
	// 	if *userName != "" && *userPassword != "" {
	// 		Config.AccessKey = basicAuth(*userName, *userPassword)
	// 	}
	// 	WriteConfiguration(CurrentUser.HomeDir)
	// 	return
	// }

	err := Zipit(upload, withName)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := sendPost(ZipFileName)
	if err != nil {
		log.Println(err)
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
	r.Header.Set("Authorization", "Basic "+Config.AccessKey)
	client := &http.Client{}
	response, err := client.Do(r)
	return response, err
}
