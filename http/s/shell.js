(function() {
    // Shell command handling
    let shellPrompt = document.querySelector(".prompt.shell");
    let output = document.querySelector(".shell-output");
    shellPrompt.addEventListener('keydown', ev => {
	// disabling the global shortcuts to be called
	ev.stopPropagation();
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
})();
