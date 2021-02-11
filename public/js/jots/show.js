const saveDelayMs = 750;
const quillArgs = ['#body', {placeholder: 'jot something...', theme: 'snow', modules: { toolbar: '#toolbar' }}];

const script = document.currentScript;
const id = script.getAttribute('data-id');
const contents = JSON.parse(script.getAttribute('data-contents'));

document.addEventListener('DOMContentLoaded', () => {
  const quill = new Quill(...quillArgs);
  const syncStatus = document.getElementById('sync-status');

  // Start quill
  quill.setContents(contents);
  quill.focus();
  quill.setSelection(quill.getLength());

  // Sync request
  const sync = () => {
    syncStatus.innerText = 'saving...';
    const fetchOptions = {
      method: 'put',
      body: JSON.stringify({body: quill.getText(), contents: quill.getContents()})
    };

    fetch(`/jot/${id}/sync`, fetchOptions)
      .catch(console.log)
      .then((response) => response.json())
      .then((json) => syncStatus.innerText = 'saved');
  };

  // Handle text change event
  let lastUpdate = null;
  let lastTimeout = null;
  quill.on('text-change', (event) => {
    console.log(event);

    lastUpdate = Date.now();
    syncStatus.innerText = 'not saved';

    if (lastTimeout) clearTimeout(lastTimeout);
    lastTimeout = setTimeout(() => {
      if (lastUpdate && (new Date()) - lastUpdate > saveDelayMs) sync();
    }, saveDelayMs);
  });

  // Save when page is unloaded
  window.addEventListener('beforeunload', () => {
    if (lastTimeout) clearTimeout(lastTimeout);
    sync();
  });

  // Delete link
  document.getElementById('delete').addEventListener('click', (event) => {
    event.preventDefault();

    if (confirm('are you sure you want to delete this jot?')) {
      fetch(`/jot/${id}`, {method: 'delete'}).then(() => window.location.href = '/home')
    }
  });
});
