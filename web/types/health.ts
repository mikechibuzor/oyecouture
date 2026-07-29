export enum ApplicationHealthStatus {
  Degraded = 'degraded',
  Healthy = 'ok',
}

export enum ServiceHealthStatus {
  Available = 'ok',
  Unavailable = 'unavailable',
}

export interface ApiHealthResponse {
  status: ApplicationHealthStatus
  services: {
    api: ServiceHealthStatus.Available
    database: ServiceHealthStatus
  }
}
