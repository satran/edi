(function(){
    // Text area auto grow as described here: https://stackoverflow.com/a/24676492
    function auto_grow(element) {
        element.style.height = "5px";
        element.style.height = (element.scrollHeight)+"px";
    }

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

	  const editor = document.querySelector("#editor");
	  auto_grow(editor);
    editor.addEventListener("input", e => {
        auto_grow(e.target);
    });

	  // Save stuff to localstorage
	  let filename = document.location.pathname;
	  let value = localStorage.getItem(filename);
	  if (value && value.length > 0) {
	      if (confirm('Load unsaved draft?')) {
	          editor.value = value;
	      }
	  }

	  let cancel;
	  editor.addEventListener("keyup", event => {
	      if (cancel) clearTimeout(cancel);
	      cancel = setTimeout(() => {
		        localStorage.setItem(filename, event.target.value);
		    }, 1000);
	  });
})();
