define(["backbone"], function(Backbone) {
var Command = Backbone.Model.extend({		
	generate: function(options) {
		var attr = this.attributes;
		if (options.Command !== undefined)
			attr.Command = options.Command;
		if (options.Args !== undefined)
			attr.Args = options.Args;
		return attr;
	}
});

var CommandList = Backbone.Collection.extend({
	model: Command,
	
	get: function(id) {
		return this.findWhere({Id: id});
	},
	
	getOrCreate: function(id) {
		var model = this.get(id);
		if (model === undefined) {
			model = this.create({Id: id});	
		}
		return model;
	},
	
	// Creates a new Command.
	new: function(attributes) {
		var model = new Command(attributes);
		model.set({"Id": model.cid});
		this.add(model);
		return model;
	}
});
	
return {
	Command: Command,
	Commands: CommandList
};
});