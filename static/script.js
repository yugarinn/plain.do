function init() {
  htmx.defineExtension('inline', {
    isInlineSwap: (swapStyle) => true,
  })
  // https://www.reddit.com/r/htmx/comments/1acmvso/hxswapoob_swaps_inner_html_of_a_component/
  // htmx.logAll()

  document.addEventListener('htmx:wsBeforeMessage', deleteTodoListener)
  addCurrentListHashToUrl()
}

function sleep(ms) {
  return new Promise(r => setTimeout(r, ms))
}

function addCurrentListHashToUrl() {
  const currentUrl = window.location.href
  const currentUrlHash = window.location.pathname.split("/").pop()

  if (currentUrlHash == 'about') return

  const currentListHash = document.getElementById('list-hash').getAttribute('data-list-hash')

  if (currentListHash == currentUrlHash) return

  const baseUrl = window.location.origin
  const newUrl = `${baseUrl}/${currentListHash}`

  window.history.replaceState({ path: newUrl }, '', newUrl)
}

function clearTodoInput() {
  document.getElementById('add-todo-input').value = '';
}

function copyListLinkToClipboard() {
    navigator.clipboard.writeText(window.location.href)
             .then(showCopyNotification)
             .catch(err => console.error('Error copying URL to clipboard: ', err))
}

function showCopyNotification() {
    const notification = document.getElementById('copyNotification')
    notification.classList.add('show')

    setTimeout(() => notification.classList.remove('show'), 2000)
}

function focusInput(target) {
  sleep(1000).then(() => {
    const input = document.getElementById(target.id)
    if (! input) return

    const value = input.value

    input.focus()
    input.value = ''
    input.value = value
  })
}

function deleteTodoListener(event) {
  if (! event.detail.message) return

  try {
    const message = JSON.parse(event.detail.message)
    if (message.action !== 'deleteTodo') return

    document.getElementById(`todo-${message.todoID}`).remove()
  } catch {}
}

window.onload = init()
