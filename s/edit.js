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
        editor.style.height = (editor.scrollHeight) + 'px';
    }

    let cancel;
    editor.addEventListener("keydown", function (event) {
	// disabling the global shortcuts to be called
	event.stopPropagation();

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
                editor.insertText("\n" + bullet.bullet);
            }
        }

	      if (cancel) clearTimeout(cancel);
	      cancel = setTimeout(() => {
		        localStorage.setItem(filename, event.target.value);
		    }, 1000);
    });

    editor.addEventListener('input', resize);

    const fileEl = document.querySelector('.file');
    ['dragenter', 'dragover'].forEach(eventName => {
        fileEl.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        fileEl.addEventListener(eventName, unhighlight, false);
    });

    function highlight(e) {
        fileEl.classList.add('highlight');
    }

    function unhighlight(e) {
        fileEl.classList.remove('highlight');
    }

    fileEl.addEventListener("dragstart", ev => {
        ev.dataTransfer.setData("text",  "data");
        ev.dataTransfer.effectAllowed = "move";
    });

    fileEl.addEventListener("drop", ev => {
        console.log('File(s) dropped');
        ev.preventDefault();
        let dt = ev.dataTransfer;
        let files = dt.files;
        ([...files]).forEach(file => {
            let formdata = new FormData();
            formdata.append("name", file.name);
            formdata.append("file", file);
            fetch("/_new", {method: "POST", body: formdata})
                .then(data => {
                    editor.insertText('!['+file.name+']('+ file.name +')');
                })
                .catch((error) => {
                    console.log("err:", error);
                });
        });
    });

    fileEl.addEventListener("dragover", ev => {
        // Prevent default behavior (Prevent file from being opened)
        ev.preventDefault();
    });

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
