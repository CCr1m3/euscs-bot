# todo

- first check commands before adding it (use a flag to force update)
- Punish dodging people in some way (reduce match priority by reducing quality of matches with dodging players?)
- /autojoin (rejoin queue after match if match was not canceled)
- improve messages fetching (goroutine per channel ?), keep messages log and try to built missing from there instead of restarting from scratch
- chat message improvements (using neural instead of markov to be able to get contextual informations ?)
- /givecredits @user 100 for mods
- better prod tools (restarting, log access)

# bot permissions on discord dev portal
![perms.png](screens/perms.png)
url : https://discord.com/api/oauth2/authorize?client_id={clientid}&permissions=397552987216&scope=applications.commands%20bot

# features

- Chat response using markov chains

![chat.png](screens/chat.png)

- matchmaking (supporting predictions)

![matches.png](screens/matches.png)

- discord role based on in game rank

![rankup.png](screens/rankup.png)