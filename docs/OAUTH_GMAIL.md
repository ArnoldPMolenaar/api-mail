# Oauth2 client in Google Cloud Console
This document explains how to register an account inside the Google Cloud Console and create an OAuth2 client to use the Gmail API.

## Step 1: Create a project
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Click on the project selector and then click on `New Project`.
3. Enter a project name and click on `Create`.

## Step 2: Create credentials
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Click on the project selector and select the project you created in the previous step.
3. Click on the `APIs & Services` menu and then click on `Credentials`.
4. Click on `Create credentials` and then click on `OAuth client ID`.
5. Select `Web application` as the application type.
6. Enter a `Name` for the OAuth client.
7. Enter the following URI in the `Authorized redirect URIs` field:
    - `{DOMAIN_NAME}/v1/oauth2/gmails/callback`
8. Click on `Create`.
9. Write down the `Client ID` and `Client Secret` values.

## Step 3: Configure the consent screen
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Click on the project selector and select the project you created in the previous step.
3. Click on the `APIs & Services` menu and then click on `OAuth consent screen`.
4. Enter the `App name`.
5. Enter the `User support email`.
6. Enter the `Developer contact information`.
7. Click on `Save`.
8. Click on `Add or remove scopes`.
9. Select the following scopes:
    - `openid`
    - `https://www.googleapis.com/auth/gmail.send`
    - `profile`
    - `email`
10. Click on `Update`.
11. Click on `Save`.
12. (optional) Add test users (for example your own gmail account).
13. Click on `Save`.
14. Click on `Back to Dashboard`.

## Step 4: Enable the Gmail API
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Click on the project selector and select the project you created in the previous step.
3. Click on the `APIs & Services` menu and then click on `Library`.
4. Search for `Gmail API` and click on it.
5. Click on `Enable`.

## Step 5: Configure the environment variables
1. Create a `.env` file in the root of the project.
2. Make sure that `DOMAIN_NAME` is set to its respective domain.

## Step 6: Run the application
1. Run the application.
2. Add a gmail by using the post.
3. In the response go to the `AuthCodeURL` field.
4. Follow the consent screen steps.
5. If the results end's in a json response with a token, the process was successful.
