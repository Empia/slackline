# Share channels between Slack accounts!

We are big fans of [Slack][slack] and have been using it for a while
now. Our friends at [Vizzuality][vizzuality] started using it recently
and we missed not being able to speak with them, so I just created a
quick service to link two channels from two Slack organizations.

## How can I set it up?

Dead simple, follow this steps from each organization's account.

 1. Create a channel you want to share with another organization.
 2. Create an Incoming WebHook itntegration and select the channel you created.
 3. Copy the Incoming WebHook token (you can find it in the left column
    from the integration page).
 4. Create a URL with the following format: ```http://slackline.herokuapp.com/bridge/?token=[TOKEN]&domain=[YOUR_SLACK_DOMAIN]``` send it to the person setting up the other organization.
 5. The person setting up the other organization will send you a similar
    URL with their domain and token, create an Outgoing WebHook with
    that URL and the channel you created in step 1.

Once you have done this in both organizations, you will have a chat-room
shared by both organizations.

Here you have an example of a Outgoing WebHook URL:

```
http://slackline.herokuapp.com/bridge/?token=bcaa5867b1d42142b74eDVA4&domain=avengers.slack.com
```

## How does it work?

We are just bridging hooks, we don't store any messages going through
the bridge.

## ToDo

We don't support custom avatars right now.

If somebody from Slack is reading this, it would be great to be able to
set the ```icon_url``` for the message to be something like
```http://[DOMAIN]/avatar/[USER_ID]``` to set the USER_ID based on the
information from the Outgoing WebHook with no extra requests or state.

[slack]: http://slack.com
[vizzuality]: http://vizzuality.com
