package mock

import (
	"github.com/hashicorp/vault/api"
	"os"
	"ritchie-server/server"
	"ritchie-server/server/config"
)

const (
	keycloakUrl = "KEYCLOAK_URL"
	oauthUrl = "OAUTH_URL"
	cliVersionUrl = "CLI_VERSION_URL"
	remoteUrl =  "REMOTE_URL"
)

func DummyConfig(args ...string) server.Config {
	return config.Configuration{
		Configs:DummyConfigMap(args...),
		SecurityConstraints: server.SecurityConstraints{
			Constraints: []server.DenyMatcher{{
				Pattern:      "/validate",
				RoleMappings: map[string][]string{"admin": {"POST", "GET"}},
			}},
			PublicConstraints: []server.PermitMatcher{{
				Pattern: "/public",
				Methods: []string{"POST", "GET"},
			}},
		},
	}
}

func DummyConfigMap(args ...string) map[string]*server.ConfigFile {
	keycloakUrl := getEnv(keycloakUrl, "http://localhost:8080")
	remoteUrl := getEnv(remoteUrl, "http://localhost:8882")
	realm := "ritchie"
	clientId := "user-login"
	clientSecret := "user-login"
	if len(args) > 0 && args[0] != "" {
		keycloakUrl = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		realm = args[1]
	}
	if len(args) > 2 && args[2] != "" {
		clientId = args[2]
	}
	if len(args) > 3 && args[3] != "" {
		clientSecret = args[3]
	}
	return map[string]*server.ConfigFile{
		"zup": {
			KeycloakConfig: &server.KeycloakConfig{
				Url:          keycloakUrl,
				Realm:        realm,
				ClientId:     clientId,
				ClientSecret: clientSecret,
			},
			OauthConfig: &server.OauthConfig{
				Url:      getEnv(oauthUrl, "http://localhost:8080/auth/realms/ritchie"),
				ClientId: "oauth",
			},
			CredentialConfig: map[string][]server.CredentialConfig{
				"credential1": {{Field: "Field", Type: "type"}},
				"credential2": {{Field: "field2", Type: "type"}},
			},
			CliVersionConfig: server.CliVersionConfig{
				Url:      getEnv(cliVersionUrl, "http://localhost:8882/s3-version-mock"),
				Provider: "s3",
			},
			RepositoryConfig: []server.Repository{
				{
					Name:           "commons",
					Priority:       0,
					TreePath:       "/tree/tree.json",
					Remote:         remoteUrl,
					ServerUrl:      "http://localhost:3000",
					ReplaceRepoUrl: "http://localhost:3000/formulas",
					Username:       "",
					Password:       "",
				},
				{
					Name:           "test1",
					Priority:       1,
					TreePath:       "/tree/tree-test1.json",
					Remote:         remoteUrl,
					ServerUrl:      "http://localhost:3000",
					ReplaceRepoUrl: "http://localhost:3000/formulas",
					Username:       "",
					Password:       "",
				},
			},
		}}
}

//Cli Version
func DummyConfigCliVersionUrlNotFound() server.Config {
	return config.Configuration{
		Configs: map[string]*server.ConfigFile{
			"zup": {
				CliVersionConfig: server.CliVersionConfig{
					Provider: "s3",
				},
			}},
	}
}
func DummyConfigCliVersionUrlWrong() server.Config {
	return config.Configuration{
		Configs: map[string]*server.ConfigFile{
			"zup": {
				CliVersionConfig: server.CliVersionConfig{
					Url:      "wrong",
					Provider: "s3",
				},
			}},
	}
}
func DummySecurityConstraints() server.SecurityConstraints {
	return server.SecurityConstraints{
		Constraints: []server.DenyMatcher{{
			Pattern:      "/test",
			RoleMappings: map[string][]string{"user": {"POST", "GET"}},
		}},
		PublicConstraints: []server.PermitMatcher{{
			Pattern: "/public",
			Methods: []string{"POST", "GET"},
		}},
	}
}

//Credential
func DummyCredential() string {
	return `{
	"service": "credential1",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}
func DummyCredentialEmpty() string {
	return `{
	"username": "Ubijara",
	"service": "",
		"credential": {
		}
	}`
}
func DummyCredentialAdmin() string {
	return `{
	"username": "Ubijara",
	"service": "credential1",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}
func DummyCredentialBadRequest() string {
	return `{
	"service": "invalid",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}

func DummyRepo() server.Repository {
	remote := getEnv(remoteUrl, "http://localhost:8882")
	return server.Repository{
		Name:           "commons",
		Priority:       0,
		TreePath:       "/tree/tree.json",
		Remote:         remote,
		ServerUrl:      "http://localhost:3000",
		ReplaceRepoUrl: "http://localhost:3000/formulas",
		Username:       "",
		Password:       "",
	}
}

func DummyRepoList() []server.Repository {
	remote := getEnv(remoteUrl, "http://localhost:8882")
	return []server.Repository{
		{
			Name:           "commons",
			Priority:       0,
			TreePath:       "/tree/tree.json",
			Remote:         remote,
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Username:       "",
			Password:       "",
		},
		{
			Name:           "test1",
			Priority:       1,
			TreePath:       "/tree/tree-test1.json",
			Remote:         remote,
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Username:       "",
			Password:       "",
		},
		{
			Name:           "test2",
			Priority:       2,
			TreePath:       "/tree/tree-test2.json",
			Remote:         remote,
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Username:       "",
			Password:       "",
		},
	}
}

//server.KeycloakManager mock
type KeycloakMock struct {
	Token string
	Code  int
	Err   error
}

func (k KeycloakMock) CreateUser(server.CreateUser, string) (string, error) {
	return k.Token, k.Err
}
func (k KeycloakMock) DeleteUser(string, string) error {
	return k.Err
}
func (k KeycloakMock) Login(string, string, string) (string, int, error) {
	return k.Token, k.Code, k.Err
}

//server.ValtManager mock
type VaultMock struct {
	Err     error
	ErrList error
	Keys    []interface{}
}

func (v VaultMock) Write(string, map[string]interface{}) error {
	return v.Err
}
func (v VaultMock) Read(string) (map[string]interface{}, error) {
	return nil, v.Err
}
func (v VaultMock) List(string) ([]interface{}, error) {
	return v.Keys, v.ErrList
}
func (v VaultMock) Delete(string) error {
	return v.Err
}
func (v VaultMock) Start(*api.Client) {
}

type AuthorizationMock struct {
	B bool
	E error
	R []string
}

func (d AuthorizationMock) AuthorizationPath(bearerToken, path, method, org string) (bool, error) {
	return d.B, d.E
}
func (d AuthorizationMock) ValidatePublicConstraints(path, method string) bool {
	return d.B
}
func (d AuthorizationMock) ListRealmRoles(bearerToken, org string) ([]interface{}, error) {
	if d.E != nil {
		return nil, d.E
	}
	new := make([]interface{}, len(d.R))
	for i, v := range d.R {
		new[i] = v
	}
	return new, d.E
}


func getEnv(key, def string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return def
}
