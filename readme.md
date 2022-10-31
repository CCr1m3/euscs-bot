# todo


- Add queue time into match quality calculations. 
- first check commands before adding it (use a flag to force update)
- Punish dodging people in some way (reduce match priority by reducing quality of matches with dodging players?)
- /autojoin (rejoin queue after match if match was not canceled)
- improve messages fetching (goroutine per channel ?), keep messages log and try to built missing from there instead of restarting from scratch
- chat message improvements (using neural instead of markov to be able to get contextual informations ?)
- /givecredits @user 100 for mods
- better prod tools (restarting, log access)
- ~~Pinned message on dedicated channel with howto use bot~~
- ~~#request channel for bot talk and commands~~
- ~~Automatically leave queue after a time (1h?). Probably tag the person who was forced out.~~ 
- ~~Currency system. Get currency from winning (more, 20?), losing (less, 10?) and predicting correctly (5?). Pay with currency for @ Ai.Mi (20). (/currency)~~
- ~~Match predictions. Disallow people playing in a match from predicting.~~
- ~~update bot permission link in this file~~
- ~~/result reserved to people in the match only (probably reserved to loosing team too)~~
- ~~/result integer checking for winning state ( score between teams >=2)~~
- ~~better matchmaking algorithms (now it takes 6 players randomly in the pool, and is often first come first serve since it's matching every 15 secs)~~
- ~~rating update (elo based)~~
- ~~rating based matchmaking~~
- ~~Better interfacing between slashcommands -> handlers -> db (now its a bit messy)~~
- ~~/cancel instead of /result 0 0   (in case someone is not here)~~

# bot permissions on discord dev portal
![image.png](image.png)
url : https://discord.com/api/oauth2/authorize?client_id={clientid}&permissions=397552987216&scope=applications.commands%20bot
