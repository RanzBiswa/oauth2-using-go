# oauth2-using-go

This code helps you to authenticate users via OAuth2 mechanism using gmail only

## Pre Requisites
Whenever we need to play around with OAuth2 gmail mechanim , Below are few list of set up that needs to be done:

-  Go to Google Cloud Platform
-  Create new project or use an existing one
-  Go to Credentials
-  Click “Create credentials”
-  Choose “OAuth client ID”
-  Add authorized redirect URL, in our case it will be localhost:8181/userprofile
-  Get client id and client secret
-  Download the JSON to your local.

Set up is done.


## Overview of code structure?
There are 3 files
- login.html  {login UI where Login with Google option is there which will go through OAuth2 mechanism}
- main.go     {base code}
- userprofile.html { redirect page after passing authentication}


