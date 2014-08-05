define(["jquery", "underscore", "backbone"], function($, _, Backbone) {
return Backbone.View.extend({
        tagName: "div",
        className: "stdout",
        template: _.template($("#output-tmpl").html()),
        
        events: {
		'click .prompt': "toggle",
		'click .cancel': "cancel",
		'contextmenu pre': "context",
		'contextmenu .command': "context"
        },
        
        initialize: function(options) {
		this.parent = options.parent;
		this.toggled = false;
		this.initActive = (options.active === undefined) ? true: options.active;
        },

	$cancel: function(){
		return this.$el.find(".cancel");
	},
        
	$out: function(){
		return this.$el.find(".out");
	},
        
	$prompt: function(){
		return this.$el.find(".prompt");
	},
        
        render: function(command, output) {
		this.$el.html(this.template({
			'command': this.model.raw,
			'active': this.initActive
		}));
		if (output !== undefined) this.append(output);
		return this;
        },
        
        append: function(lines) {
		if (!(lines instanceof Array)) {
			lines = lines.split("\n");
		}
		for (var i = 0; i < lines.length; i++) {
			var node = document.createElement("pre");
			if (lines[i].length == 0) continue;
			var pre = document.createTextNode(lines[i]);
			node.appendChild(pre);
			this.$out().append(node);
		}
        },
        
        toggle: function() {
		if (this.toggled) {
			this.$prompt().html("&#187;");
			this.$out().show();
			this.toggled = false;
		} else {
			this.$prompt().html("+");
			this.$out().hide();
			this.toggled = true;
		}
        },
	
	cancel: function(){
		// If the command is not active lets remove it.
		if (!this.$prompt().hasClass("active")) {
			this.remove();	
			return;
		}
		this.parent.send({
			Command: "Cancel", 
			Args: [], 
			Id: this.model.get("Id")
		});
	},
      
	context: function(evt){
		var filenames = evt.currentTarget.innerHTML.trim("\n").split(" ");
		this.parent.send({
			Command: "Context",
			Args: filenames,
			Pwd: this.model.get("Pwd")
		});
		evt.preventDefault();
	},
	
	done: function(){
		this.$prompt().removeClass("active");
	}
});
});
