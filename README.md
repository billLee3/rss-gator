### Setup Required To Run the Program
---
(*Linux*) 
Install postgres on your local machine using the following code:
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
If on linux then you will need to update your password, likeso:
`sudo passwd postgres`
(*Linux*) 
Install Golang on your machine using the following code:
```
sudo apt install golang-go
```
### Install RSS-Gator
---
Now that you have Golang installed on your machine use the following golang command from the terminal
```
go install https://github.com/billLee3/rss-gator@latest
```
The program should be on your machine now. 

### Setting up your local config file
---
In your home directory, create a `.gatorconfig.json` file from the terminal using the following bash command:
`touch .gatorconfig.json`
Within the `.gatorconfig.json` create json to hold your database url. Example below:
```
{
  db_url: "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```
### Migrate the DB
---
From rss-gator/sql/schema directory run:
`goose postgres <connection_string> up`
This will align the database structure to what is necessary to run the application
### Running the Program
---
From the root of the rss-gator directory (created from go install step), run `go build`.

Some different commands that you can run from the root of the rss-gator directory are as follows:
`go run . register <name>` - registers a user
`go run . login <name>` - logs in a user
`go run . users` - displays a list of all users registered to the db
`go run . addFeed <title> <url>` - adds a feed to the database
`go run . feeds` - displays all of the feeds being tracked in the database
`go run . agg <timeDuration: 1m>` - aggregates the data from the feeds in the feeds table
`go run . follow <url>` - allows users to follow certain feeds
`go run . following` - displays which feeds the logged in user is following
`go run . unfollow <url>` - unfollows a feed followed by the user
`go run . browse` - displays posts by the logged in user



