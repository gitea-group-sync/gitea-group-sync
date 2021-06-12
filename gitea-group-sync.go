package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/ldap.v3"
	"gopkg.in/yaml.v2"
)

func AddUsersToTeam(apiKeys GiteaKeys, users []Account, team int) bool {

	for i := 0; i < len(users); i++ {

		fullusername := url.PathEscape(fmt.Sprintf("%s", users[i].Full_name))
		apiKeys.Command = "/api/v1/users/search?q=" + fullusername + "&access_token="
		foundUsers := RequestSearchResults(apiKeys)

		for j := 0; j < len(foundUsers.Data); j++ {

			if strings.EqualFold(users[i].Login, foundUsers.Data[j].Login) {
				apiKeys.Command = "/api/v1/teams/" + fmt.Sprintf("%d", team) + "/members/" + foundUsers.Data[j].Login + "?access_token="
				error := RequestPut(apiKeys)
				if len(error) > 0 {
					log.Println("Error (Team does not exist or Not Found User) :", parseJson(error).(map[string]interface{})["message"])
				}
			}
		}
	}
	return true
}

func DelUsersFromTeam(apiKeys GiteaKeys, Users []Account, team int) bool {

	for i := 0; i < len(Users); i++ {

		apiKeys.Command = "/api/v1/users/search?uid=" + fmt.Sprintf("%d", Users[i].Id) + "&access_token="

		foundUser := RequestSearchResults(apiKeys)

		apiKeys.Command = "/api/v1/teams/" + fmt.Sprintf("%d", team) + "/members/" + foundUser.Data[0].Login + "?access_token="
		RequestDel(apiKeys)
	}
	return true
}

var configFlag = flag.String("config", "config.yaml", "Specify YAML Configuration File")

func main() {
	// Parse flags of programm
	flag.Parse()
	mainJob() // First run for check settings

	var repTime string
	if len(os.Getenv("REP_TIME")) == 0 {

	} else {
		repTime = os.Getenv("REP_TIME")
	}

	c := cron.New()
	c.AddFunc(repTime, mainJob)
	c.Start()
	fmt.Println(c.Entries())
	for true {
		time.Sleep(100 * time.Second)
	}
}

// This Function parses the enviroment for application specific variables and returns a Config struct.
// Used for setting all required settings in the application
func importEnvVars() Config {

	// Create temporary structs for creating the final config
	envConfig := Config{}

	// ApiKeys
	envConfig.ApiKeys = GiteaKeys{}
	envConfig.ApiKeys.TokenKey = strings.Split(os.Getenv("GITEA_TOKEN"), ",")
	envConfig.ApiKeys.BaseUrl = os.Getenv("GITEA_URL")

	// LDAP Config
	envConfig.LdapURL = os.Getenv("LDAP_URL")
	envConfig.LdapBindDN = os.Getenv("BIND_DN")
	envConfig.LdapBindPassword = os.Getenv("BIND_PASSWORD")
	envConfig.LdapFilter = os.Getenv("LDAP_FILTER")
	envConfig.LdapUserSearchBase = os.Getenv("LDAP_USER_SEARCH_BASE")

	// Check TLS Settings
	if len(os.Getenv("LDAP_TLS_PORT")) > 0 {
		port, err := strconv.Atoi(os.Getenv("LDAP_TLS_PORT"))
		envConfig.LdapPort = port
		envConfig.LdapTLS = true
		log.Printf("DialTLS:=%v:%d", envConfig.LdapURL, envConfig.LdapPort)
		if err != nil {
			log.Println("LDAP_TLS_PORT is invalid.")
		}
	} else {
		if len(os.Getenv("LDAP_PORT")) > 0 {
			port, err := strconv.Atoi(os.Getenv("LDAP_PORT"))
			envConfig.LdapPort = port
			envConfig.LdapTLS = false
			log.Printf("Dial:=%v:%d", envConfig.LdapURL, envConfig.LdapPort)
			if err != nil {
				log.Println("LDAP_PORT is invalid.")
			}
		}
	}
	// Set defaults for user Attributes
	if len(os.Getenv("LDAP_USER_IDENTITY_ATTRIBUTE")) == 0 {
		envConfig.LdapUserIdentityAttribute = "uid"
		log.Println("By default LDAP_USER_IDENTITY_ATTRIBUTE = 'uid'")
	} else {
		envConfig.LdapUserIdentityAttribute = os.Getenv("LDAP_USER_IDENTITY_ATTRIBUTE")
	}

	if len(os.Getenv("LDAP_USER_FULL_NAME")) == 0 {
		envConfig.LdapUserFullName = "sn" //change to cn if you need it
		log.Println("By default LDAP_USER_FULL_NAME = 'sn'")
	} else {
		envConfig.LdapUserFullName = os.Getenv("LDAP_USER_FULL_NAME")
	}

	return envConfig // return the config struct for use.
}

func importYAMLConfig(path string) (Config, error) {
	// Open Config File
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err // Aborting
	}
	defer f.Close()

	// Parse File into Config Struct
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err // Aborting
	}
	return cfg, nil
}

func (c Config) checkConfig() {
	if len(c.ApiKeys.TokenKey) <= 0 {
		log.Println("GITEA_TOKEN is empty or invalid.")
	}
	if len(c.ApiKeys.BaseUrl) == 0 {
		log.Println("GITEA_URL is empty")
	}
	if len(c.LdapURL) == 0 {
		log.Println("LDAP_URL is empty")
	}
	if c.LdapPort <= 0 {
		log.Println("LDAP_TLS_PORT is invalid.")
	} else {
		log.Printf("DialTLS:=%v:%d", c.LdapURL, c.LdapPort)
	}
	if (len(c.LdapBindDN) > 0 && len(c.LdapBindPassword) == 0) {
		log.Println("BIND_DN supplied, but BIND_PASSWORD is empty")
	}
	if len(c.LdapFilter) == 0 {
		log.Println("LDAP_FILTER is empty")
	}
	if len(c.LdapUserSearchBase) == 0 {
		log.Println("LDAP_USER_SEARCH_BASE is empty")
	}
	if len(c.LdapUserIdentityAttribute) == 0 {
		c.LdapUserIdentityAttribute = "uid"
		log.Println("By default LDAP_USER_IDENTITY_ATTRIBUTE = 'uid'")
	}
	if len(c.LdapUserFullName) == 0 {
		c.LdapUserFullName = "sn" //change to cn if you need it
		log.Println("By default LDAP_USER_FULL_NAME = 'sn'")
	}
}

func mainJob() {

	//------------------------------
	//  Check and Set input settings
	//------------------------------
	var cfg Config

	cfg, importErr := importYAMLConfig(*configFlag)
	if importErr != nil {
		log.Println("Fallback: Importing Settings from Enviroment Variables ")
		cfg = importEnvVars()
	} else {
		log.Println("Successfully imported YAML Config from " + *configFlag)
		fmt.Println(cfg)
	}
	// Checks Config
	cfg.checkConfig()
	log.Println("Checked config elements")

	// Prepare LDAP Connection
	var l *ldap.Conn
	var err error
	if cfg.LdapTLS {
		l, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", cfg.LdapURL, cfg.LdapPort), &tls.Config{InsecureSkipVerify: true})
	} else {
		l, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", cfg.LdapURL, cfg.LdapPort))
	}

	if err != nil {
		fmt.Println(err)
		fmt.Println("Please set the correct values for all specifics.")
		os.Exit(1)
	}
	defer l.Close()

        if len(cfg.LdapBindDN) == 0 {
	   err = l.UnauthenticatedBind("")
	} else {
	  err = l.Bind(cfg.LdapBindDN, cfg.LdapBindPassword)
	}
	
	if err != nil {
		log.Fatal(err)
	}

	page := 1
	cfg.ApiKeys.BruteforceTokenKey = 0
	cfg.ApiKeys.Command = "/api/v1/admin/orgs?page=" + fmt.Sprintf("%d", page) + "&limit=20&access_token=" // List all organizations
	organizationList := RequestOrganizationList(cfg.ApiKeys)

	log.Printf("%d organizations were found on the server: %s", len(organizationList), cfg.ApiKeys.BaseUrl)

	for 1 < len(organizationList) {

		for i := 0; i < len(organizationList); i++ {

			log.Println(organizationList)

			log.Printf("Begin an organization review: OrganizationName= %v, OrganizationId= %d \n", organizationList[i].Name, organizationList[i].Id)

			cfg.ApiKeys.Command = "/api/v1/orgs/" + organizationList[i].Name + "/teams?access_token="
			teamList := RequestTeamList(cfg.ApiKeys)
			log.Printf("%d teams were found in %s organization", len(teamList), organizationList[i].Name)
			log.Printf("Skip synchronization in the Owners team")
			cfg.ApiKeys.BruteforceTokenKey = 0

			for j := 1; j < len(teamList); j++ {

				// preparing request to ldap server
				filter := fmt.Sprintf(cfg.LdapFilter, teamList[j].Name)
				searchRequest := ldap.NewSearchRequest(
					cfg.LdapUserSearchBase, // The base dn to search
					ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
					filter, // The filter to apply
					[]string{"cn", "uid", "mailPrimaryAddress, sn", cfg.LdapUserIdentityAttribute}, // A list attributes to retrieve
					nil,
				)
				// make request to ldap server
				sr, err := l.Search(searchRequest)
				if err != nil {
					log.Fatal(err)
				}
				AccountsLdap := make(map[string]Account)
				AccountsGitea := make(map[string]Account)
				var addUserToTeamList, delUserToTeamlist []Account
				if len(sr.Entries) > 0 {
					log.Printf("The LDAP %s has %d users corresponding to team %s", cfg.LdapURL, len(sr.Entries), teamList[j].Name)
					for _, entry := range sr.Entries {

						AccountsLdap[entry.GetAttributeValue(cfg.LdapUserIdentityAttribute)] = Account{
							Full_name: entry.GetAttributeValue(cfg.LdapUserFullName),
							Login:     entry.GetAttributeValue(cfg.LdapUserIdentityAttribute),
						}
					}

					cfg.ApiKeys.Command = "/api/v1/teams/" + fmt.Sprintf("%d", teamList[j].Id) + "/members?access_token="
					AccountsGitea, cfg.ApiKeys.BruteforceTokenKey = RequestUsersList(cfg.ApiKeys)
					log.Printf("The gitea %s has %d users corresponding to team %s Teamid=%d", cfg.ApiKeys.BaseUrl, len(AccountsGitea), teamList[j].Name, teamList[j].Id)

					for k, v := range AccountsLdap {
						if AccountsGitea[k].Login != v.Login {
							addUserToTeamList = append(addUserToTeamList, v)
						}
					}
					log.Printf("can be added users list %v", addUserToTeamList)
					AddUsersToTeam(cfg.ApiKeys, addUserToTeamList, teamList[j].Id)

					for k, v := range AccountsGitea {
						if AccountsLdap[k].Login != v.Login {
							delUserToTeamlist = append(delUserToTeamlist, v)
						}
					}
					log.Printf("must be del users list %v", delUserToTeamlist)
					DelUsersFromTeam(cfg.ApiKeys, delUserToTeamlist, teamList[j].Id)

				} else {
					log.Printf("The LDAP %s not found users corresponding to team %s", cfg.LdapURL, teamList[j].Name)
				}
			}
		}

		page++
		cfg.ApiKeys.BruteforceTokenKey = 0
		cfg.ApiKeys.Command = "/api/v1/admin/orgs?page=" + fmt.Sprintf("%d", page) + "&limit=20&access_token=" // List all organizations
		organizationList = RequestOrganizationList(cfg.ApiKeys)
		log.Printf("%d organizations were found on the server: %s", len(organizationList), cfg.ApiKeys.BaseUrl)
	}
}
