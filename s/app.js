(function(){
    const shellbtn = document.getElementById("shell-btn");
    let w = document.querySelector("#shell-window");
    w.hidden = true;
    shellbtn.addEventListener('click', event => {
	      w.hidden = !w.hidden;
    });

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
		            output.innerHTML = data.output;
		            prompt.value = "";
	          })
	          .catch((error) => {
		            console.error('Error:', error);
	          });
    });

    const menuBtn = document.querySelector("#menu-btn");
	  const menuWindow = document.querySelector("#menu-window");
	  menuWindow.hidden = true;
	  menuBtn.addEventListener('click', event => {
	      menuWindow.hidden = !menuWindow.hidden;
	  });

	  fetch('/_menu', {method: 'GET'})
	      .then(response => response.text())
	      .then(data => {
		        menuWindow.innerHTML = data;
	      })
	      .catch((error) => {
		        console.error('Error:', error);
	      });
})();
