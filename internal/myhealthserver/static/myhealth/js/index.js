$( document ).ready(function() {
  const input = document.getElementById('messageInput');
  const btn = document.getElementById('sendBtn');
  const chat = document.getElementById('chatWindow');

  btn.onclick = () => {
      if (!input.value.trim()) return;
      const msg = document.createElement('div');
      msg.className = 'message sent';
      msg.innerHTML = input.value + `<div class="message-time">${new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}</div>`;
      chat.appendChild(msg);
      input.value = '';
      chat.scrollTop = chat.scrollHeight;
  };

  input.addEventListener("keypress", (e) => {
      if (e.key === "Enter") btn.click();
  });
});
