# Wikipedia Edit Tracker

This program tracks recent Wikipedia edits by monitoring the Wikimedia event stream. It stores edit counts per language and date in an SQLite database and allows users to set their preferred language.

## Features

- Tracks Wikipedia edits in real-time.
- Stores edit counts by language and date.
- Allows users to set and retrieve their preferred language.
- Uses SQLite for lightweight, local data storage.
- Can be integrated into a Discord server as a bot.
- Output the most recent five updates for requested language and date.

## Requirements

- Go 1.18+
- SQLite3
- Discord Bot Token (for Discord integration)

## Installation

1. **Clone the repository**
   ```sh
   git clone https://github.com/dilya-gitit/wiki-tracker
   cd wiki-tracker
   ```
2. **Install dependencies**
   ```sh
   go mod tidy
   ```

## Usage

### Running the Program

1. **Compile and run the application**

   ```sh
   go run .
   ```

2. **Tracking Wikipedia Edits**

   - The program will listen to Wikimedia's event stream and store edit counts in `wikipedia_bot.db`.

### Launching the Bot in a Discord Server

1. **Create a Discord Bot**
   - Go to the [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application and add a bot.
   - Copy the bot token.
2. **Add the Bot to a Server**
   - Generate an OAuth2 invite link with `bot` and `messages.read` permissions.
   - Invite the bot to your Discord server.
3. **Run the Discord Bot**
   - Set environment variable "DISCORD_BOT_TOKEN" to the bot token.
   - Start the bot with:
     ```sh
     go run .
     ```

### Discord Bot Commands

#### `!recent`
Retrieves the most recent changes for the current or specified language.

#### `!setLang [language_code]`
Sets a default language for the user/server session.

#### `!stats [yyyy-mm-dd]`
Returns the number of edits per day for the chosen language.

## CI/CD
1. Currently, the repository has a basic workflow to run Go tests and build on every push and PR to main.

## Design Decisions and Trade-Offs

1. The app reads from the streaming API and updates the count of messages read for a date-language pair in memory. However, every 100 messages, we update the count in the database. This ensures that we are not running expensive SQL queries on every single message, while also providing a decent persistence of statistical data in case the app crashes and has to reboot.
2. When the app restarts, it will bootstrap the count from the database, save it into memory, and continue counting. When a user asks for the stats, it will return immediately from an in-memory map.

### Additional Feature
- The app also stores the latest Kafka offset read for each language-date pair, which is part of the payload that Wikipedia's streaming API provides. This was initially done to avoid counting the same message twice. However, later it was realized that the app always reads the latest message when it restarts, making this feature unnecessary for now, but it can be utilized if conditions change.

### Database Structure

The SQLite database contains two tables:

#### `stats`

| Column | Type    | Description                             |
| ------ | ------- | --------------------------------------- |
| date   | TEXT    | Date of the edit (YYYY-MM-DD)           |
| lang   | TEXT    | Wikipedia language (e.g., `en`, `fr`)   |
| edits  | INTEGER | Number of edits for the date & language |
| offset | INTEGER | Latest offset read for date-lang pair   |

#### `user_lang`

| Column  | Type | Description        |
| ------- | ---- | ------------------ |
| user_id | TEXT | Unique user ID     |
| lang    | TEXT | Preferred language |

## Ideas on Scaling Using Kafka

1. If we need to perform distributed processing of these messages from the stream, we could design the following architecture:

### Kafka Producer
- Have a single Kafka producer that performs a lightweight operation of reading from the stream and publishing this data into an internal Kafka topic.
- Use multiple partitions for the topic and make the producer do Round-Robin publishing for those partitions.

### Kafka Consumer
- Have multiple consumers depending on the load and assign them to different partitions.
- These consumers can now parallelize the processing of the messages.
