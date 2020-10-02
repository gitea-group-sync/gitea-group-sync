package main

type Organization struct {
	Id          int    `json:"id"`
	Avatar_url  string `json:"avatar_url"`
	Description string `json:"description"`
	Full_name   string `json:"full_name"`
	Location    string `json:"location"`
	Name        string `json:"username"`
	Visibility  string `json:"visibility"`
	Website     string `json:"website"`
}

type Team struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permission  string `json:"permission"`
}
type User struct {
	Id         int    `json:"id"`
	Avatar_url string `json:"avatar_url"`
	Created    string `json:"created"`
	Email      string `json:"email"`
	Full_name  string `json:"full_name"`
	Is_admin   bool   `json:"is_admin"`
	Language   string `json:"language"`
	Last_login string `json:"last_login"`
	Login      string `json:"login"`
}

type Account struct {
	Id        int    `json:"id"`
	Full_name string `json:"full_name"`
	Login     string `json:"login"`
}

type SearchResults struct {
	Data []User `json:"data"`
	Ok   bool   `json:"ok"`
}

type GiteaKeys struct {
	TokenKey           []string `yaml:"TokenKey"`
	BaseUrl            string   `yaml:"BaseUrl"`
	Command            string
	BruteforceTokenKey int
}

// Config describes the settings of the application. This structure is used in the settings-import process
type Config struct {
	ApiKeys                   GiteaKeys `yaml:"ApiKeys"`
	LdapURL                   string    `yaml:"LdapURL"`
	LdapPort                  int       `yaml:"LdapPort"`
	LdapTLS                   bool      `yaml:"LdapTLS"`
	LdapBindDN                string    `yaml:"LdapBindDN"`
	LdapBindPassword          string    `yaml:"LdapBindPassword"`
	LdapFilter                string    `yaml:"LdapFilter"`
	LdapUserSearchBase        string    `yaml:"LdapUserSearchBase"`
	ReqTime                   string    `yaml:"ReqTime"`
	LdapUserIdentityAttribute string    `yaml:"LdapUserIdentityAttribute"`
	LdapUserFullName          string    `yaml:"LdapUserFullName"`
} //!TODO! Implement check if valid
