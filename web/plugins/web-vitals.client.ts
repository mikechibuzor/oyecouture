import {
  onCLS,
  onINP,
  onLCP,
  type CLSMetric,
  type INPMetric,
  type LCPMetric,
} from 'web-vitals'

type SupportedWebVitalMetric = CLSMetric | INPMetric | LCPMetric

interface WebVitalDetail {
  readonly name: SupportedWebVitalMetric['name']
  readonly value: number
  readonly delta: number
  readonly rating: SupportedWebVitalMetric['rating']
  readonly path: string
}

export default defineNuxtPlugin(() => {
  const navigationPath = window.location.pathname

  const reportMetric = (metric: SupportedWebVitalMetric): void => {
    const detail: WebVitalDetail = {
      name: metric.name,
      value: Number(metric.value.toFixed(metric.name === 'CLS' ? 4 : 0)),
      delta: Number(metric.delta.toFixed(metric.name === 'CLS' ? 4 : 0)),
      rating: metric.rating,
      path: navigationPath,
    }

    window.dispatchEvent(new CustomEvent<WebVitalDetail>('oyecouture:web-vital', {
      detail,
    }))

    if (import.meta.dev) {
      console.info('[web-vital]', detail)
    }
  }

  onCLS(reportMetric)
  onINP(reportMetric)
  onLCP(reportMetric)
})
