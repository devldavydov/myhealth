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

        addMessage('received', 'OK');

        $chat.scrollTop($chat.prop('scrollHeight'));
    });

    $('#messageInput').on('keydown', function(e) {
        if (e.key === 'Enter') {
            $('#sendBtn').click();
        }
    });
    
    addMessage('received', 'Привет! Введи команду....')
});
