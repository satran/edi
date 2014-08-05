require.config({
	// This reloads required scripts.
	urlArgs: "bust=" + (new Date()).getTime()
});

require(["jquery", "command", "editor", "connection"], function($, Command, Editor, Connection) {
	var debug = true;
	var app = {};

	app.el = {
		status: $("#status"),
		// Editor's element is only used for CM, 
		// that's why we need the actual DOM element.
		editor: $("#editor")[0]
	}

	app.editor = Editor(app.el.editor, app.el.status);
	app.editor.setKeymap("vim");
	app.editor.setTheme("solarized light");
	
	app.command = new Command({editor: app.editor});
	$("#container").prepend(app.command.render().el);
	app.command.focus();
	app.command.on("send", function(cmd) {
		app.conn.send(cmd);	
	});

	app.conn = Connection({
		command: app.command	
	});

	if (debug !== undefined && debug){
		window.app = app;
	}
	
	require(["jqxcore"], function(){
		require(["jqxbuttons", "jqxsplitter"], function(){
			$("#container").jqxSplitter({width: '100%', height: '100%'});
		});
	});

	document.addEventListener("keydown", function(event){
		// Shift focus to command on ctrl+space
		// and Editor on ctrl+shift+space
		if (event.which === 32){
			if (event.ctrlKey){
				if (document.activeElement === app.command.$el.find(".bar")[0]){
					app.editor.focus();	
					return;
				}
				app.command.focus();
				event.preventDefault();
			}
		}
	});
});
