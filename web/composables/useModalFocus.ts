interface ModalFocusOptions {
  readonly isOpen: Readonly<Ref<boolean>>
  readonly container: Readonly<Ref<HTMLElement | null>>
  readonly initialFocus?: Readonly<Ref<HTMLElement | null>>
  readonly onEscape: () => void
  readonly lockScroll?: boolean
}

const MODAL_FOCUSABLE_SELECTOR = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  '[tabindex]:not([tabindex="-1"])',
].join(',')

export const useModalFocus = ({
  isOpen,
  container,
  initialFocus,
  onEscape,
  lockScroll = true,
}: ModalFocusOptions): void => {
  let previouslyFocusedElement: HTMLElement | null = null
  let previousBodyOverflow = ''

  const visibleFocusableElements = (): HTMLElement[] => {
    if (!container.value) {
      return []
    }

    return Array.from(
      container.value.querySelectorAll<HTMLElement>(MODAL_FOCUSABLE_SELECTOR),
    ).filter((element) => element.offsetParent !== null)
  }

  const handleKeydown = (keyboardEvent: KeyboardEvent): void => {
    if (!isOpen.value) {
      return
    }

    if (keyboardEvent.key === 'Escape') {
      keyboardEvent.preventDefault()
      onEscape()
      return
    }

    if (keyboardEvent.key !== 'Tab') {
      return
    }

    const focusableElements = visibleFocusableElements()

    if (focusableElements.length === 0) {
      keyboardEvent.preventDefault()
      container.value?.focus()
      return
    }

    const firstFocusableElement = focusableElements[0]!
    const lastFocusableElement = focusableElements[focusableElements.length - 1]!

    if (keyboardEvent.shiftKey && document.activeElement === firstFocusableElement) {
      keyboardEvent.preventDefault()
      lastFocusableElement.focus()
    } else if (!keyboardEvent.shiftKey && document.activeElement === lastFocusableElement) {
      keyboardEvent.preventDefault()
      firstFocusableElement.focus()
    }
  }

  watch(isOpen, async (modalIsOpen) => {
    if (!import.meta.client) {
      return
    }

    if (modalIsOpen) {
      previouslyFocusedElement = document.activeElement instanceof HTMLElement
        ? document.activeElement
        : null

      if (lockScroll) {
        previousBodyOverflow = document.body.style.overflow
        document.body.style.overflow = 'hidden'
      }

      document.addEventListener('keydown', handleKeydown)
      await nextTick()
      const initialFocusTarget = initialFocus?.value
        ?? visibleFocusableElements()[0]
        ?? container.value
      initialFocusTarget?.focus()
      return
    }

    document.removeEventListener('keydown', handleKeydown)

    if (lockScroll) {
      document.body.style.overflow = previousBodyOverflow
    }

    await nextTick()
    previouslyFocusedElement?.focus()
    previouslyFocusedElement = null
  })

  onBeforeUnmount(() => {
    document.removeEventListener('keydown', handleKeydown)

    if (lockScroll) {
      document.body.style.overflow = previousBodyOverflow
    }
  })
}
