package jenkins

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	jenkins "github.com/bndr/gojenkins"
)

type jenkinsClient interface {
	GetLocalUser(username string) (jenkinsLocalUser, error)
}

type jenkinsLocalUser struct {
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	PasswordHash string `json:"password_hash"`
	Username     string `json:"username"`
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
	payload := [...]string{
		"import hudson.security.HudsonPrivateSecurityRealm.Details",
		"import hudson.tasks.Mailer",
		"import groovy.json.JsonOutput",
		"def response = [:]",
		fmt.Sprintf("def user = jenkins.model.Jenkins.instance.securityRealm.getUser('%s')", username),
		"if (user != null) {",
		`response["username"] = user.getId()`,
		`response["fullname"] = user.getFullName()`,
		`response["email"] = user.getProperty(Mailer.UserProperty.class).getAddress()`,
		`response["password_hash"] = user.getProperty(Details.class).getPassword()`,
		"}",
		"println(JsonOutput.toJson(response))",
	}
	finalPayload := url.Values{}
	finalPayload.Set("script", strings.Join(payload[:], "\n"))

	localUser := jenkinsLocalUser{}
	var respStruct interface{} = &localUser

	resp, err := j.Requester.Post("/scriptText", strings.NewReader(finalPayload.Encode()), respStruct, map[string]string{})

	if err != nil {
		return jenkinsLocalUser{}, fmt.Errorf("Error making request to Jenkins: %v", err)
	}

	if resp.StatusCode != 200 {
		return jenkinsLocalUser{}, fmt.Errorf("Call to jenkins return non 200 response code: %d, %v", resp.StatusCode, resp)
	}

	return localUser, nil
}
