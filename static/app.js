(function() {
    var commands = 0;
    var ws = new WebSocket("ws://192.168.1.138:8080/ws");
    ws.ononpen = function(evt) {
	console.log("open");
    }
    
    ws.onclose = function(evt) {
	console.log("close");
    }

    ws.onmessage = function(evt) {
	var msg = JSON.parse(evt.data);
	$("#command-"+msg.id + " .output").append(msg.output);
    }
    
    ws.onerror = function(evt) {
	console.log("error: " + evt.data);
    }
    
    $("#prompt").on("keypress", function(e) {
	var key = e.which || e.keyCode;
	if (key !== 13) {
	    return;
	}
	var cmd = $("#prompt").val();
	var html = '<div class="command" id="command-' + commands + '"><p class="cmd">' + cmd + '</p><pre class="output"></pre></div>';
	$("#container").append(html);
	$("#prompt").val("");
	var msg = {
	    id: commands,
	    cmd: cmd
	}
	commands += 1;
	ws.send(JSON.stringify(msg));
    });

    $("#prompt").focus();

})();
