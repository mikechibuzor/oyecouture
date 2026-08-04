export const useAppLoading = () => {
  const pendingCount = useState<number>('app-loading-count', () => 0)
  const loadingLabel = useState<string>('app-loading-label', () => '')

  const isAppLoading = computed(() => pendingCount.value > 0)

  const startLoading = (label = ''): (() => void) => {
    pendingCount.value += 1

    if (label) {
      loadingLabel.value = label
    }

    let isFinished = false

    return () => {
      if (isFinished) {
        return
      }

      isFinished = true
      pendingCount.value = Math.max(0, pendingCount.value - 1)

      if (pendingCount.value === 0) {
        loadingLabel.value = ''
      }
    }
  }

  const withLoading = async <T>(
    action: () => Promise<T>,
    label = '',
  ): Promise<T> => {
    const finishLoading = startLoading(label)

    try {
      return await action()
    } finally {
      finishLoading()
    }
  }

  return {
    isAppLoading,
    loadingLabel,
    startLoading,
    withLoading,
  }
}
