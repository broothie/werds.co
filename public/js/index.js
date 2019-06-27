
document.addEventListener('DOMContentLoaded', function() {
    var textarea = document.getElementById('werds');
    var form = document.getElementById('form');

    // Grows textarea vertically
    function growTextArea() {
        if (textarea.scrollHeight > textarea.clientHeight) {
            textarea.style.height = textarea.scrollHeight + 'px';
        }
    }

    growTextArea();

    // Adds newlines on shift+enter
    textarea.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();

            if (e.shiftKey) {
                textarea.value = textarea.value + "\n";
            } else {
                form.submit();
            }
        }
    });

    textarea.addEventListener('keyup', growTextArea)
});
