$( document ).ready(function() {
    var $chat = $('#chatWindow');

    function addMessage(type, msg) {
        $('<div>').
            addClass('message ' + type).
            html(msg + `<div class="message-time">${new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}</div>`).
            appendTo($chat);
    }

    $('#sendBtn').click(function() {
        let cmd = $('#messageInput').val();
        if (!cmd.trim())
            return;

        addMessage('sent', cmd);
    
        $('#messageInput').val('');

        $.ajax({
            type: "POST",
            url: "/api",
            async: false,
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({'cmd': cmd}),
            dataType: "json",
            success: function(response) {
                if (response.error !== '') {
                    // error
                    addMessage('received', response.error);
                } else if (response.isFile) {
                    // file
                } else {
                    // text
                    addMessage('received', response.textResponse);
                }
            },
            error: function(xhr, status, error) {
                addMessage('received', error);
            }
        });

        $chat.scrollTop($chat.prop('scrollHeight'));
    });

    $('#messageInput').on('keydown', function(e) {
        if (e.key === 'Enter') {
            $('#sendBtn').click();
        }
    });
    
    addMessage('received', 'Привет! Введи команду....')
});
