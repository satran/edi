define(["jquery", "underscore", "backbone", "models", "output", "prompt"], function($, _, Backbone, models, Output, Prompt) {
return Backbone.View.extend({
	className: "cols",
	
	id: "command_col",
	
        template: _.template($("#command-tmpl").html()),
	
	// Commands that need not to be logged 
	exempted: ["W", "Context"],
	
	initialize: function(options){
		this.editor = options.editor;
		this.commands = new models.Commands();
		this.prompt = new Prompt({parent: this});
		this.prompt.on("execute", this.execute, this);
		this.prompt.on("clear", this.clear, this);
	},
	
        render: function() {
		this.$el.html(this.template());
		this.$el.append(this.prompt.render().el);
		return this;
        },
	
	$output: function(){
		return this.$el.find("#output");	
	},
	
	initView: function(model) {
		var view = new Output({model: model, parent: this});
		this.$output().append(view.render().el);
		return view;
	},

	clear: function() {
		this.$output().empty();
		this.commands.reset();
	},

	execute: function(text){
		var args = text.split(" ");
		var command = args.shift();

		switch(command){
			case "W":
				this.saveBuffer(command, args);
				break;
			default:
				// Send to server.
				this.send({Command: text, raw: text});
		};

	},

	evaluate: function(response){
		if (response.Command === undefined) {
			console.log("Server passed wrong format.", response);
			return;
		}

		switch (response.Command) {
			case "console":
				this.add(response);
				break;
			case "newFile":
				this.editor.newBuffer(response.Args[0],
					     response.Args[1], response.Args[2]);
				break;
			case "openFile":
				this.editor.openBuffer(response.Args[0]);
				break;
			case "openDir":
				this.openDir(response);
				break;
			case "done":
				this.done(response);
				break;
			default:
				this.log("Can't process command.");
		}
	},

	// Send to the server.
	send: function(options) {
		var command;
		if (options.Id !== undefined) {
			command = this.commands.get(options.Id);	
			if (command === undefined){
				console.log("Cant find existing command.");	
				return;
			}
			this.trigger("send", JSON.stringify(command.generate(options)));
		} else {
			command = this.commands.new({
				Command: options.Command,
				Args: options.Args,
				Pwd: options.Pwd
			});
			this.trigger("send", JSON.stringify(command));
		}
		
		// If the view was already created lets continue
		if (command.view !== undefined){
			return;
		}
		
		// raw is the command typed in the prompt.
		if (options.raw !== undefined){
			command.raw = options.raw;
		}

		// Commands that are exempted should not have a view
		if (this.exempted.indexOf(command.get("Command")) == -1){
			command.view = this.initView(command);
		}
		
	},

	saveBuffer: function(command, args) {
		if (args === undefined || args.length == 0){
			args = [this.editor.getCurrentBufferName()];
		}
		for (var i=0; i< args.length; i++){
			this.send({
				Command: command,
				Args: [args[i], this.editor.getText(args[i])]
			});
		}
	},

	add: function(cmd){
		var model = this.commands.get(cmd.Id);
		if (model === undefined) {
			this.commands.new(cmd);	
		} else {
			// Update necessary attributes	
			if (cmd.Pwd !== undefined) {
				model.set({Pwd: cmd.Pwd});	
			}
		}
		if (cmd.raw !== undefined) {
			model.raw = cmd.raw;	
		}
		if (model.view === undefined){
			model.view = this.initView(model);
		}
		model.view.append(cmd.Args);
	},

	openDir: function(cmd){
		cmd.raw = cmd.Args.shift();
		cmd.Args = cmd.Args.shift();
		this.add(cmd);
	},

	done: function(cmd) {
		var model = this.commands.findWhere({Id: cmd.Id});
		if (model.view === undefined) {
			return;	
		}
		model.view.done();
	},

	focus: function(){
		this.prompt.$bar().focus();
	}
});
});
