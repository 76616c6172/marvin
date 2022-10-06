![alt text](https://github.com/76616c6172/marvin/blob/master/marvin.jpg)

Discord bot for Ersatz.


## Setup guide
**Step 1** Get a token so the bot can interact with the discord API

Go to:
[https://discord.com/developers/applications/](https://discord.com/developers/applications/)

- Create a new named "application" and save the `APPLICATION_ID`
- Click the Bot tab and add a named bot
- Reset the `TOKEN` for the bot and save it
- Uncheck the slider for "PUBLIC BOT"
- Check the slider for "PRESENCE INTENT"
- Check the slider for "SERVER MEMBERS INTENT"
- Check the slider for "MESSAGE CONTENT INTENT"


**Step 2** Invite the bot to your server

Modify this link with the `APPLICATION_ID` from step 1.  
`https://discord.com/oauth2/authorize?client_id=<APPLICATION_ID>&permissions=8&scope=bot`  

Then copy paste the in the browser and invite the bot to your server.

**Step 3** Compile the bot

[Install Go](https://go.dev/doc/install) if you haven't already.

Modify the code to your liking and compile the bot

```bash
./build_bot.sh
```

**Step 4** Add the authorization token

Edit run_bot.sh and add the `TOKEN` from step 1 and Execute run_bot.sh
```bash
./run.sh
```