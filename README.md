#pocket2dropbox: <img src="https://raw.githubusercontent.com/msoap/pocket2dropbox/misc/img/pocket_icon.png" height="32" width="32"> ⇢ <img src="https://raw.githubusercontent.com/msoap/pocket2dropbox/misc/img/dropbox_icon.png" height="32" width="32">

Backup [Pocket](http://getpocket.com/) articles to dropbox.

##Install

From source:

    go get -u github.com/msoap/pocket2dropbox
    ln -s $GOPATH/bin/pocket2dropbox ~/bin/pocket2dropbox

##Usage

 * Get Pocket/Dorpbox app_id, keys and tokens - see links below.
 * Create config file.
 * Download wgethtml.pl to PATH dir.
 * Add to cron: `0 * * * * pocket2dropbox`

##Configuration

By config file `~/.config/pocket2dropbox.cfg` ([example](https://raw.githubusercontent.com/github.com/msoap/pocket2dropbox/misc/pocket2dropbox.cfg)):

    {
        "pocket_key": "***",
        "pocket_token": "***",
        "db_client_id": "***",
        "db_client_secret": "***",
        "db_token": "***",
        "get_since_days": 30,
        "favorites": false
    }

or through environment vars:

	# Pocket settings
	POCKET_KEY
	POCKET_TOKEN

	# Dropbox settings
	DB_CLIENTID
	DB_CLIENTSECRET
	DB_TOKEN

options:

    pocket2dropbox [options]
    options
        -favorites        : save favorites articles only
        -get-since-days=N : get articles since this days
        -version
        -help

##Dependencies

[wgethtml.pl](https://gist.github.com/msoap/2567074) - for download html

##Links

 * [Get pocket keys/tokens](https://getpocket.com/developer/docs/authentication)
 * [Get dropbox keys/tokens](https://www.dropbox.com/developers/apps/create)
 * [Pocket API](https://getpocket.com/developer/docs/overview)
 * [Go client library for the Dropbox](https://github.com/stacktic/dropbox)
 * [Save html page with embedding css/js/images in file](https://gist.github.com/msoap/2567074)

##TODO

 * Delete deleted articles in Pocket
 * Rewrite wgethtml.pl on Go
