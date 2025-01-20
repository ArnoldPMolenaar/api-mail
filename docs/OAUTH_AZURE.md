# Oauth2 client in Azure portal
This document explains how to register an account inside the Azure portal and create an OAuth2 client to use the Outlook API.

## Step 1: Create an application
1. Go to the [Azure portal](https://portal.azure.com/).
2. Click or search for `Microsoft Entra ID`.
3. Click on `App registrations`.
4. Click on `New registration`.
5. Enter a `Name` for the application.
6. Select `Accounts in this organizational directory only`.
7. Enter the following URI in the `Redirect URI` field:
    - `{DOMAIN_NAME}/v1/oauth2/azures/callback`
8. Select platforms:
    - `Web`
9. Click on `Register`.

> Write down the `Application (client) ID` and `Directory (tenant) ID` values.

## Step 2: Create a client secret
1. Under `Manage` click on `Certificates & secrets`.
2. Click on `New client secret`.
3. Enter a `Description` for the client secret.
4. Select an `Expiration` date.
5. Click on `Add`.

> Write down the `Value` of the client secret.

## Step 3: Configure the API permissions
1. Under `Manage` click on `API permissions`.
2. Click on `Add a permission`.
3. Click on `Microsoft Graph`.
4. Click on `Delegated permissions`.
5. Select the following permissions:
    - `offline_access`
    - `User.Read`
    - `Mail.Send`
    - `openid`
6. Click on `Add permissions`.

> That's it! You have successfully created an OAuth2 client in the Azure portal.
> You can now add this account to the application by using the POST.

## Step 4: Configure the environment variables
1. Create a `.env` file in the root of the project.
2. Make sure that `DOMAIN_NAME` is set to its respective domain.

## Step 5: Run the application
1. Run the application.
2. Add a azure by using the post.
3. In the response go to the `AuthCodeURL` field.
4. Follow the consent screen steps.
5. If the results end's in a json response with a token, the process was successful.
