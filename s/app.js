// Text area auto grow as described here: https://stackoverflow.com/a/24676492
function auto_grow(element) {
    element.style.height = "5px";
    element.style.height = (element.scrollHeight)+"px";
}
