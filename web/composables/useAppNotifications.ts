export type AppNotificationTone = 'success' | 'error' | 'info'

export interface AppNotification {
  readonly id: string
  readonly tone: AppNotificationTone
  readonly title: string
  readonly message?: string
}

interface NotifyOptions {
  readonly tone: AppNotificationTone
  readonly title: string
  readonly message?: string
  readonly durationMs?: number
}

const DEFAULT_NOTIFICATION_DURATION_MS = 4200

const createNotificationId = (): string =>
  `notice-${Date.now()}-${Math.random().toString(36).slice(2)}`

export const useAppNotifications = () => {
  const notifications = useState<AppNotification[]>('app-notifications', () => [])

  const dismiss = (notificationId: string): void => {
    notifications.value = notifications.value.filter(
      (notification) => notification.id !== notificationId,
    )
  }

  const notify = (options: NotifyOptions): string => {
    const notification: AppNotification = {
      id: createNotificationId(),
      tone: options.tone,
      title: options.title,
      ...(options.message ? { message: options.message } : {}),
    }

    notifications.value = [notification, ...notifications.value].slice(0, 4)

    if (import.meta.client) {
      window.setTimeout(
        () => dismiss(notification.id),
        options.durationMs ?? DEFAULT_NOTIFICATION_DURATION_MS,
      )
    }

    return notification.id
  }

  const notifySuccess = (title: string, message?: string): string =>
    notify({ tone: 'success', title, message })

  const notifyError = (title: string, message?: string): string =>
    notify({ tone: 'error', title, message })

  const notifyInfo = (title: string, message?: string): string =>
    notify({ tone: 'info', title, message })

  return {
    notifications,
    dismiss,
    notify,
    notifySuccess,
    notifyError,
    notifyInfo,
  }
}
