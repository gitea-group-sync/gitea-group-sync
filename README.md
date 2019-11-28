# Gitea-group-sync

This application adds users to the appropriate groups. 

You must have configured your LDAP with gitea

Here I will give the settings for a simple [LDAP](https://github.com/rroemhild/docker-test-openldap), you can configure by changing the code as you like

If you configured the gitea <=> [LDAP](https://github.com/rroemhild/docker-test-openldap) connection correctly, you should have users

![](images/Image1.png)

You need to create Manage Access Tokens and add key to run.sh or docker-compose.yml the configuration file

The application supports several keys, since to add people to the group you must be the owner of the organization.

![](images/Image2.png)