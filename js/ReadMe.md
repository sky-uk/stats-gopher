# Stats Gopher js client

## Installing

`npm run dist` builds and minified the library to `./dist`

Serve `dist/stats_gopher[.min].js` from within your app


## Usage

### Initialization

```
var statsGopher = new StatsGopher({
  jQuery: $,
  endpoint: "http://example.net"
})

```

where `jQuery` is expected to be a jQuery-like interface with a `.ajax` method

### Sending Data

```
  statsGopher.send({
    eventType: 'your-event-type',
    ...
  })
```

* add arbitary data to the event
* the `eventType` field is mandatory
* the `send` method stamps the event with a `sendTime`
