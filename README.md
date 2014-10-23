Stats Gopher
============

[![Circle CI](https://circleci.com/gh/sjltaylor/stats-gopher.png?style=badge)](https://circleci.com/gh/sjltaylor/stats-gopher)


A web stats endpoint for relaying stats to NewRelic Insights view the browser.

Environment variables:
  * `PORT` defaults to `80`
  * `NEW_RELIC_INSIGHTS_KEY`
  * `NEW_RELIC_INSIGHTS_ENDPOINT`
  * `STDOUT_LISTENER=1` prints all events received to STDOUT

`NEW_RELIC_INSIGHTS_*` values come from INSIGHTS ("Manage Data" / "API Keys")
