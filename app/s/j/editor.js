// Defines all the listeners and functions for the editor view.
define(["lib/codemirror"], function(CodeMirror) {
return function(editor, status){
	var self = CodeMirror(editor, {
		lineNumbers: true,
		lineWrapping: true,
		electricChars: true,
	    	indentUnit: 8,
	    	tabSize: 8,
	    	indentWithTabs: true,
		smartIndent: true
	});

	require(["addon/selection/active-line"], function() {
		self.setOption("styleActiveLine", true);
	});
	
	self.setTheme = function(theme){
		var url = "s/c/theme/" + theme + ".css";
		var link = document.createElement("link");
		link.type = "text/css";
		link.rel = "stylesheet";
		link.href = url;
		document.getElementsByTagName("head")[0].appendChild(link);
		self.setOption("theme", theme);
	}

	self.setMode = function(mode){
		require(["mode/" + mode + "/" + mode], function() {
			self.setOption("mode", mode);
		});
	};

	self.setKeymap = function(keymap){
		require(["keymap/" + keymap], function() {
			self.setOption("keyMap", keymap);
		});
	};

	self.getCurrentLine = function() {
		var line = self.getCursor().line;
		return self.getLine(line);
	};
	
	self.append = function(text){
		self.replaceRange(
			text,
			CodeMirror.Pos(self.lastLine())
		);
		var lastline = self.lastLine();
		self.setCursor(lastline);
	}

	// Stores file name and CodeMirror.Doc
	var buffers = {};

	self.newBuffer = function(filename, file, mode) {
		var doc;
		if (mode !== undefined && mode !== ""){
			require(["mode/" + mode + "/" + mode], function() {
				doc = CodeMirror.Doc(file, mode, 0);
				doc.filename = filename;
				buffers[filename] = doc;
				self.swapDoc(doc);
			});
		} else {
			doc = CodeMirror.Doc(file, mode, 0);
			doc.filename = filename;
			buffers[filename] = doc;
			self.swapDoc(doc);
		}
		self.focus();
		status.html(filename);
	}

	self.openBuffer = function(filename) {
		var doc = buffers[filename];
		self.swapDoc(doc);
		self.focus();
		status.html(filename);
	}

	self.getCurrentBufferName = function(){
		return self.getDoc().filename;
	}

	self.getText = function(filename) {
		if (filename === undefined){
			filename = self.getCurrentBufferName();
		}
		return self.getDoc().getValue();
	}
	return self;
}
});
