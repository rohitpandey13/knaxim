package main

var helpstr = `knaximctl is for admin controls of a knaxim server
command form: knaximctl [Options]... [Action] [Arguments]...

Actions

addUser
userInfo
addRole
removeRole
updateFileCount
updateFileSpace
initDB
addAcronyms
help
`

var helpstrs = map[string]string{
	"help":            "display help message of actions\nknaximctl help",
	"addrole":         "adds a role to a user\nknaximctl addRole [username] [role]",
	"removerole":      "removes a role to a user\nknaximctl removeRole [username] [role]",
	"updatefilecount": "change the file count limit for a user\nknaimxctl updateFileCount [username] [amount]",
	"updatefilespace": "change the file space limit for a user\nknaximctl updateFileSpace [username] [amount]",
	"initdb":          "initialize database\nknaximctl initDB",
	"addacronyms":     "add acyonyms to database\nknaximctl addAcronyms [filepath]",
	"adduser":         "add user to database\nknaximctl addUser [username] [email] [password,optional]",
	"userinfo":        "display information about a user\nknaximctl userInfo [username]",
}