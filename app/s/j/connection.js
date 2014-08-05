define([], function() {
return function(options){
	var self = {
		command: options.command
	}
	
	// Flag denoting if web socket connection is active
	self.ready = false;
	
	// We shall use the location of the window if location is not defined.
	if (self.location === undefined) {
		self.location = "ws://" + window.location.hostname + ":" + window.location.port + "/ws";
	}
	
	self.ws = new WebSocket(self.location);
	
	self.ws.onopen = function () {
		self.ready = true;	
	}

	self.ws.onmessage = function(evt){
		var response = JSON.parse(evt.data);
		self.command.evaluate(response);
	}
	
	self.ws.onerror = function(evt){
		console.log("error", evt);	
	}
	
	self.send = function(data) {
		self.ws.send(data);	
	}
	
	return self;
}
});
