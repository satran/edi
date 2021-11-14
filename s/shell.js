(function(){
    let prompt = document.querySelector("#prompt");
    let output = document.querySelector("#shell-output");
    prompt.addEventListener('keyup', ev => {
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
		            output.innerHTML += "\n" + "$ " + prompt.value + "\n" + data.output;
		            prompt.value = "";
	          })
	          .catch((error) => {
		            console.error('Error:', error);
	          });
    });
})();
