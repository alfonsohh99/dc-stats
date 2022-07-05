# DC-STATS

DC-Stats (Discord Statistics) is a Discord bot written in go.

Using [discordgo](https://github.com/bwmarrin/discordgo) to interact with the discord API

Periodically gathers data from every guild and processes it.
Commands:

- [ !topVoice ] Users can see which one of them has spent the most on a voice channel using the command
- [ !myVoice ] Users can see their spent time in detail for every channel in the guild using the command
- [ !topMessage ] Users can see which one of them has sent more chat messages with the command
- [ !myMessage ] Users can see the number of messages sent per chat with the command

The bot makes use of the [chrono](https://github.com/procyon-projects/chrono) library for scheduling the background tasks:

- Gathering information about wether users are in voice chat or not, then adding up their total time
- Processing voice chat information so it does not have to be calculated for every command invocation
- Gathering chat message information from every chat channel and storing it
- Processing chat message informa so it does not have to be calculated for every command invocation
- Off-site backups to [AWS](https://github.com/aws/aws-sdk-go-v2) S3

Users can access their guild's data through the [Echo API](https://github.com/labstack/echo) which could also be fetched from a front-end

The bot also takes advantaje of goroutines for improved performance when managing asynchronous tasks.

All this data is stored using mongoDB.
