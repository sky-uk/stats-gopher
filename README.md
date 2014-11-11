Stats Gopher
============

[![Circle CI](https://circleci.com/gh/sjltaylor/stats-gopher.png?style=badge)](https://circleci.com/gh/sjltaylor/stats-gopher)


A web stats endpoint for relaying stats to NewRelic Insights from the browser.

Environment variables:
  * `PORT` defaults to `80`
  * `NEW_RELIC_INSIGHTS_KEY`
  * `NEW_RELIC_INSIGHTS_ENDPOINT`
  * `STDOUT_LISTENER=1` prints all events received to STDOUT

`NEW_RELIC_INSIGHTS_*` values come from Insights (See "Manage Data" / "API Keys" within Insights)

## Events

`http://wherever.its.hosted.net/gopher/`

`POST` Any data here to have it relayed to insights

Stats Gopher does no processing of the data; events are relayed unchanged

## Presence

`http://wherever.its.hosted.net/presence/`

This allows a browser to send `heartbeat` and `user-activity` notifications.
The stats gopher is currently (statically) configured to monitor for only
these two types of presence. Other notifications would be ignored.

Expected usage scenario:

Browsers sends a regular (~10s) heartbeat notification:

```
{
  key: "<session-id>@my-site",
  code: "heartbeat"
}
```

and sends useractivity notifications whenever a window mouse event is triggered:

```
{
  key: "<session-id>@my-site",
  code: "user-activity"
}
```

The stats gopher is configured with a heartbeat timeout of `45s` and a user
activity timeout of `45min`. If either timeout occurs the session is condisered
dead and an event is sent to new relic with the session duration

Note the key. The key could be used to monitor a session, sitewide as above,
or a single page: `<session-id>@my-site/my-page` or even all users on a page:
`my-site/my-page`
