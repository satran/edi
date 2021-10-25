// Text area auto grow as described here: https://stackoverflow.com/a/24676492
function auto_grow(element) {
    element.style.height = "5px";
    element.style.height = (element.scrollHeight)+"px";
}

function dropHandler(ev) {
    let fileInput = document.getElementById("file-input");
    console.log('File(s) dropped');

    // Prevent default behavior (Prevent file from being opened)
    ev.preventDefault();
    let dT = new DataTransfer();
    if (ev.dataTransfer.items) {
	// Use DataTransferItemList interface to access the file(s)
	for (var i = 0; i < ev.dataTransfer.items.length; i++) {
	    // If dropped items aren't files, reject them
	    if (ev.dataTransfer.items[i].kind === 'file') {
		var file = ev.dataTransfer.items[i].getAsFile();
		console.log('... file[' + i + '].name = ' + file.name);
		dT.items.add(file);
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
	    }
	}
    }
    fileInput.files = dT.files;
}

function dragOverHandler(ev) {
    // Prevent default behavior (Prevent file from being opened)
    ev.preventDefault();
}
