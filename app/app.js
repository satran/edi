var ConsoleView = Backbone.View.extend({
	className: "console-window",
	template: _.template($('#console-tmpl').html()),
	events: {
		"keydown .command": "keystroke"
	},

	initialize: function(){
		this.render();
	},

	render: function(){
		this.$el.html(this.template());
		return this;
	},

	keystroke: function(ev) {
		switch(ev.which) {
			// Esc
			case 27:
				this.$el.hide();
				break;
		}
	}
});
var WindowView = Backbone.View.extend({
	tagName: "table",
	className: "window",
	template: _.template($('#window-tmpl').html()),

	events: {
		"mousedown .drag": "drag",
		"focus .editor": "highlight",
		"focusout .editor": "unhighlight",
		"keydown .editor": "keystroke"
	},

	initialize: function(args){
		this.parent = args.parent;
		this.render();
	},

	render: function(){
		this.$el.html(this.template({id: this.cid, name: "New File"}));
		return this;
	},

	drag: function(ev){
		this.parent.drag_init(this.el);
	},

	highlight: function(ev) {
		$(this.el).addClass("highlight");
	},

	unhighlight: function(ev) {
		$(this.el).removeClass("highlight");
	},

	keystroke: function(ev) {
		var target = ev.target;
		switch(ev.which){
		//tab
		case 9:
			ev.preventDefault();
			var start = $(target).get(0).selectionStart;
			var end = $(target).get(0).selectionEnd;

			// set textarea value to: text before caret + tab + text after caret
			$(target).val($(target).val().substring(0, start)
				+ "\t"
				+ $(target).val().substring(end));

			// put caret at right position again
			$(target).get(0).selectionStart =
			$(target).get(0).selectionEnd = start + 1;	
			break;
		}
	}
});

var App = Backbone.View.extend({
	el: document,

	events: {
		"mousemove": "move_elem",
		"mouseup": "destroy"
	},

	initialize: function(){
		this.selected = null; // Object of the element to be moved
		this.x_pos = 0;
		this.y_pos = 0; // Stores x & y coordinates of the mouse pointer
		this.x_elem = 0;
		this.y_elem = 0; // Stores top, left values (edge) of the element

		this.command = new ConsoleView();
		$("body").append(this.command.el);
		this.addWindow();
	},

	// Will be called when user starts dragging an element
	drag_init: function(elem) {
		// Store the object of the element which needs to be moved
		this.selected = elem;
		this.x_elem = this.x_pos - this.selected.offsetLeft;
		this.y_elem = this.y_pos - this.selected.offsetTop;
	},

	// Will be called when user dragging an element
	move_elem: function (e) {
		this.x_pos = document.all ? window.event.clientX : e.pageX;
		this.y_pos = document.all ? window.event.clientY : e.pageY;
		if (this.selected !== null) {
			this.selected.style.left = (this.x_pos - this.x_elem) + 'px';
			this.selected.style.top = (this.y_pos - this.y_elem) + 'px';
		}
	},

	// Destroy the object when we are done
	destroy: function () {
		this.selected = null;
	},

	addWindow: function() {
		var w = new WindowView({parent: this});
		$("body").append(w.el);
	}
});


