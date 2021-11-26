function hideParent(ev) {
    let w = ev.target.closest(".window");
    w.hidden = true;
}

(function(){

    // Ensure all windows are closed on open
    let windows = document.querySelectorAll(".window");
    windows.forEach(w => w.hidden = true);

    const openBtn = document.getElementById("open-btn");
    openBtn.addEventListener("click", function (event) {
	let parents = document.querySelectorAll(".window.open");
	parents.forEach(w => {
	    w.hidden = false;
	    let prompt = w.querySelector(".prompt");
	    prompt.focus();
	});
    });

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

    const prompt = document.querySelector(".prompt.open");
    let waiting = false;
    prompt.addEventListener("keydown", function (event) {
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
        showResults(event.target.closest(".window"), prompt.value);
    });
})();
