(function() {
    const editor = document.getElementById("editor");
    const tabsize = 4;
    const keymap = {
        "<": { value: "<>", pos: 1 },
        "(": { value: "()", pos: 1 },
        "{": { value: "{}", pos: 1 },
        "[": { value: "[]", pos: 1 },
        "'": { value: "''", pos: 1 },
        '"': { value: '""', pos: 1 },
        "“": { value: "“”", pos: 1 },
        "`": { value: "``", pos: 1 },
        "‘": { value: "‘’", pos: 1 }
    };
    const snipmap = {
        // These make no sense but I'll add them for completeness
        "1#": "# ",
        "2#": "## ",

        // These make sense
        "3#": "### ",
        "4#": "#### ",
        "5#": "##### ",
        "6#": "###### ",

        // Might be a good idea to add a snippet for a tables sometime.
        "@l": '((link ""))'
    };
	  const filename = document.location.pathname;

    editor.moveCursor = function(i) {
        const pos = editor.selectionStart;
        editor.setSelectionRange(pos+i, pos+i);
    };

    editor.insertText = function(text) {
        editor.focus();
        const pos = editor.selectionStart;
        editor.setRangeText(text, pos, pos, "end");
    };

    function getWord(text, caretPos) {
        let preText = text.substring(0, caretPos);
        let split = preText.split(/\s/);
        return split[split.length - 1].trim();
    }

    function looksLikeBullet(text, caretPos) {
        let bulletRegex = /^([ \t]*[\*\-\+]\s*).*/gim;
        let line = text
            .substring(0, caretPos)
            .split(/\r?\n|\r/)
            .pop();
        let numberedListRegex = /^([ \t]*\d+\.\s*).*/gim;
        if (bulletRegex.test(line)) {
            return {
                bullet: line.replace(bulletRegex, "$1")
            };
        } else if (numberedListRegex.test(line)) {
            return {
                bullet: line
                    .replace(numberedListRegex, "$1")
                    .replace(/\d+/, (number) => +number + 1)
            };
        }
        return false;
    }

    function resize() {
        editor.style.height = '5px';
        editor.style.height = (editor.scrollHeight) + 'px';
    }

	  let cancel;
    editor.addEventListener("keydown", function (event) {
        if (keymap[event.key]) {
            event.preventDefault();
            const pos = editor.selectionStart;
            editor.setRangeText(keymap[event.key].value);
            editor.moveCursor(keymap[event.key].pos);
        }

        if (event.key === "Tab") {
            const word = getWord(editor.value, editor.selectionStart);
            if (word && snipmap[word]) {
                event.preventDefault();
                const pos = editor.selectionStart;
                editor.value =
                    editor.value.slice(0, pos - word.length) +
                    snipmap[word] +
                    editor.value.slice(editor.selectionEnd);

                editor.selectionStart = editor.selectionEnd =
                    pos + (snipmap[word].length - 1);
            } else {
                event.preventDefault();
                editor.insertText(" ".repeat(tabsize));
            }
        }

        if (event.key === "Enter") {
            let bullet = looksLikeBullet(editor.value, editor.selectionStart);
            if (bullet) {
                event.preventDefault();
                let addition = editor.value.substring(editor.selectionStart);
                editor.value = editor.value.substring(0, editor.selectionStart);
                editor.value += "\n" + bullet.bullet + addition;
            }
        }
        //auto_grow(editor);

	      if (cancel) clearTimeout(cancel);
	      cancel = setTimeout(() => {
		        localStorage.setItem(filename, event.target.value);
		    }, 1000);
    });

    editor.addEventListener('input', resize);

    function dropHandler(ev) {
        const fileInput = document.getElementById("file-input");
        console.log('File(s) dropped');

        // Prevent default behavior (Prevent file from being opened)
        ev.preventDefault();
        let names = "";
        let dT = new DataTransfer();
        if (ev.dataTransfer.items) {
	          // Use DataTransferItemList interface to access the file(s)
	          for (var i = 0; i < ev.dataTransfer.items.length; i++) {
	              // If dropped items aren't files, reject them
	              if (ev.dataTransfer.items[i].kind === 'file') {
		                var file = ev.dataTransfer.items[i].getAsFile();
		                console.log('... file[' + i + '].name = ' + file.name);
		                dT.items.add(file);
		                names += file.name + ",";
	              }
	          }
        } else {
	          // Use DataTransfer interface to access the file(s)
	          for (var i = 0; i < ev.dataTransfer.files.length; i++) {
	              console.log('... file[' + i + '].name = ' + ev.dataTransfer.files[i].name);
	              if (ev.dataTransfer.files[i].kind === 'file') {
		                var file = ev.dataTransfer.files[i].getAsFile();
		                console.log('... file[' + i + '].name = ' + file.name);
		                dT.items.add(file);
		                names += file.name + ",";
	              }
	          }
        }
        fileInput.files = dT.files;
        const nameInput = document.querySelector('input[name="name"]');
        nameInput.value = names.slice(0, -1);
        document.querySelector(".editor").hidden = true;
    }

    function dragOverHandler(ev) {
        // Prevent default behavior (Prevent file from being opened)
        ev.preventDefault();
    }

    const savebtn = document.getElementById("save-btn");
	  const form = document.getElementById("new-form");
	  savebtn.addEventListener('click', event => {
	      localStorage.removeItem(filename);
	      form.submit();
	  });

	  // Load draft from localstorage
	  let value = localStorage.getItem(filename);
	  if (value && value.length > 0) {
	      if (confirm('Load unsaved draft?')) {
	          editor.value = value;
	      }
	  }

    // on load for the first time resize the editor to full page
    resize();
})();
