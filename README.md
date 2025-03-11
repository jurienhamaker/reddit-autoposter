# reddit-autoposter
A simple go script to post reddit posts on a cron schedule

## First Steps
Go to your app preferences. Click the "Create app" or "Create another app" button. Fill out the form like so:

    name: My Example App
    App type: Choose the script option
    description: You can leave this blank
    about url: You can leave this blank
    redirect url: http://www.example.com/unused/redirect/uri (We won't be using this as a redirect)

Note: These examples will only work for script type apps, which will ONLY have access to accounts registered as "developers" of the app and require the application to know the user's password. Read more about app types.

Hit the "create app" button. Make note of the client ID and client secret. 

Note: You should NEVER post your client secret (or your reddit password) in public. If you create a bot, you should take steps to ensure that the bot's password and the app's client secret are secured against digital theft. The client IDs, secrets, tokens and passwords used here are, obviously, fake and invalid.

## Running the script 

- Copy the `example.config.toml` to `config.toml`
- Edit accordingly (You can set multiple cron jobs)
- Copy `.env.example` to `.env`
- Fill the `.env` with the values from the **First Steps** instructions
- Run `reddit-autoposter`
