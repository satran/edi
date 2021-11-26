function hideParent(ev) {
    let w = ev.target.closest(".window");
    w.hidden = true;
}

(function(){
    // Ensure all windows are closed on open
    let windows = document.querySelectorAll(".window");
    windows.forEach(w => w.hidden = true);

    function openWindow(win) {
	let parents = document.querySelectorAll(".window."+win);
	parents.forEach(w => {
	    w.hidden = false;
	    let prompt = w.querySelector(".prompt");
	    prompt.focus();
	});
    }
    const windowBtns = document.querySelectorAll(".btn");
    windowBtns.forEach(w => {
	w.addEventListener("click", function (event) {
	    let win = event.target.closest(".btn").dataset.window;
	    openWindow(win);
	});
    });

    // Open files handling
    let files;
    function filterItems(arr, query) {
	return arr.filter(el => el.toLowerCase().indexOf(query.toLowerCase()) !== -1)
    }

    function showResults(parent, val) {
        let res = parent.querySelector(".file-list");
        res.innerHTML = '';
        let list = '';
        // From https://github.com/farzher/fuzzysort
        let terms = filterItems(files, val)
        for (i=0; i<terms.length; i++) {
            file = terms[i];
            list += '<a href="' + file + '">' + file + '</a>';
        }
        res.innerHTML = list;
    }

    const openPrompt = document.querySelector(".prompt.open");
    let waiting = false;
    openPrompt.addEventListener("keydown", function (event) {
	// disabling the global shortcuts to be called
	event.stopPropagation();
	if (event.key === "Escape") {
	    hideParent(event);
	}
	if (!files && !waiting) {
	    waiting = true;
	    fetch("/_ls")
		.then(response => response.json())
		.then(data => files = data);
	    return;
	}
        showResults(event.target.closest(".window"), openPrompt.value);
    });


    // Shell command handling
    let shellPrompt = document.querySelector(".prompt.shell");
    let output = document.querySelector(".shell-output");
    shellPrompt.addEventListener('keydown', ev => {
	// disabling the global shortcuts to be called
	ev.stopPropagation();
	if (event.key === "Escape") {
	    hideParent(event);
	}

	if (ev.keyCode !== 13) return;

	let cmd = {cmd: shellPrompt.value};
	fetch('_sh', {
	    method: 'POST',
	    headers: {
		'Content-Type': 'application/json',
	    },
	    body: JSON.stringify(cmd),
	})
	    .then(response => response.json())
	    .then(data => {
		output.innerHTML += "\n" + "$ " + shellPrompt.value + "\n" + data.output;
		shellPrompt.value = "";
	    })
	    .catch((error) => {
		console.error('Error:', error);
	    });
    });

    // Setup keybinding
    document.addEventListener('keydown', function(ev) {
	switch (ev.key){
	case "e":
	    if (document.location.pathname.startsWith("/edit")) {
		break;
	    }
	    document.location = "/edit" + document.location.pathname;
	    break;
	case "o":
	    openWindow("open");
	    break;
	case "s":
	    openWindow("shell");
	    break;
	}
	ev.stopPropagation();
    });
})();
