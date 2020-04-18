#!/bin/sh
export GITEA_TOKEN=c00c810bb668c63ce7cd8057411d2f560eac469c,2c02df6959d012dee8f5da3539f63223417c4bbe
export GITEA_URL=http://localhost:3000
export LDAP_URL=localhost
export LDAP_TLS_PORT=636
export BIND_DN='cn=admin,dc=planetexpress,dc=com'
export BIND_PASSWORD=GoodNewsEveryone
export LDAP_FILTER='(&(objectClass=person)(memberOf=cn=%s,ou=people,dc=planetexpress,dc=com))'
export LDAP_USER_SEARCH_BASE='ou=people,dc=planetexpress,dc=com'
export LDAP_USER_LOGIN_ATTRIBUTE='uid'
export REP_TIME='@every 1m'
go run .
