var chai   = require('chai');
var expect = chai.expect;
var sinon = require('sinon');
var StatsGopher = require('stats_gopher').StatsGopher;

describe('StatsGopher', function() {
  function catchError(fn) {
    try {
      fn()
    } catch (e) {
      return e.toString()
    }
  }

  describe('StatsGopher(options)', function() {
    it('throws if the options object does not have a jQuery.ajax function', function () {
      expect(catchError(function () {
        new StatsGopher()
      })).to.equal("Error: no 'jQuery.ajax' option specified")
    });
    it('throws if the options object does not have an endpoint', function () {
      expect(catchError(function () {
        new StatsGopher({
          jQuery: { ajax: function () {} }
        })
      })).to.equal("Error: no 'endpoint' option specified")
    });
    it('assigns the options to the statsGopher instance', function () {
      var options = {
        jQuery: { ajax: function () {} },
        endpoint: "meow"
      }
      var statsGopher = new StatsGopher(options)
      expect(statsGopher.options).to.eq(options)
    });
    it('assigns an empty buffer', function () {
      var options = {
        jQuery: { ajax: function () {} },
        endpoint: "meow"
      }
      var statsGopher = new StatsGopher(options)
      expect(statsGopher.buffer).to.have.length.of(0);
    });
  });
  describe('instance methods', function() {
    var statsGopher, options;

    beforeEach(function () {
      options = {
        jQuery: {
          ajax: sinon.spy()
        },
        endpoint: 'test-endpoint'
      }
      statsGopher = new StatsGopher(options)
    });
    describe('send(datum)', function() {
      it('stamps the datum with a sendTime of now', function() {
        var datum = {
          eventType: 'test-event'
        };
        var t0 = new Date().valueOf();
        statsGopher.send(datum);
        var t1 = new Date().valueOf();
        expect(datum.sendTime).to.be.at.least(t0)
        expect(datum.sendTime).to.be.at.most(t1)
      });
      it('calls starts the timeout', function () {
        statsGopher.startTimeout = sinon.spy()
        statsGopher.send({
          eventType: 'test-event'
        });
        expect(statsGopher.startTimeout.called).to.equal(true)
      });
      it('puts the datum in the buffer', function () {
        var datum = {
          eventType: 'test-event'
        };

        statsGopher.send(datum);
        debugger
        expect(statsGopher.buffer.indexOf(datum)).to.equal(0)
      });
    });
    describe('startTimeout', function() {
      describe('when there is no current timeout', function() {
        it('sets a timeout', function () {
          expect(statsGopher.timeout).to.equal(undefined)
          statsGopher.startTimeout()
          expect(statsGopher.timeout).not.to.equal(undefined)
        });
      });
      describe('when there is a current timeout', function() {
        it('does not change it', function () {
          statsGopher.timeout = 34567
          statsGopher.startTimeout()
          expect(statsGopher.timeout).to.equal(34567)
        });
      });
    });
    describe('flush()', function() {
      it('reinitializes the buffer', function () {
        statsGopher.send({eventType: 'test-event'});
        var originalBuffer = statsGopher.buffer
        expect(statsGopher.buffer.length).not.to.equal(0)
        statsGopher.flush()
        expect(statsGopher.buffer.length).to.equal(0)
        expect(statsGopher.buffer).not.to.equal(originalBuffer)
      });
      it('returns the current buffer', function () {
        var e = { eventType: 'test-event' }
        statsGopher.send(e);
        var originalBuffer = statsGopher.buffer
        var flushedBuffer = statsGopher.flush()
        expect(flushedBuffer).to.equal(originalBuffer)
        expect(flushedBuffer[0]).to.equal(e)
      });
    });
    describe('onTimeout', function() {
      var buffer = [{eventType:"test"}], ajaxOptions;

      beforeEach(function () {
        statsGopher.flush = sinon.spy(function () {
          return buffer
        })
        statsGopher.options.jQuery.ajax = sinon.spy(function (options) {
          ajaxOptions = options
        })
      })
      it('deletes the timeout', function () {
        statsGopher.startTimeout()
        expect(statsGopher.timeout).not.to.equal(undefined)
        statsGopher.onTimeout()
        expect(statsGopher.timeout).to.equal(undefined)
      });
      it('flushes the buffer', function () {
        statsGopher.onTimeout()
        expect(statsGopher.flush.called).to.equal(true)
      });
      it('calls jQuery.ajax with the correct options', function () {
        statsGopher.onTimeout();
        expect(statsGopher.options.endpoint.length).not.to.equal(0)
        expect(ajaxOptions.url).to.equal(statsGopher.options.endpoint)
        expect(ajaxOptions.type).to.equal('POST')
        expect(ajaxOptions.dataType).to.equal('json')
        expect(ajaxOptions.data).to.equal(JSON.stringify(buffer))
        expect(ajaxOptions.cache).to.equal(false)
      });
    });
  });
});
