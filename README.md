# DC-STATS

DC-Stats (Discord Statistics) is a Discord bot written in go.

Periodically gathers data from every guild and processes it for the users.
Users can see which one of them has spent the most on a voice channel using the !top command
Users can see their spent time in detail for every channel in the guild using the !myStats command

The bot makes use of the [github.com/procyon-projects/chrono](https://github.com/procyon-projects/chrono) library for scheduling the 2 main task

- Gathering information (frequently) wether users are in chat or not, then adding up the time
- Processing that information (less frequently) so it does not have to be calculated for every command invocation

The bot also takes advantaje of goroutines for improved performance managing asynchronous tasks
