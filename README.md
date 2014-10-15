Stats Gopher
============

A web stats endpoint for relaying stats to NewRelic Insights view the browser.

Environment variables:
  * `PORT` defaults to `80`
  * `NEWRELIC_INSIGHTS_KEY`
  * `NEWRELIC_INSIGHTS_ENDPOINT`
  * `STDOUT_LISTENER=1` prints all events received to STDOUT

`NEWRELIC_*` values come from INSIGHTS ("Manage Data" / "API Keys")
