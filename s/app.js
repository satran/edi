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


(function(){
    const shellbtn = document.getElementById("shell-btn");
    let w = document.querySelector("#shell-window");
    w.hidden = true;
    shellbtn.addEventListener('click', event => {
	w.hidden = !w.hidden;
    });

    let prompt = document.querySelector("#prompt");
    let output = document.querySelector("#shell-output");
    prompt.addEventListener('keydown', ev => {
	if (ev.keyCode !== 13) return;
	
	let cmd = {cmd: prompt.value};
	fetch('_sh', {
	    method: 'POST',
	    headers: {
		'Content-Type': 'application/json',
	    },
	    body: JSON.stringify(cmd),
	})
	    .then(response => response.json())
	    .then(data => {
		console.log('Success:', data);
		output.innerHTML = data.output;
	    })
	    .catch((error) => {
		console.error('Error:', error);
	    });
    });
})();
