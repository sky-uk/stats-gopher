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
}

StatsGopher.prototype = {
  send: function (datum) {
    datum.sendTime = new Date().valueOf()
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

exports.StatsGopher = StatsGopher