
$(function() {
    var container = $('.container');

    function handleResize() {
        console.log('resizing');
        container.textfill({ maxFontPixels: -1 });
    }

    handleResize();
    container.css('color', 'black');
    $(window).resize(handleResize);
});
