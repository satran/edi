define(["jquery", "underscore", "backbone"], function($, _, Backbone) {
return Backbone.View.extend({
	id: "input", 
	
        template: _.template($("#prompt-tmpl").html()),
	
	events: {
		"keydown": "shortcuts"	
	},
	
	initialize: function(options) {
		this.parent = options.parent;	
	},

	$bar: function(){
		return this.$el.find(".bar");	
	},
	
        render: function() {
		this.$el.html(this.template());
		return this;
        },
	
	shortcuts: function(event){
		switch (event.which){	
		// Enter to execute the command
		case 13:
			this.execute(event);
			break;
		// Ctrl + L to clear screen
		case 76:
			if (!event.ctrlKey){
				return;
			}
			this.clearConsole(event);
			break;
		}
	},
	execute: function(event) {
		this.parent.execute(this.$bar().val());
		this.$bar().val("");
		event.preventDefault();
	},
	clearConsole: function(event) {
		this.parent.clear();
		event.preventDefault();
	}
});
});