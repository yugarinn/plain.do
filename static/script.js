function init() {
  // https://www.reddit.com/r/htmx/comments/1acmvso/hxswapoob_swaps_inner_html_of_a_component/
  htmx.defineExtension('inline', {
    isInlineSwap: (swapStyle) => true,
  })

  document.addEventListener('htmx:wsBeforeMessage', deleteTodoListener)
  addCurrentListHashToUrl()
  addShortcutsListeners()
}

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms))
}

function addCurrentListHashToUrl() {
  const currentUrl = window.location.href
  const currentUrlHash = window.location.pathname.split('/').pop()

  if (currentUrlHash == 'about') return

  const currentListHash = document
    .getElementById('list-hash')
    .getAttribute('data-list-hash')

  if (currentListHash == currentUrlHash) return

  const baseUrl = window.location.origin
  const newUrl = `${baseUrl}/${currentListHash}`

  window.history.replaceState({ path: newUrl }, '', newUrl)
}

function clearTodoInput() {
  document.getElementById('add-todo-input').value = ''
}

function copyListLinkToClipboard() {
  navigator.clipboard
    .writeText(window.location.href)
    .then(showCopyNotification)
    .catch((err) => console.error('Error copying URL to clipboard: ', err))
}

function showCopyNotification() {
  const notification = document.getElementById('copyNotification')
  notification.classList.add('show')

  setTimeout(() => notification.classList.remove('show'), 2000)
}

function focusInput(target, delayInMiliseconds = 200) {
  sleep(delayInMiliseconds).then(() => {
    const input = document.getElementById(target.id)
    if (!input) return

    const value = input.value

    input.value = ''
    input.value = value
    document.getElementById(target.id).focus()
  })
}

function deleteTodoListener(event) {
  if (!event.detail.message) return

  try {
    const message = JSON.parse(event.detail.message)
    if (message.action !== 'deleteTodo') return

    const removeTarget = document.getElementById(`todo-${message.todoID}`)
    const upperTodoElement = findUpperTodoInput(removeTarget)

    if (upperTodoElement) {
      focusInput(upperTodoElement, 0)
    }

    removeTarget.remove()
  } catch {}
}

function findUpperTodoInput(currentInput) {
  let currentForm = currentInput.closest('.todo-element')

  if (!currentForm) return null

  let previousForm = currentForm.previousElementSibling
  while (previousForm && !previousForm.classList.contains('todo-element')) {
    previousForm = previousForm.previousElementSibling
  }

  if (previousForm) return previousForm.querySelector('.todo-input')

  let nextForm = currentForm.nextElementSibling
  while (nextForm && !nextForm.classList.contains('todo-element')) {
    nextForm = nextForm.nextElementSibling
  }

  if (nextForm) return nextForm.querySelector('.todo-input')

  return null
}

function addShortcutsListeners() {
  const newTaskInput = document.getElementById('add-todo-input')

  const focusNewTask = function () {
    newTaskInput.focus()
  }

  const toggleSelectedTaskStatus = function () {
    const currentFocusedElement = document.activeElement
    if (! currentFocusedElement.classList.contains('todo-input')) return

    const checkElement = currentFocusedElement.previousElementSibling
    checkElement.click()
    focusInput(currentFocusedElement)
  }

  const deleteSelectedTask = function () {
    const currentFocusedElement = document.activeElement
    if (! currentFocusedElement.classList.contains('todo-input')) return

    currentFocusedElement.value = ''
  }

  const focusClosestTaskInput = function (direction) {
    if (! ['up', 'down'].includes(direction)) return

    const currentElement = document.activeElement
    const elements = Array.from(document.getElementsByClassName('todo-input'))

    if (elements.length == 0) return

    if (currentElement.classList.contains('todo-input')) {
      let targetTaskIndex
      const currentSelectedTaksIndex = elements.findIndex((element) => element.id == currentElement.id)

      if (direction == 'up') targetTaskIndex = currentSelectedTaksIndex - 1
      if (direction == 'down') targetTaskIndex = currentSelectedTaksIndex + 1

      if (targetTaskIndex < 0 || targetTaskIndex > elements.length - 1) {
        focusNewTask()
        return
      }

      const targetElement = elements[targetTaskIndex]
      focusInput(targetElement, 0)

      return
    }

    if (direction == 'up') {
      const targetElement = elements[elements.length - 1]
      focusInput(targetElement)
    }

    if (direction == 'down') {
      const targetElement = elements[0]
      focusInput(targetElement)
    }
  }

  document.addEventListener('keydown', function(event) {
    if (event.ctrlKey && event.shiftKey && event.key === 'N') {
      event.preventDefault()
      focusNewTask()
    }

    if (event.ctrlKey && event.shiftKey && event.code === 'Space') {
      event.preventDefault()
      toggleSelectedTaskStatus()
    }

    if (event.ctrlKey && event.shiftKey && event.key === 'D') {
      event.preventDefault()
      deleteSelectedTask()
    }

    if (event.key === 'ArrowUp' || event.keyCode === 38) {
      event.preventDefault()
      focusClosestTaskInput('up')
    }

    if (event.key === 'ArrowDown' || event.keyCode === 40) {
      event.preventDefault()
      focusClosestTaskInput('down')
    }
  })
}

window.onload = init()
