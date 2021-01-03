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

const getLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import hudson.security.HudsonPrivateSecurityRealm.Details
import hudson.tasks.Mailer
import groovy.json.JsonOutput

def result = [:]

def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

user = secRealm.getUser('{{ .Username }}')
if (user != null) {
  	result['error'] = false
  	result['msg'] = ''
	result['data'] = [:]
  	result['data']['username'] = user.getId()
	result['data']['fullname'] = user.getFullName()
  	result['data']['password_hash'] = user.getProperty(Details.class).getPassword()
  	result['data']['email'] = user.getProperty(Mailer.UserProperty.class).getAddress()
  	result['data']['description'] = user.getDescription() != null ? user.getDescription() : ''
} else {
	result['error'] = false
  	result['msg'] = ''
  	result['data'] = [:]
}

return println(JsonOutput.toJson(result))
`

const createLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import hudson.tasks.Mailer
import groovy.json.JsonOutput

def result = [:]
def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

def user = secRealm.createAccount('{{ .Username }}', '{{ .Password }}')
user.addProperty(new Mailer.UserProperty('{{ .Email }}'))
user.setFullName('{{ .Fullname }}')
user.setDescription('{{ .Description }}')
result['error'] = false
result['msg'] = 'User {{ .Username }} successfully created'
result['data'] = [:]

return println(JsonOutput.toJson(result))
`

const deleteLocalUserCommand = `
import hudson.security.HudsonPrivateSecurityRealm
import groovy.json.JsonOutput

def result = [:]
def secRealm = jenkins.model.Jenkins.instance.getSecurityRealm()
if (!secRealm instanceof HudsonPrivateSecurityRealm) {
  result['error'] = true
  result['msg'] = 'Jenkins is not using local user database'
  result['data'] = [:]
  return println(JsonOutput.toJson(result))
}

user = secRealm.getUser('{{ .Username }}')
user.delete()
result['error'] = false
result['msg'] = 'User {{ .Username }} successfully created'
result['data'] = [:]

return println(JsonOutput.toJson(result))
`

type jenkinsClient interface {
	GetLocalUser(username string) (jenkinsLocalUser, error)
	CreateLocalUser(username string, password string, fullname string, email string, description string) error
	DeleteLocalUser(username string) error
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
	Error   bool   `json:"error"`
	Message string `json:"msg"`
	Data    jenkinsLocalUser
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
