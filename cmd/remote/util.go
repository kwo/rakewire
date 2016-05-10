package remote

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	errMissingHost     = errors.New("Missing host/port, try setting the --host option")
	errMissingUsername = errors.New("Missing username, try setting the --username option")
	errMissingPassword = errors.New("Missing password, try setting the --password option")
)

func makeRequest(c *cli.Context, path string, req interface{}, rsp interface{}) error {

	addr, username, password, token, errCredentials := getHostUsernamePasswordToken(c)
	if errCredentials != nil {
		return errCredentials
	}

	var auth string
	if len(token) == 0 {
		content := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		auth = "Basic " + content
	} else {
		auth = "Bearer " + token
	}

	reqData, errMarshal := json.Marshal(req)
	if errMarshal != nil {
		return errMarshal
	}

	request, errRequest := http.NewRequest(http.MethodPost, "https://"+addr+"/api/"+path, bytes.NewBuffer(reqData))
	if errRequest != nil {
		return errRequest
	}
	request.Header.Set("User-Agent", getAppNameAndVersion(c))
	request.Header.Add("Authorization", auth)
	client := &http.Client{}
	response, errResponse := client.Do(request)
	if errResponse != nil {
		return errResponse
	}

	defer response.Body.Close()
	rspData, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		return errRead
	}
	if errUnmarshal := json.Unmarshal(rspData, rsp); errUnmarshal != nil {
		return errUnmarshal
	}

	return nil

}

func getAppNameAndVersion(c *cli.Context) string {
	appName := strings.Fields(c.App.Name)[0]
	version := c.App.Version
	return strings.TrimSpace(appName + " " + version)
}

func getHostUsernamePasswordToken(c *cli.Context) (host, username, password, token string, err error) {

	host = c.Parent().String("host")
	username = c.Parent().String("username")
	password = c.Parent().String("password")
	token = c.Parent().String("token")

	if len(host) == 0 {
		err = errMissingHost
		return
	}

	if len(token) == 0 {

		if len(username) == 0 {
			err = errMissingUsername
			return
		}

		if len(password) == 0 {
			err = errMissingPassword
			return
		}

	}

	return

}
