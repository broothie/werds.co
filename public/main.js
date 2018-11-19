
// https://github.com/30-seconds/30-seconds-of-code#copytoclipboard-
function copyToClipboard(str) {
    const el = document.createElement('textarea');
    el.value = str;
    el.setAttribute('readonly', '');
    el.style.position = 'absolute';
    el.style.left = '-9999px';
    document.body.appendChild(el);
    const selected =
        document.getSelection().rangeCount > 0 ? document.getSelection().getRangeAt(0) : false;
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
    if (selected) {
        document.getSelection().removeAllRanges();
        document.getSelection().addRange(selected);
    }
}


$(function() {
    var container = $('.container');

    function handleResize() {
        console.log('resizing');
        container.textfill({ maxFontPixels: -1 });
    }

    handleResize();
    container.css('color', 'black');
    $(window).resize(handleResize);

    var copyLink = $('.copy-link');
    copyLink.click(function() {
        copyToClipboard(window.location.href);
        copyLink.text('copied!');
        setTimeout(function() {
            copyLink.text('copy link');
        }, 1250)
    });
});
