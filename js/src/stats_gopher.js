function StatsGopher (options) {
  options = options || {}

  if (!(('jQuery' in options) && (typeof options.jQuery.ajax === 'function'))) {
    throw new Error("no 'jQuery.ajax' option specified")
  }

  if (typeof options.endpoint !== 'string') {
    throw new Error("no 'endpoint' option specified")
  }

  this.options = options
  this.buffer = []
  this.sid = StatsGopher.sid();
}

StatsGopher.prototype = {
  send: function (datum) {
    datum.sendTime = new Date().valueOf()
    datum.sid = this.sid
    this.startTimeout()
    this.buffer.push(datum)
  },
  startTimeout: function () {
    if ('timeout' in this) return
    this.timeout = setTimeout(this.onTimeout.bind(this), 100)
  },
  flush: function () {
    var buffer = this.buffer
    this.buffer = []
    return buffer
  },
  onTimeout: function () {
    delete this.timeout
    var data = this.flush()
    var options = this.options
    var jQuery = options.jQuery

    return jQuery.ajax({
      type: "POST",
      url: options.endpoint,
      data: JSON.stringify(data),
      // because CORS doesn't allow application/json
      dataType: 'text',
      cache: false
    });
  }
}

if ('window' in this) {
  exports = window
}

StatsGopher.sid = function () {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
    return v.toString(16);
  });
}

exports.StatsGopher = StatsGopher
