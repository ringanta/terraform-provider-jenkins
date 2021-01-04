package jenkins

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"text/template"

	jenkins "github.com/bndr/gojenkins"
)

type jenkinsClient interface {
	GetLocalUser(username string) (jenkinsLocalUser, error)
	CreateLocalUser(username string, password string, fullname string, email string, description string) error
	DeleteLocalUser(username string) error
	GetUserPermissions(username string) (jenkinsUserPermissions, error)
	CreateUserPermissions(username string, permissions []string) error
	PostScript(payload bytes.Buffer, respStruct interface{}) error
}

type jenkinsLocalUser struct {
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	PasswordHash string `json:"password_hash"`
	Username     string `json:"username"`
	Description  string `json:"description"`
}

type jenkinsLocalUserCreate struct {
	Password string
	jenkinsLocalUser
}

type jenkinsResponse struct {
	Error   bool             `json:"error"`
	Message string           `json:"msg"`
	Data    jenkinsLocalUser `json:"data"`
}

type jenkinsUserPermissions struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

type jenkinsResponseUserPermissions struct {
	Error   bool                   `json:"error"`
	Message string                 `json:"msg"`
	Data    jenkinsUserPermissions `json:"data"`
}

// jenkinsAdapter wraps the Jenkins client, enabling additional functionality
type jenkinsAdapter struct {
	*jenkins.Jenkins
}

// Config is the set of parameters needed to configure the Jenkins provider.
type Config struct {
	ServerURL string
	CACert    io.Reader
	Username  string
	Password  string
}

func newJenkinsClient(c *Config) *jenkinsAdapter {
	client := jenkins.CreateJenkins(nil, c.ServerURL, c.Username, c.Password)
	if c.CACert != nil {
		// provide CA certificate if server is using self-signed certificate
		client.Requester.CACert, _ = ioutil.ReadAll(c.CACert)
	}

	// return the Jenkins API client
	return &jenkinsAdapter{Jenkins: client}
}

func (j *jenkinsAdapter) GetLocalUser(username string) (jenkinsLocalUser, error) {
	var command bytes.Buffer
	commandTemplate := template.Must(template.New("command").Parse(getLocalUserCommand))

	err := commandTemplate.Execute(&command, jenkinsLocalUser{Username: username})
	if err != nil {
		return jenkinsLocalUser{}, fmt.Errorf("Failed parsing groovy commands to get local user: %v", err)
	}

	response := jenkinsResponse{}
	var respStruct interface{} = &response

	j.PostScript(command, respStruct)

	if response.Error {
		return jenkinsLocalUser{}, fmt.Errorf(response.Message)
	}

	return response.Data, nil
}

func (j *jenkinsAdapter) CreateLocalUser(username string, password string, fullname string, email string, description string) error {
	var command bytes.Buffer
	commandTemplate := template.Must(template.New("command").Parse(createLocalUserCommand))
	data := jenkinsLocalUserCreate{
		Password: password,
		jenkinsLocalUser: jenkinsLocalUser{
			Username:    username,
			Fullname:    fullname,
			Email:       email,
			Description: description,
		},
	}

	err := commandTemplate.Execute(&command, data)
	if err != nil {
		return fmt.Errorf("Failed parsing groovy commands to get local user: %v", err)
	}

	response := jenkinsResponse{}
	var respStruct interface{} = &response
	j.PostScript(command, respStruct)

	if response.Error {
		return fmt.Errorf(response.Message)
	}

	return nil
}

func (j *jenkinsAdapter) DeleteLocalUser(username string) error {
	var command bytes.Buffer
	commandTemplate := template.Must(template.New("command").Parse(deleteLocalUserCommand))
	err := commandTemplate.Execute(&command, jenkinsLocalUser{Username: username})
	if err != nil {
		return fmt.Errorf("Failed parsing groovy commands to get local user: %v", err)
	}

	response := jenkinsResponse{}
	var respStruct interface{} = &response

	j.PostScript(command, respStruct)

	if response.Error {
		return fmt.Errorf(response.Message)
	}

	return nil
}

func (j *jenkinsAdapter) GetUserPermissions(username string) (jenkinsUserPermissions, error) {
	var command bytes.Buffer

	commandTemplate := template.Must(template.New("command").Parse(getUserPermissionsCommand))
	err := commandTemplate.Execute(&command, jenkinsUserPermissions{Username: username})
	if err != nil {
		return jenkinsUserPermissions{}, fmt.Errorf("Error parsing groovy commands to get user permissions: %v", err)
	}

	response := jenkinsResponseUserPermissions{}
	var respStruct interface{} = &response

	j.PostScript(command, respStruct)
	if response.Error {
		return jenkinsUserPermissions{}, fmt.Errorf(response.Message)
	}

	return response.Data, nil
}

func (j *jenkinsAdapter) CreateUserPermissions(username string, permissions []string) error {
	var command bytes.Buffer

	commandTemplate := template.Must(template.New("command").Parse(createUserPermissionsCommand))
	err := commandTemplate.Execute(&command, jenkinsUserPermissions{Username: username, Permissions: permissions})
	if err != nil {
		return fmt.Errorf("Error parsing groovy commands to create user permissions: %v", err)
	}

	response := jenkinsResponseUserPermissions{}
	var respStruct interface{} = &response

	j.PostScript(command, respStruct)
	if response.Error {
		return fmt.Errorf(response.Message)
	}

	return nil
}

func (j *jenkinsAdapter) PostScript(payload bytes.Buffer, respStruct interface{}) error {
	finalPayload := url.Values{}
	finalPayload.Set("script", payload.String())

	resp, err := j.Requester.Post("/scriptText", strings.NewReader(finalPayload.Encode()), respStruct, map[string]string{})
	if err != nil {
		return fmt.Errorf("Error making request to Jenkins: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Call to jenkins return non 200 response code: %d, %v", resp.StatusCode, resp)
	}

	return nil
}
