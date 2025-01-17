package main

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command holds data related to a command executable by any message system
type Command struct {
	Function            func([]string, *CommandEnvironment) *discordgo.MessageEmbed //The function value of what to execute when the command is ran
	HelpText            string                                                      //The text that will display in the help message
	Arguments           []CommandArgument                                           //The arguments required for this command
	RequiredArguments   []string                                                    //The minimum required arguments by name that must exist for the function to execute; default = 0
	RequiredPermissions int64                                                       //The permission(s) a user must have for the command to be executed by them

	IsAlternateOf string //If this is an alternate command, point to the original command

	IsAdministrative bool //Whether or not this command requires the user to be a bot admin

	IsAdvancedCommand bool                                                                 //Whether or not this command uses advanced parameters
	AdvancedFunction  func([]CommandArgument, *CommandEnvironment) *discordgo.MessageEmbed //The function value of what to execute when the command is ran
}

// CommandArgument holds data related to an argument available or required by a command
type CommandArgument struct {
	//Used for help text
	Name        string //The name of the argument
	ArgType     string //The argument's type
	Description string //A description of the argument

	//Used for command argument parsing
	Value string //The value supplied with the argument
}

// CommandEnvironment holds data related to the environment a command can utilize for data or functionality
type CommandEnvironment struct {
	Channel *discordgo.Channel //The channel the command was executed in
	Guild   *discordgo.Guild   //The guild the command was executed in
	Message *discordgo.Message //The message that triggered the command execution
	User    *discordgo.User    //The user that executed the command
	Member  *discordgo.Member  //The guild member that executed the command

	Command   string //The command used to execute the command with this environment (in the event of a command alias)
	BotPrefix string //The bot prefix used to execute this command (useful for command lists and example commands)

	UpdatedMessageEvent bool
}

func initCommands() {
	//Initialize the commands map
	botData.Commands = make(map[string]*Command)

	//All user-accessible commands with no parameters
	botData.Commands["about"] = &Command{Function: commandAbout, HelpText: "Displays information about " + botData.BotName + " and how to use it."}
	botData.Commands["invite"] = &Command{Function: commandInvite, HelpText: "Displays available invite links for " + botData.BotName + "."}
	botData.Commands["source"] = &Command{Function: commandSource, HelpText: "Displays available source code links for " + botData.BotName + "."}
	botData.Commands["version"] = &Command{Function: commandVersion, HelpText: "Displays the current version of " + botData.BotName + "."}
	botData.Commands["credits"] = &Command{Function: commandCredits, HelpText: "Displays a list of credits for the creation and functionality of " + botData.BotName + "."}
	botData.Commands["roll"] = &Command{Function: commandRoll, HelpText: "Rolls a dice."}
	botData.Commands["doubleroll"] = &Command{Function: commandDoubleRoll, HelpText: "Rolls two die."}
	botData.Commands["coinflip"] = &Command{Function: commandCoinFlip, HelpText: "Flips a coin."}
	botData.Commands["join"] = &Command{Function: commandVoiceJoin, HelpText: "Joins the current voice channel.", RequiredPermissions: discordgo.PermissionVoiceConnect}
	botData.Commands["leave"] = &Command{Function: commandVoiceLeave, HelpText: "Leaves the current voice channel.", RequiredPermissions: discordgo.PermissionVoiceConnect}
	botData.Commands["ping"] = &Command{Function: commandPing, HelpText: "Returns the ping average to Discord."}

	//All user-accessible info commands with or without parameters
	botData.Commands["botinfo"] = &Command{Function: commandBotInfo, HelpText: "Displays info about the bot's current state."}
	botData.Commands["serverinfo"] = &Command{Function: commandServerInfo, HelpText: "Displays info about the current server."}
	botData.Commands["userinfo"] = &Command{
		Function: commandUserInfo,
		HelpText: "Displays info about the current or specified user.",
		Arguments: []CommandArgument{
			{Name: "user", Description: "The user to view info about", ArgType: "mention/user ID"},
		},
	}

	//All user-accessible commands with parameters
	botData.Commands["help"] = &Command{
		Function: commandHelp,
		HelpText: "Displays a list of commands you have permission to use.",
		Arguments: []CommandArgument{
			{Name: "page", Description: "The help page to view", ArgType: "number"},
			{Name: "command", Description: "The command to view help for", ArgType: "string"},
		},
	}
	botData.Commands["translate"] = &Command{
		Function: commandTranslate,
		HelpText: "Translates a given message to the specified language.",
		RequiredArguments: []string{
			"[source language]",
			"(target language) message",
		},
		Arguments: []CommandArgument{
			{Name: "source language", Description: "Optional, the language to translate from", ArgType: "string"},
			{Name: "target language", Description: "The language to translate to", ArgType: "string"},
			{Name: "message", Description: "The message to translate", ArgType: "string"},
		},
	}
	botData.Commands["nnid"] = &Command{
		Function: commandNNID,
		HelpText: "Checks whether the specified NNID exists or not.",
		RequiredArguments: []string{
			"username",
		},
		Arguments: []CommandArgument{
			{Name: "username", Description: "The NNID to check for", ArgType: "string"},
		},
	}
	botData.Commands["remind"] = &Command{
		Function: commandRemind,
		HelpText: "Reminds you with the written message at the specified time.",
		RequiredArguments: []string{
			"(message and time)/other",
		},
		Arguments: []CommandArgument{
			{Name: "message and time", Description: "The message to remind you with, including what time to remind you at", ArgType: "string"},
			{Name: "list", Description: "Lists your remind entries on an optionally specified page", ArgType: "<page>"},
			{Name: "remove", Description: "Deletes the specified remind entry or entries", ArgType: "number(s)"},
		},
	}
	botData.Commands["hewwo"] = &Command{
		Function: commandHewwo,
		HelpText: "Hewwo!!! (´・ω・｀)",
		RequiredArguments: []string{
			"message",
		},
		Arguments: []CommandArgument{
			{Name: "message", Description: "The text to translate to Hewwo", ArgType: "string"},
		},
	}
	botData.Commands["minecraft"] = &Command{
		Function: commandMinecraft,
		HelpText: "Displays information about a specified user or server.",
		RequiredArguments: []string{
			"user/server",
			"name/host",
		},
		Arguments: []CommandArgument{
			{Name: "user", Description: "Displays information about the specified user", ArgType: "string"},
			{Name: "server", Description: "Displays infromation about the specified server", ArgType: "ip(:port)"},
		},
	}
	botData.Commands["zalgo"] = &Command{
		Function: commandZalgo,
		HelpText: "Mystifies your text.",
		RequiredArguments: []string{
			"message",
		},
		Arguments: []CommandArgument{
			{Name: "message", Description: "The text to mystify", ArgType: "string"},
		},
	}
	botData.Commands["nlp"] = &Command{
		Function: commandNLP,
		HelpText: "Raw natural language processing in Discord. Powered by Prose:tm:.",
		RequiredArguments: []string{
			"message",
		},
		Arguments: []CommandArgument{
			{Name: "message", Description: "The message to parse", ArgType: "string"},
		},
	}
	botData.Commands["image"] = &Command{
		IsAdvancedCommand: true,
		AdvancedFunction:  commandImageAdv,
		HelpText:          "Allows you to manipulate images with various effects.",
		RequiredArguments: []string{
			"-effect (value)",
		},
		Arguments: []CommandArgument{
			{Name: "backgroundcolor", Description: "Sets the background color (if transparent; set before other effects)", ArgType: "#hex, rgb(), rgba()"},
			{Name: "brightness", Description: "Sets the brightness", ArgType: "percentage"},
			{Name: "contrast", Description: "Sets the contrast", ArgType: "percentage"},
			{Name: "flip", Description: "Flips the image", ArgType: "horizontal/vertical"},
			{Name: "gamma", Description: "Sets the gamma", ArgType: "percentage"},
			{Name: "gaussian", Description: "Applies the gaussian blur effect at the specified intensity", ArgType: "percentage"},
			{Name: "grayscale/greyscale", Description: "Applies a grayscale effect"},
			{Name: "height", Description: "Sets the height", ArgType: "number"},
			{Name: "interpolation", Description: "Sets the interpolation (set before other effects)", ArgType: "cubic, linear, nearest"},
			{Name: "invert", Description: "Inverts the colors"},
			{Name: "pixelate", Description: "Applies a pixelation effect at the specified intensity", ArgType: "number"},
			{Name: "resampling", Description: "Sets the resampling (set before other effects)", ArgType: "box, cubic, lanczos, linear, nearest"},
			{Name: "rotate", Description: "Rotates by the specified degrees", ArgType: "circular degrees"},
			{Name: "saturation", Description: "Applies a saturation effect at the specified intensity", ArgType: "percentage"},
			{Name: "sepia", Description: "Applies a sepia effect at the specified intensity", ArgType: "percentage"},
			{Name: "sobel", Description: "Applies the sobel filter"},
			{Name: "threshold", Description: "Applies a black/white threshold at the specified intensity", ArgType: "percentage"},
			{Name: "transpose", Description: "Flips the image horizontally and rotates it 90° counter-clockwise"},
			{Name: "transverse", Description: "Flips the image vertically and rotates it 90° counter-clockwise"},
			{Name: "width", Description: "Sets the width", ArgType: "number"},
		},
	}
	botData.Commands["screenshot"] = &Command{
		Function: commandScreenshot,
		HelpText: "Takes a screenshot of a website.",
		RequiredArguments: []string{
			"url",
		},
		Arguments: []CommandArgument{
			{Name: "url", Description: "The URL to take a screenshot of", ArgType: "url"},
		},
	}
	botData.Commands["cve"] = &Command{
		Function: commandCVE,
		HelpText: "Fetches information about a specified CVE.",
		RequiredArguments: []string{
			"CVE ID",
		},
		Arguments: []CommandArgument{
			{Name: "cve", Description: "The CVE ID to fetch information about", ArgType: "string"},
		},
	}
	botData.Commands["geoip"] = &Command{
		Function: commandGeoIP,
		HelpText: "Performs a GeoIP lookup on the specified IP/hostname.",
		RequiredArguments: []string{
			"IP/hostname",
		},
		Arguments: []CommandArgument{
			{Name: "IP/hostname", Description: "The IP or hostname to perform a GeoIP lookup on", ArgType: "IP address/hostname"},
		},
	}
	if botData.BotOptions.UseXKCD {
		botData.Commands["xkcd"] = &Command{
			Function: commandXKCD,
			HelpText: "Displays an XKCD comic depending on the requested type or comic number.",
			RequiredArguments: []string{
				"(comic number|latest|random)",
			},
			Arguments: []CommandArgument{
				{Name: "comic number", Description: "The number of an existing XKCD comic", ArgType: "number"},
				{Name: "latest", Description: "Fetches the latest XKCD comic", ArgType: ""},
				{Name: "random", Description: "Fetches a random XKCD comic", ArgType: ""},
			},
		}
	}
	if botData.BotOptions.UseImgur {
		botData.Commands["imgur"] = &Command{
			Function: commandImgur,
			HelpText: "Displays info about the specified Imgur image or album URL.",
			RequiredArguments: []string{
				"url",
			},
			Arguments: []CommandArgument{
				{Name: "url", Description: "The Imgur image or album URL", ArgType: "string"},
			},
		}
	}
	if botData.BotOptions.UseGitHub {
		botData.Commands["github"] = &Command{
			Function: commandGitHub,
			HelpText: "Displays info about the specified GitHub user or repo and fetches trending users and repositories.",
			RequiredArguments: []string{
				"username(/repo) **OR** trending repo/user today/week/month (language)",
			},
			Arguments: []CommandArgument{
				{Name: "username", Description: "The GitHub user to fetch info about", ArgType: "string"},
				{Name: "username/repo", Description: "The GitHub repo to fetch info about", ArgType: "string"},
				{Name: "trending repo/user today/week/month (language)", Description: "The trending type to fetch info about", ArgType: "multiple string"},
			},
		}
	}
	botData.Commands["urbandictionary"] = &Command{
		Function: commandUrbanDictionary,
		HelpText: "Displays the definition of a term according to the Urban Dictionary.",
		RequiredArguments: []string{
			"term",
		},
		Arguments: []CommandArgument{
			{Name: "term", Description: "The term to fetch a definition for", ArgType: "string"},
		},
	}
	botData.Commands["balance"] = &Command{
		Function: commandBalance,
		HelpText: "Displays the user's current balance.",
	}
	botData.Commands["daily"] = &Command{
		Function: commandDaily,
		HelpText: "Lets the user receive credits daily.",
	}
	botData.Commands["transfer"] = &Command{
		Function: commandTransfer,
		HelpText: "Transfers credits to another user.",
		RequiredArguments: []string{
			"amount",
			"target",
		},
		Arguments: []CommandArgument{
			{Name: "amount", Description: "The amount of credits to transfer", ArgType: "number"},
			{Name: "target", Description: "The target user to transfer credits to", ArgType: "mention"},
		},
	}

	//Voice commands
	botData.Commands["play"] = &Command{
		Function: commandPlay,
		HelpText: "Plays either the first result from a YouTube search query or the specified stream URL in the user's voice channel.",
		Arguments: []CommandArgument{
			{Name: "search query", Description: "The YouTube search query to use when fetching a video to play", ArgType: "string"},
			{Name: "url", Description: "The YouTube, Spotify, SoundCloud, Bandcamp or direct audio/video URL to play", ArgType: "string"},
		},
	}
	botData.Commands["stop"] = &Command{
		Function: commandStop,
		HelpText: "Stops the audio playback in the user's voice channel.",
	}
	botData.Commands["skip"] = &Command{
		Function: commandSkip,
		HelpText: "Skips to the next queue entry in the user's voice channel.",
	}
	botData.Commands["pause"] = &Command{
		Function: commandPause,
		HelpText: "Pauses the audio playback in the user's voice channel.",
	}
	botData.Commands["resume"] = &Command{
		Function: commandResume,
		HelpText: "Resumes the audio playback in the user's voice channel.",
	}
	botData.Commands["volume"] = &Command{
		Function: commandVolume,
		HelpText: "Sets the volume level for the next audio playback.",
		RequiredArguments: []string{
			"volume",
		},
		Arguments: []CommandArgument{
			{Name: "volume", Description: "The volume level to use", ArgType: "number [0 - 512]"},
		},
	}
	botData.Commands["repeat"] = &Command{
		Function: commandRepeat,
		HelpText: "Switches queue playback between three modes: no repeat, repeat queue, and repeat now playing.",
		Arguments: []CommandArgument{
			{Name: "disable", Description: "Disables repeat mode", ArgType: ""},
			{Name: "queue", Description: "Enables repeat queue mode", ArgType: ""},
			{Name: "now playing", Description: "Enables repeat now playing mode", ArgType: ""},
		},
	}
	botData.Commands["shuffle"] = &Command{
		Function: commandShuffle,
		HelpText: "Toggles queue shuffling during playback.",
	}
	botData.Commands["youtube"] = &Command{
		Function: commandYouTube,
		HelpText: "Allows you to navigate YouTube search results to select what to add to the queue.",
		RequiredArguments: []string{
			"command (value)",
		},
		Arguments: []CommandArgument{
			{Name: "search", Description: "Searches the specified query value", ArgType: "string"},
			{Name: "next", Description: "Navigates forward in a search result's pages", ArgType: ""},
			{Name: "previous", Description: "Navigates backward in a search result's pages", ArgType: ""},
			{Name: "cancel", Description: "Cancels the search result", ArgType: ""},
			{Name: "play", Description: "Plays the chosen search result from the current page", ArgType: "number"},
		},
	}
	botData.Commands["spotify"] = &Command{
		Function: commandSpotify,
		HelpText: "Allows you to search Spotify search results and playlists to select to what to add to the queue.",
		RequiredArguments: []string{
			"command (value)",
		},
		Arguments: []CommandArgument{
			{Name: "search", Description: "Displays track results for the specified search query", ArgType: "string"},
			{Name: "playlist", Description: "Displays track results for the specified playlist", ArgType: "playlist"},
			{Name: "next", Description: "Navigates forward in a playlist's pages", ArgType: ""},
			{Name: "previous", Description: "Navigates backward in a playlist's pages", ArgType: ""},
			{Name: "page/jump", Description: "Jumps to the specified page", ArgType: "number"},
			{Name: "cancel", Description: "Cancels the search/playlist session", ArgType: ""},
			{Name: "play", Description: "Plays the chosen result (single track, 10 popular artist tracks, full album, or list a playlist)", ArgType: "number"},
			{Name: "play all", Description: "Plays every track result", ArgType: ""},
			{Name: "play view", Description: "Plays every track result on the current page", ArgType: ""},
		},
	}
	botData.Commands["queue"] = &Command{
		Function: commandQueue,
		HelpText: "Lists and manages entries in the queue.",
		Arguments: []CommandArgument{
			{Name: "clear", Description: "Clears the queue", ArgType: ""},
			{Name: "remove", Description: "Removes the specified queue entry or entries", ArgType: "number"},
		},
	}
	botData.Commands["nowplaying"] = &Command{
		Function: commandNowPlaying,
		HelpText: "Displays the now playing entry.",
	}
	botData.Commands["lyrics"] = &Command{
		Function: commandLyrics,
		HelpText: "Displays the lyrics for the currently playing track.",
	}

	//All moderation commands with parameters
	botData.Commands["purge"] = &Command{
		Function:            commandPurge,
		HelpText:            "Purges the specified amount of messages from the channel, up to 100 messages at a time.",
		RequiredPermissions: discordgo.PermissionManageMessages,
		RequiredArguments: []string{
			"amount (user1) (user2) (user3)",
		},
		Arguments: []CommandArgument{
			{Name: "message count", Description: "The amount of messages to delete", ArgType: "number"},
			{Name: "user(s)", Description: "The user(s) to delete the messages from within the specified amount of messages", ArgType: "mention"},
		},
	}
	botData.Commands["kick"] = &Command{
		Function:            commandKick,
		HelpText:            "Kicks the specified user(s) from the server.",
		RequiredPermissions: discordgo.PermissionKickMembers,
		RequiredArguments: []string{
			"user1 (user2) (user3) (reason for kick)",
		},
		Arguments: []CommandArgument{
			{Name: "user(s)", Description: "The user(s) to kick", ArgType: "mention"},
			{Name: "reason", Description: "The reason for the kick", ArgType: "string"},
		},
	}
	botData.Commands["ban"] = &Command{
		Function:            commandBan,
		HelpText:            "Bans the specified user(s) from the server.",
		RequiredPermissions: discordgo.PermissionBanMembers,
		RequiredArguments: []string{
			"(days) user1 (user2) (user3) (reason for ban)",
		},
		Arguments: []CommandArgument{
			{Name: "days", Description: "How many days worth of messages to delete from the specified user(s)", ArgType: "number"},
			{Name: "user(s)", Description: "The user(s) to ban", ArgType: "mention"},
			{Name: "reason", Description: "The reason for the ban", ArgType: "string"},
		},
	}
	botData.Commands["hackban"] = &Command{
		IsAdvancedCommand:   true,
		AdvancedFunction:    commandHackBan,
		HelpText:            "Bans the specified user ID(s) from the server.",
		RequiredPermissions: discordgo.PermissionBanMembers,
		RequiredArguments: []string{
			"(-days days) -id user1 (-id user2) (-id user3) (-reason reason for ban)",
		},
		Arguments: []CommandArgument{
			{Name: "days", Description: "How many days worth of messages to delete from the specified user(s)", ArgType: "number"},
			{Name: "id", Description: "The user ID to ban", ArgType: "user ID"},
			{Name: "reason", Description: "The reason for the ban"},
		},
	}

	botData.Commands["server"] = &Command{
		Function:            commandSettingsServer,
		HelpText:            "Changes the specified settings for the server.",
		RequiredPermissions: discordgo.PermissionAdministrator,
		RequiredArguments: []string{
			"setting (value)",
		},
		Arguments: []CommandArgument{
			{Name: "filter", Description: "Manages the swear filter", ArgType: ""},
			{Name: "joinmsg", Description: "Sets the join message for this channel", ArgType: "string"},
			{Name: "leavemsg", Description: "Sets the leave message for this channel", ArgType: "string"},
			{Name: "log", Description: "Manages the logging events", ArgType: ""},
			{Name: "tips", Description: "Enables or disables logging events for this channel", ArgType: "enable/disable"},
			{Name: "autosendnowplaying", Description: "Enables or disables automatically sending now playing embeds without user interaction", ArgType: "enable/disable"},
			{Name: "invitegen", Description: "Manages invite link generation via the API", ArgType: ""},
			{Name: "reset", Description: "Resets the specified setting to the default/empty value", ArgType: "string"},
		},
	}

	botData.Commands["roleme"] = &Command{
		IsAdvancedCommand:   true,
		AdvancedFunction:    commandRoleMe,
		HelpText:            "Allows you to manage the roleme events list. No arguments will list the roleme events.",
		RequiredPermissions: discordgo.PermissionAdministrator,
		Arguments: []CommandArgument{
			{Name: "addrole", Description: "Adds the role to add when this event triggers", ArgType: "role"},
			{Name: "removerole", Description: "Adds the role to remove when this event triggers", ArgType: "role"},
			{Name: "casesensitive", Description: "Whether or not the triggers specified are case-sensitive", ArgType: "boolean"},
			{Name: "channel", Description: "Adds the channel to the channel whitelist for triggering this event", ArgType: "channel"},
			{Name: "trigger", Description: "Adds the message used to trigger this event", ArgType: "message"},
			{Name: "delete", Description: "Deletes the specified roleme entry", ArgType: "number"},
		},
	}
	botData.Commands["bot"] = &Command{
		Function:            commandSettingsBot,
		HelpText:            "Changes the specified settings for the bot within this server.",
		RequiredPermissions: discordgo.PermissionAdministrator,
		RequiredArguments: []string{
			"setting (value)",
		},
		Arguments: []CommandArgument{
			{Name: "prefix", Description: "Sets the bot command prefix", ArgType: "string"},
		},
	}
	botData.Commands["user"] = &Command{
		Function: commandSettingsUser,
		HelpText: "Changes the specified settings for the user.",
		RequiredArguments: []string{
			"setting (value)",
		},
		Arguments: []CommandArgument{
			{Name: "about/aboutme/description/desc/info", Description: "Sets your aboutme or views the aboutme of another user", ArgType: "string/mention"},
			{Name: "timezone", Description: "Sets the timezone to use", ArgType: "timezone"},
			{Name: "social", Description: "Manages your socials", ArgType: ""},
		},
	}

	botData.Commands["starboard"] = &Command{
		Function:            commandStarboard,
		HelpText:            "Manages the guild's starboard.",
		RequiredPermissions: discordgo.PermissionAdministrator,
		RequiredArguments: []string{
			"setting (value)",
		},
		Arguments: []CommandArgument{
			{Name: "enable", Description: "Enables the starboard", ArgType: ""},
			{Name: "disable", Description: "Disables the starboard", ArgType: ""},
			{Name: "channel (set/remove)", Description: "Either returns the current starboard channel or optionally sets it to the current channel or removes the current channel in place", ArgType: ""},
			{Name: "nsfwchannel (set/remove)", Description: "Either returns the current NSFW channel or optionally sets it to the current channel (if marked as NSFW) or removes the current channel in place", ArgType: ""},
			{Name: "emoji (emoji)", Description: "Either returns the current emoji or optionally sets it to the specified emoji", ArgType: "emoji"},
			{Name: "nsfwemoji (emoji)", Description: "Either returns the current NSFW emoji or optionally sets it to the specified emoji", ArgType: "emoji"},
			{Name: "selfstar", Description: "Sets whether or not selfstars are permitted", ArgType: "boolean"},
			{Name: "minimum", Description: "Either returns the current minimum reaction requirement or sets it to the specified amount", ArgType: "number"},
		},
	}

	botData.Commands["feed"] = &Command{
		IsAdvancedCommand:   true,
		AdvancedFunction:    commandFeed,
		HelpText:            "Manages the guild's various RSS and Atom feeds.",
		RequiredPermissions: discordgo.PermissionAdministrator,
		RequiredArguments: []string{
			"-action (value)",
		},
		Arguments: []CommandArgument{
			{Name: "add", Description: "Adds an RSS/Atom feed", ArgType: "url"},
			{Name: "edit", Description: "Edits a given feed entry", ArgType: "id"},
			{Name: "remove", Description: "Removes a given feed entry", ArgType: "id"},
			{Name: "frequency", Description: "How often in seconds to check for a new post, can be used with ``-edit``", ArgType: "seconds"},
			{Name: "setchannel", Description: "When used with ``-edit``, sets the current channel as the feed's post channel", ArgType: ""},
			{Name: "list", Description: "Lists all active feeds", ArgType: ""},
		},
	}

	/*botData.Commands["autoit"] = &Command{
		Function: commandAutoIt,
		HelpText: "Executes AutoIt code using the [AutoGo](https://github.com/JoshuaDoes/autogo) VM.",
		RequiredArguments: []string{
			"script",
		},
		Arguments: []CommandArgument{
			{Name: "script", Description: "The AutoIt code to run", ArgType: "string"},
		},
	}
	botData.Commands["au3"] = &Command{IsAlternateOf: "autoit"}
	botData.Commands["autoit3"] = &Command{IsAlternateOf: "autoit"}
	botData.Commands["autogo"] = &Command{IsAlternateOf: "autoit"}*/

	botData.Commands["pluspro"] = &Command{
		Function: func(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
			patreonTxt, err := os.ReadFile("patreon.txt")
			if err != nil {
				return NewErrorEmbed("Clinet+ Error", "Error reading the latest about Clinet+: %v", err)
			}
			return NewGenericEmbed("Clinet+", string(patreonTxt))
		},
		HelpText: "Explains the Clinet+ subscription offered via Patreon.",
	}
	botData.Commands["donate"] = &Command{IsAlternateOf: "pluspro"}
	botData.Commands["patreon"] = &Command{IsAlternateOf: "pluspro"}
	botData.Commands["paypal"] = &Command{IsAlternateOf: "pluspro"}
	botData.Commands["plus"] = &Command{IsAlternateOf: "pluspro"}
	botData.Commands["pro"] = &Command{IsAlternateOf: "pluspro"}
	botData.Commands["pp"] = &Command{IsAlternateOf: "pluspro"}

	//Alternate commands for pre-established commands
	botData.Commands["?"] = &Command{IsAlternateOf: "help"}
	botData.Commands["commands"] = &Command{IsAlternateOf: "help"}
	botData.Commands["ver"] = &Command{IsAlternateOf: "version"}
	botData.Commands["v"] = &Command{IsAlternateOf: "version"}
	botData.Commands["rolldouble"] = &Command{IsAlternateOf: "doubleroll"}
	botData.Commands["flipcoin"] = &Command{IsAlternateOf: "coinflip"}
	botData.Commands["img"] = &Command{IsAlternateOf: "image"}
	botData.Commands["gh"] = &Command{IsAlternateOf: "github"}
	botData.Commands["yt"] = &Command{IsAlternateOf: "youtube"}
	botData.Commands["sp"] = &Command{IsAlternateOf: "sp"}
	botData.Commands["np"] = &Command{IsAlternateOf: "nowplaying"}
	botData.Commands["q"] = &Command{IsAlternateOf: "queue"}
	botData.Commands["loop"] = &Command{IsAlternateOf: "repeat"}
	botData.Commands["next"] = &Command{IsAlternateOf: "skip"}
	botData.Commands["ud"] = &Command{IsAlternateOf: "urbandictionary"}
	botData.Commands["owo"] = &Command{IsAlternateOf: "hewwo"}
	botData.Commands["uwu"] = &Command{IsAlternateOf: "hewwo"}
	botData.Commands["mc"] = &Command{IsAlternateOf: "minecraft"}
	botData.Commands["guildinfo"] = &Command{IsAlternateOf: "serverinfo"}
	botData.Commands["ss"] = &Command{IsAlternateOf: "screenshot"}
	botData.Commands["credits"] = &Command{IsAlternateOf: "balance"}
	botData.Commands["cash"] = &Command{IsAlternateOf: "balance"}
	botData.Commands["money"] = &Command{IsAlternateOf: "balance"}
	botData.Commands["nightly"] = &Command{IsAlternateOf: "daily"}
	botData.Commands["send"] = &Command{IsAlternateOf: "transfer"}
	botData.Commands["googletranslate"] = &Command{IsAlternateOf: "translate"}
	botData.Commands["gtranslate"] = &Command{IsAlternateOf: "translate"}

	//Google settings and controls
	botData.Commands["gcpauth"] = &Command{Function: commandGCPAuth, HelpText: "Manages Google Cloud Platform authentication.", IsAdministrative: true}

	//Administrative commands for bot owners
	botData.Commands["reload"] = &Command{Function: commandReload, HelpText: "Reloads the bot configuration.", IsAdministrative: true}
	botData.Commands["restart"] = &Command{Function: commandRestart, HelpText: "Restarts the bot in case something goes awry.", IsAdministrative: true}
	botData.Commands["update"] = &Command{Function: commandUpdate, HelpText: "Updates the bot to the latest git repo commit.", IsAdministrative: true}
	botData.Commands["debug"] = &Command{Function: commandDebug, HelpText: "Toggles debug mode.", IsAdministrative: true}
	botData.Commands["sudo"] = &Command{
		Function:         commandSudo,
		HelpText:         "Runs a command as the specified user.",
		IsAdministrative: true,
		RequiredArguments: []string{
			"user", "command (arguments)",
		},
		Arguments: []CommandArgument{
			{Name: "user", Description: "The user to run a command as", ArgType: "mention/ID"},
			{Name: "command", Description: "The command to run as the specified user", ArgType: "string"},
			{Name: "arguments", Description: "Optional additional arguments to pass to the command", ArgType: "N/A"},
		},
	}
	botData.Commands["status"] = &Command{
		Function:         commandStatus,
		HelpText:         "Sets the bot's status message.",
		IsAdministrative: true,
		RequiredArguments: []string{
			"type", "status",
		},
		Arguments: []CommandArgument{
			{Name: "playing", Description: "Sets the playing status", ArgType: "string"},
			{Name: "streaming", Description: "Sets the streaming status and the stream URL", ArgType: "url string"},
			{Name: "listening", Description: "Sets the listening to status", ArgType: "string"},
			{Name: "watching", Description: "Sets the watching status", ArgType: "string"},
		},
	}
}

func callCommand(commandName string, args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	if command, exists := botData.Commands[commandName]; exists {
		if command.IsAlternateOf != "" {
			if commandAlternate, exists := botData.Commands[command.IsAlternateOf]; exists {
				command = commandAlternate
			} else {
				return nil
			}
		}
		if command.IsAdministrative && env.User.ID != botData.BotOwnerID {
			return NewErrorEmbed("Command Error - Not Authorized (NA)", "I'm sorry Dave, I'm afraid I can't do that.")
		}
		if command.RequiredPermissions != 0 {
			if permissionsAllowed, _ := MemberHasPermission(botData.DiscordSession, env.Guild.ID, env.User.ID, env.Channel.ID, command.RequiredPermissions); permissionsAllowed == false {
				return NewErrorEmbed("Command Error - No Permissions (NP)", "Just what do you think you're doing, Dave?")
			}
		}
		if len(args) >= len(command.RequiredArguments) {
			if command.IsAdvancedCommand {
				advancedArgs := make([]CommandArgument, 0)

				//Make sure each legacy argument value is either an argument identifier or an argument value
				for i := 0; i < len(args); i++ {
					if strings.HasPrefix(args[i], "-") {
						if i+1 < len(args) {
							if !strings.HasPrefix(args[i+1], "-") {
								advancedArgs = append(advancedArgs, CommandArgument{Name: strings.TrimSpace(strings.TrimPrefix(args[i], "-")), Value: strings.TrimSpace(args[i+1])})
								i++
								continue
							} else {
								advancedArgs = append(advancedArgs, CommandArgument{Name: strings.TrimSpace(strings.TrimPrefix(args[i], "-")), Value: ""})
								continue
							}
						}
						advancedArgs = append(advancedArgs, CommandArgument{Name: strings.TrimSpace(strings.TrimPrefix(args[i], "-")), Value: ""})
						continue
					} else {
						return getCommandUsage(commandName, "Command Error - Loose Argument Value (LAV)", env)
					}
				}

				return command.AdvancedFunction(advancedArgs, env)
			}
			return command.Function(args, env)
		}
		return getCommandUsage(commandName, "Command Error - Not Enough Parameters (NEP)", env)
	}
	return nil
}

func getCommandUsage(commandName, title string, env *CommandEnvironment) *discordgo.MessageEmbed {
	command := botData.Commands[commandName]
	if command.IsAlternateOf != "" {
		command = botData.Commands[command.IsAlternateOf]
	}

	parameterFields := []*discordgo.MessageEmbedField{}
	parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: "Usage", Value: env.BotPrefix + commandName + " " + strings.Join(command.RequiredArguments, " ")})
	for i := 0; i < len(command.Arguments); i++ {
		if command.IsAdvancedCommand {
			name := "-" + command.Arguments[i].Name
			if command.Arguments[i].ArgType != "" {
				name += " (" + command.Arguments[i].ArgType + ")"
			}
			parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: name, Value: command.Arguments[i].Description, Inline: true})
			continue
		}
		parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: command.Arguments[i].Name + " (" + command.Arguments[i].ArgType + ")", Value: command.Arguments[i].Description, Inline: true})
	}

	usageEmbed := NewEmbed().
		SetTitle(title).
		SetDescription("**" + commandName + "**: " + command.HelpText).
		SetColor(0xFF0000).MessageEmbed
	usageEmbed.Fields = parameterFields

	return usageEmbed
}

func getCustomCommandUsage(command *Command, commandName, title string, env *CommandEnvironment) *discordgo.MessageEmbed {
	parameterFields := []*discordgo.MessageEmbedField{}
	parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: "Usage", Value: env.BotPrefix + commandName + " " + strings.Join(command.RequiredArguments, " ")})
	for i := 0; i < len(command.Arguments); i++ {
		if command.IsAdvancedCommand {
			name := "-" + command.Arguments[i].Name
			if command.Arguments[i].ArgType != "" {
				name += " (" + command.Arguments[i].ArgType + ")"
			}
			parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: name, Value: command.Arguments[i].Description, Inline: true})
			continue
		}
		parameterFields = append(parameterFields, &discordgo.MessageEmbedField{Name: command.Arguments[i].Name + " (" + command.Arguments[i].ArgType + ")", Value: command.Arguments[i].Description, Inline: true})
	}

	usageEmbed := NewEmbed().
		SetTitle(title).
		SetDescription("**" + commandName + "**: " + command.HelpText).
		SetColor(0xFF0000).MessageEmbed
	usageEmbed.Fields = parameterFields

	return usageEmbed
}
