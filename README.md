# Project overview:

1.  Problem:
    There is a lack of bots/assistants that can save links/urls in discord. 

2.  Audience: Anyone looking to save links in discord. (Other social media might be included later)

3.  MVP Solution:

    1. Create a discord bot that saves, updates, and deletes links. 
    2. A website to show the bot's functionality along with its invitiation link.

4.  Tech stack

    1. Frontend: HTML (HTMX is planned to be added later)
    2. Backend: Golang
    3. Database: Supabase, easy to use serverless postgres.

5.  Project management:
    1. Trello for task management

# How to run:
1. Clone the repo:
```
git clone https://github.com/yohanc3/link-vault
```
2. Create and populate a .env file:
```
POSTGRES_DB_URI= user=<INSERT> password=<INSERT> host=<INSERT> port=<INSERT> dbname=postgres

DISCORD_BOT_TOKEN=
DISCORD_TEST_BOT_TOKEN=

#DEV or PROD
#If MODE is set to DEV, the DISCORD_TEST_BOT_TOKEN will be used. 
MODE=

```
3. Run main.go - (Dependencies will be automatically installed)
```
go run .
```
