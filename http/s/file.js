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


    // Setup keybinding
    document.addEventListener('keydown', function(ev) {
	switch (ev.key){
	case "e":
	    if (document.location.pathname.startsWith("/edit") ||
		document.location.pathname.startsWith("/_sh")) {
		break;
	    }
	    document.location = "/edit" + document.location.pathname;
	    break;
	case "h":
	    document.location = "/"
	    break;
	case "n":
	    document.location = "/_new"
	    break;
	case "o":
	    openWindow("open");
	    break;
	case "s":
	    document.location = "/_sh"
	    break;
	}
	ev.stopPropagation();
    });
})();
